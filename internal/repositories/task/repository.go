package task

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	model "github.com/nullrish/task-manager-go/internal/models/task"
)

type TaskRepository interface {
	CreateTask(ctx context.Context, task *model.TaskRequest) error
	GetTaskByID(ctx context.Context, tID uuid.UUID) (*model.Task, error)
	GetTasksByUserID(ctx context.Context, uID uuid.UUID) ([]model.Task, error)
	GetTasks(ctx context.Context) ([]model.Task, error)
	UpdateTask(ctx context.Context, task *model.TaskRequest) error
	DeleteTask(ctx context.Context, tID uuid.UUID) error
}

type taskRepo struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) TaskRepository {
	return &taskRepo{db: db}
}

func (tr *taskRepo) CreateTask(ctx context.Context, task *model.TaskRequest) error {
	return nil
}

func (tr *taskRepo) GetTaskByID(ctx context.Context, tID uuid.UUID) (*model.Task, error) {
	return nil, nil
}

func (tr *taskRepo) GetTasksByUserID(ctx context.Context, uID uuid.UUID) ([]model.Task, error) {
	return nil, nil
}

func (tr *taskRepo) GetTasks(ctx context.Context) ([]model.Task, error) {
	return nil, nil
}

func (tr *taskRepo) UpdateTask(ctx context.Context, task *model.TaskRequest) error {
	return nil
}

func (tr *taskRepo) DeleteTask(ctx context.Context, tID uuid.UUID) error {
	return nil
}
