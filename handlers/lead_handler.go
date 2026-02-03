package handlers

import (
	"github.com/gofiber/fiber/v2"
	"your-app/services"
)

// ImportLeads handles POST /leads/import with multipart form file "file" (CSV).
func ImportLeads(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "missing file: use form field 'file' with a CSV file",
		})
	}

	f, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to read uploaded file",
		})
	}
	defer f.Close()

	imported, err := services.ImportLeadsFromCSV(f)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid CSV",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":  "import completed",
		"imported": imported,
	})
}
