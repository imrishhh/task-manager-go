package router

import (
	"database/sql"

	"github.com/gofiber/fiber/v3"
	"github.com/nullrish/task-manager-go/internal/middleware/jwt"
	authrouter "github.com/nullrish/task-manager-go/internal/router/auth_router"
)

func ConfigureRoutes(app *fiber.App, db *sql.DB) {
	api := app.Group("/api")

	auth := api.Group("/auth")

	authrouter.ConfigureAuthRoutes(auth, db)

	// Configure JWT middleware here
	app.Use(jwt.Middleware())

	// Test JWT middleware restriction without Authorization Header
	auth.Get("/restricted", func(c fiber.Ctx) error {
		return c.SendString("This one is restricted")
	})
}
