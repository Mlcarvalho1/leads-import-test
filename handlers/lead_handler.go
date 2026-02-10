package handlers

import (
	"encoding/json"
	"strconv"
	"strings"

	"leads-import/models"
	"leads-import/services"
	"leads-import/validation"

	"github.com/gofiber/fiber/v3"
)

func ImportLeads(c fiber.Ctx) error {
	// Extract company_id and user_id from JWT (set by auth middleware)
	companyID, _ := c.Locals("company_id").(int)
	userIDStr, _ := c.Locals("user_id").(string)
	userID, _ := strconv.Atoi(userIDStr)

	if companyID == 0 || userID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid token: missing company_id or user_id",
		})
	}

	// Extract bearer token
	authHeader := c.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" || token == authHeader {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "missing bearer token",
		})
	}

	// Check permission
	if err := services.CheckImportPermission(token, companyID); err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Parse "data" form field as JSON
	dataStr := c.FormValue("data")
	if dataStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "missing 'data' form field",
		})
	}

	var req models.ImportRequest
	if err := json.Unmarshal([]byte(dataStr), &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid 'data' JSON",
			"details": err.Error(),
		})
	}

	// Validate request fields
	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "name is required"})
	}
	if len(req.Name) > 255 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "name must be at most 255 characters"})
	}
	if req.AccountID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "account_id is required"})
	}
	if req.SourceID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "source_id is required"})
	}
	if len(req.TagIDs) > 5 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "max 5 tag_ids allowed"})
	}

	// Parse file
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "missing 'file': use form field 'file' with a CSV or Excel file",
		})
	}

	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to read uploaded file",
		})
	}
	defer file.Close()

	rows, rowErrors, err := validation.ParseFile(file, fileHeader.Filename)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "file validation failed",
			"details": err.Error(),
		})
	}

	if len(rowErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":  "file contains invalid rows",
			"errors": rowErrors,
		})
	}

	// Start import
	importService := services.GetImportService()
	importID, err := importService.StartImport(services.StartImportInput{
		Request:   req,
		Rows:      rows,
		CompanyID: companyID,
		UserID:    userID,
		Token:     token,
	})
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"import_id": importID,
	})
}
