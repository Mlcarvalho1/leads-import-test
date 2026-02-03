package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// SetupRoutes configures all application routes
func SetupRoutes(app *fiber.App) {
	api := app.Group("/", logger.New())

	api.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "OK"})
	})

	RegisterLeadRoutes(api)
}
