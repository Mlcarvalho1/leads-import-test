package routes

import (
	"leads-import/middlewares"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

// SetupRoutes configures all application routes
func SetupRoutes(app *fiber.App) {
	api := app.Group("/", logger.New())

	api.Get("/health", func(c fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "OK"})
	})

	api.Use(middlewares.Protected())

	RegisterLeadRoutes(api)
}
