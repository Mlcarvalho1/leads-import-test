package middlewares

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// DummyAuth is a sample authentication middleware
func DummyAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authorization header missing"})
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		c.Locals("dummyToken", token)

		return c.Next()
	}
}
