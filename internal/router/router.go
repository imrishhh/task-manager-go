package router

import (
	"database/sql"

	"github.com/gofiber/fiber/v3"
	"github.com/nullrish/task-manager-go/internal/middleware/jwt"
	"github.com/nullrish/task-manager-go/internal/router/auth"
)

func ConfigureRoutes(app *fiber.App, db *sql.DB) {
	api := app.Group("/api")

	r := api.Group("/auth")

	auth.ConfigureAuthRoutes(r, db)

	// Configure JWT middleware here
	app.Use(jwt.Middleware())

	// Test JWT middleware restriction without Authorization Header
	r.Get("/restricted", func(c fiber.Ctx) error {
		return c.SendString("This one is restricted")
	})
}
