package routes

import (
	"leads-import/handlers"

	"github.com/gofiber/fiber/v3"
)

func RegisterLeadRoutes(api fiber.Router) {
	api.Post("/import", handlers.ImportLeads)
}
