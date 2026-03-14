package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	apperr "github.com/nullrish/task-manager-go/internal/errors"
	"github.com/nullrish/task-manager-go/internal/model"
	"github.com/nullrish/task-manager-go/internal/service"
	"github.com/nullrish/task-manager-go/internal/util/validator"
)

type TaskHandler struct {
	service *service.TaskService
}

func NewTaskHandler(s *service.TaskService) *TaskHandler {
	return &TaskHandler{s}
}

func (h *TaskHandler) CreateTask(c fiber.Ctx) error {
	req := &model.TaskRequest{}
	if err := c.Bind().Body(req); err != nil {
		return &apperr.ValidationError{Message: "Invalid input"}
	}

	if req.TaskTitle == "" {
		return &apperr.ValidationError{Field: "task_title", Message: "please enter task title"}
	}

	if !validator.ValidateStatus(req.Status) {
		return &apperr.ValidationError{Field: "status", Message: "invalid status value"}
	}

	task, err := h.service.CreateTask(c.Context(), req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(&model.Response{
		Message: "Successfully created task!",
		Data:    task,
	})
}

func (h *TaskHandler) GetUserTasks(c fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return &apperr.ValidationError{Field: "id", Message: "invalid user id"}
	}
	tasks, err := h.service.GetTaskByUserID(c.Context(), userID)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(&model.Response{
		Message: "Successfully retrieved tasks!",
		Data:    tasks,
	})
}

func (h *TaskHandler) GetTask(c fiber.Ctx) error {
	taskID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return &apperr.ValidationError{Field: "id", Message: "invalid task id"}
	}
	task, err := h.service.GetTaskByTaskID(c.Context(), taskID)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(&model.Response{
		Message: "Successfully retrieved tasks!",
		Data:    task,
	})
}

func (h *TaskHandler) GetTasks(c fiber.Ctx) error {
	tasks, err := h.service.GetTasks(c.Context())
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(&model.Response{
		Message: "Successfully retreived tasks!",
		Data:    tasks,
	})
}

func (h *TaskHandler) UpdateTask(c fiber.Ctx) error {
	req := &model.TaskRequest{}
	if err := c.Bind().Body(req); err != nil {
		return &apperr.ValidationError{Message: "invalid input"}
	}

	if !validator.ValidateStatus(req.Status) {
		return &apperr.ValidationError{Field: "status", Message: "invalid status value"}
	}

	task, err := h.service.UpdateTask(c.Context(), req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(&model.Response{
		Message: "Successfully updated tasks!",
		Data:    task,
	})
}

func (h *TaskHandler) DeleteTask(c fiber.Ctx) error {
	taskID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return &apperr.ValidationError{Field: "id", Message: "invalid task id"}
	}
	if err := h.service.DeleteTask(c.Context(), taskID); err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}
