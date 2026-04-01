package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
	apperr "github.com/nullrish/task-manager-go/internal/errors"
	"github.com/nullrish/task-manager-go/internal/model"
)

type TaskRepository interface {
	CreateTask(ctx context.Context, req *model.TaskRequest) (*model.Task, error)
	GetTaskByID(ctx context.Context, taskID uuid.UUID) (*model.Task, error)
	GetTasksByUserID(ctx context.Context, userID uuid.UUID) ([]model.Task, error)
	GetTasks(ctx context.Context) ([]model.Task, error)
	UpdateTask(ctx context.Context, req *model.TaskRequest) (*model.Task, error)
	DeleteTask(ctx context.Context, taskID uuid.UUID) error
}

type taskRepo struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) TaskRepository {
	return &taskRepo{db: db}
}

func (r *taskRepo) CreateTask(ctx context.Context, req *model.TaskRequest) (*model.Task, error) {
	query := `
		INSERT INTO tasks (task_title, task_description, status, user_id) VALUES ($1, $2, $3, $4)
		RETURNING id, task_title, task_description, status, created_at, updated_at, user_id;
	`
	var task model.Task
	err := r.db.QueryRowContext(ctx, query, req.TaskTitle, req.TaskDescription, req.Status, req.UserID).Scan(
		&task.ID,
		&task.TaskTitle,
		&task.TaskDescription,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.UserID,
	)
	if err != nil {
		log.Printf("(task_repository) - [CreateTask] failed for user %s: %v", req.UserID, err)
		return nil, &apperr.DatabaseError{Message: err.Error()}
	}
	return &task, nil
}

func (r *taskRepo) GetTaskByID(ctx context.Context, taskID uuid.UUID) (*model.Task, error) {
	query := `
		SELECT id, task_title, task_description, status, created_at, updated_at, user_id FROM tasks
		WHERE id = $1
	`
	var task model.Task
	err := r.db.QueryRowContext(ctx, query, taskID).Scan(
		&task.ID,
		&task.TaskTitle,
		&task.TaskDescription,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &apperr.NotFoundError{Resource: "task", ID: taskID.String()}
		}
		log.Printf("(task_repository) - [GetTaskByID] failed for task id %s: %v", taskID, err)
		return nil, &apperr.DatabaseError{Message: err.Error()}
	}
	return &task, nil
}

func (r *taskRepo) GetTasksByUserID(ctx context.Context, userID uuid.UUID) ([]model.Task, error) {
	query := `
		SELECT id, task_title, task_description, status, created_at, updated_at, user_id FROM tasks
		WHERE user_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &apperr.NotFoundError{Resource: "task", ID: userID.String()}
		}
		log.Printf("(task_repository) - [GetTasksByUserID] failed for user id %s: %v", userID, err)
		return nil, &apperr.DatabaseError{Message: err.Error()}
	}
	// This one shouldn't be handle for most of the time as underlying connection pool willl handle cleanup.
	// But for this one I just added it for showing that we can handle closing error as well if necessary.
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("(task_repository) - [GetTasksByUserID] Couldn't close the rows: %v", err)
		}
	}()

	var tasks []model.Task
	for rows.Next() {
		var t model.Task
		// Following go idiom to handle errors
		if err := rows.Scan(&t.ID, &t.TaskTitle, &t.TaskDescription, &t.Status, &t.CreatedAt, &t.UpdatedAt, &t.UserID); err != nil {
			log.Printf("(task_repository) - [GetTasksByUserID] Cannot scan the row: %v", err)
			continue
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (r *taskRepo) GetTasks(ctx context.Context) ([]model.Task, error) {
	query := `
		SELECT id, task_title, task_description, status, created_at, updated_at, user_id FROM tasks;
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &apperr.NotFoundError{Resource: "task", ID: ""}
		}
		log.Printf("(task_repository) - [GetTasks] Cannot scan the row: %v", err)
		return nil, fmt.Errorf("GetTasks: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("(task_repository) - [GetTasks] Couldn't close the rows: %v", err)
		}
	}()

	var tasks []model.Task
	for rows.Next() {
		var t model.Task
		if err := rows.Scan(&t.ID, &t.TaskTitle, &t.TaskDescription, &t.Status, &t.CreatedAt, &t.UpdatedAt, &t.UserID); err != nil {
			log.Printf("(task_repository) - [GetTasks] Cannot scan the row: %v", err)
			continue
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (r *taskRepo) UpdateTask(ctx context.Context, req *model.TaskRequest) (*model.Task, error) {
	query := `
		UPDATE tasks SET task_title = $1, task_description = $2, user_id = $3, status = $4 WHERE id = $5
		RETURNING id, task_title, task_description, status, created_at, updated_at, user_id;
	`
	var task model.Task
	err := r.db.QueryRowContext(ctx, query,
		req.TaskTitle,
		req.TaskDescription,
		req.UserID,
		req.Status,
		req.ID,
	).Scan(
		&task.ID,
		&task.TaskTitle,
		&task.TaskDescription,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.UserID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &apperr.NotFoundError{Resource: "task", ID: req.ID.String()}
		}
		log.Printf("(task_repository) - [UpdateTask] Cannot update the task id %s: %v", req.ID, err)
		return nil, &apperr.DatabaseError{Message: err.Error()}
	}
	return &task, nil
}

func (r *taskRepo) DeleteTask(ctx context.Context, taskID uuid.UUID) error {
	query := `
		DELETE FROM tasks WHERE id = $1;
	`
	result, err := r.db.ExecContext(ctx, query, taskID)
	if err != nil {
		log.Printf("(task_repository) - [DeleteTask] Couldn't delete the task id %d: %v", taskID, err)
		return &apperr.DatabaseError{Message: err.Error()}
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return &apperr.DatabaseError{Message: "something went wrong while fecthing row affected"}
	}
	if rowsAffected == 0 {
		return &apperr.NotFoundError{Resource: "task", ID: taskID.String()}
	}
	return nil
}
