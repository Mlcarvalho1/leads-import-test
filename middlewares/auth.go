// Package middlewares...
package middlewares

import (
	"fmt"
	"os"
	"strconv"

	jwtware "github.com/gofiber/contrib/v3/jwt"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

func GetJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "secret" // fallback for development
	}

	return []byte(secret)
}

func Protected() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:     jwtware.SigningKey{Key: GetJWTSecret()},
		ErrorHandler:   jwtError,
		SuccessHandler: extractUserInfo,
	})
}

func extractUserInfo(c fiber.Ctx) error {
	token := jwtware.FromContext(c)
	if token == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Authentication required",
		})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid token claims",
		})
	}

	fmt.Println("claims", claims)

	var companyID int
	var userID string

	if userMap, ok := claims["user"].(map[string]interface{}); ok {
		if companyIDFloat, ok := userMap["company_id"].(float64); ok {
			companyID = int(companyIDFloat)
		}
		if idFloat, ok := userMap["id"].(float64); ok {
			userID = strconv.Itoa(int(idFloat))
		}
	}

	c.Locals("company_id", companyID)
	c.Locals("user_id", userID)

	return c.Next()
}

func jwtError(c fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"status": "error", "message": "Missing or malformed JWT", "data": nil})
	} else {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT", "data": nil})
	}
}
