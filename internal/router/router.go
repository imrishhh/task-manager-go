package router

import (
	"database/sql"

	"github.com/gofiber/fiber/v3"
	"github.com/nullrish/task-manager-go/internal/middleware"
)

func ConfigureRoutes(app *fiber.App, db *sql.DB) {
	api := app.Group("/api")

	r := api.Group("/auth")
	configureAuthRoutes(r, db)

	// Configure JWT middleware here
	app.Use(middleware.AuthMiddleware())

	r = api.Group("/task")
	configureTaskRouter(r, db)
	// Test JWT middleware restriction without Authorization Header
	r.Get("/restricted", func(c fiber.Ctx) error {
		return c.SendString("This one is restricted")
	})
}
