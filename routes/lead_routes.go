package routes

import (
	"github.com/gofiber/fiber/v2"
	"your-app/handlers"
)

// RegisterLeadRoutes registers lead-related routes
func RegisterLeadRoutes(api fiber.Router) {
	leads := api.Group("/leads")
	leads.Post("/import", handlers.ImportLeads)
}
