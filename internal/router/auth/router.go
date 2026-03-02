package auth

import (
	"database/sql"

	"github.com/gofiber/fiber/v3"
	authHandler "github.com/nullrish/task-manager-go/internal/handlers/auth"
	userRepo "github.com/nullrish/task-manager-go/internal/repositories/user"
	authService "github.com/nullrish/task-manager-go/internal/services/auth"
)

func ConfigureAuthRoutes(r fiber.Router, db *sql.DB) {
	repo := userRepo.NewUserRepository(db)
	service := authService.NewService(repo)
	handler := authHandler.NewHandler(service)
	r.Post("/register", handler.RegisterUser)
	r.Post("/login", handler.LoginUser)
	r.Post("/refresh", handler.RefreshToken)
}
