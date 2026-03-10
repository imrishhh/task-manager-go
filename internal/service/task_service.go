package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/nullrish/task-manager-go/internal/model"
	"github.com/nullrish/task-manager-go/internal/repository"
)

type TaskService struct {
	repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) *TaskService {
	return &TaskService{repo}
}

func (s *TaskService) CreateTask(ctx context.Context, req *model.TaskRequest) (*model.Task, error) {
	task, err := s.repo.CreateTask(ctx, req)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskService) GetTaskByTaskID(ctx context.Context, taskID uuid.UUID) (*model.Task, error) {
	task, err := s.repo.GetTaskByID(ctx, taskID)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskService) GetTaskByUserID(ctx context.Context, userID uuid.UUID) ([]model.Task, error) {
	task, err := s.repo.GetTasksByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskService) GetTasks(ctx context.Context) ([]model.Task, error) {
	task, err := s.repo.GetTasks(ctx)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskService) UpdateTask(ctx context.Context, req *model.TaskRequest) (*model.Task, error) {
	task, err := s.repo.UpdateTask(ctx, req)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskService) DeleteTask(ctx context.Context, taskID uuid.UUID) error {
	return s.repo.DeleteTask(ctx, taskID)
}
