package router

import (
	"database/sql"

	"github.com/gofiber/fiber/v3"
	"github.com/nullrish/task-manager-go/internal/handler"
	"github.com/nullrish/task-manager-go/internal/repository"
	"github.com/nullrish/task-manager-go/internal/service"
)

func configureTaskRouter(r fiber.Router, db *sql.DB) {
	repo := repository.NewTaskRepository(db)
	s := service.NewTaskService(repo)
	h := handler.NewTaskHandler(s)
	r.Post("/create", h.CreateTask)
	r.Put("/update", h.UpdateTask)
	r.Get("/by-task-id/:id", h.GetTask)
	r.Get("/by-user-id/:id", h.GetUserTasks)
	r.Get("/all", h.GetTasks)
	r.Delete("/:id", h.DeleteTask)
}
