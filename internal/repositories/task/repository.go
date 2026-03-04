package task

import (
	"context"
	"database/sql"
	"log"

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
	query := `
		INSERT INTO tasks (task_title, task_description, user_id) VALUES ($1, $2, $3);
	`
	_, err := tr.db.ExecContext(ctx, query, task.TaskTitle, task.TaskDescription, task.UserID)
	return err
}

func (tr *taskRepo) GetTaskByID(ctx context.Context, tID uuid.UUID) (*model.Task, error) {
	query := `
		SELECT id, task_title, task_description, created_at, updated_at, user_id FROM tasks
		WHERE id = $1
	`
	var task model.Task
	err := tr.db.QueryRowContext(ctx, query, tID).Scan(
		&task.ID,
		&task.TaskTitle,
		&task.TaskDescription,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &task, nil
}

func (tr *taskRepo) GetTasksByUserID(ctx context.Context, uID uuid.UUID) ([]model.Task, error) {
	query := `
		SELECT id, task_title, task_description, created_at, updated_at, user_id FROM tasks
		WHERE user_id = $1
	`
	rows, err := tr.db.QueryContext(ctx, query, uID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	// This one shouldn't be handle for most of the time as underlying connection pool willl handle cleanup.
	// But for this one I just added it for showing that we can handle it too.
	defer func() {
		if err := rows.Close(); err != nil {
			log.Println("Couldn't close the rows:", err.Error())
		}
	}()

	var tasks []model.Task
	for rows.Next() {
		var t model.Task
		// Following go idiom to handle errors
		if err := rows.Scan(&t.ID, &t.TaskTitle, &t.TaskDescription, &t.CreatedAt, &t.UpdatedAt, &t.UserID); err != nil {
			log.Println("Cannot scan the row:", err.Error())
			continue
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (tr *taskRepo) GetTasks(ctx context.Context) ([]model.Task, error) {
	query := `
		SELECT id, task_title, task_description, created_at, updated_at, user_id FROM tasks;
	`
	rows, err := tr.db.QueryContext(ctx, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Println("Couldn't close the rows:", err.Error())
		}
	}()

	var tasks []model.Task
	for rows.Next() {
		var t model.Task
		if err := rows.Scan(&t.ID, &t.TaskTitle, &t.TaskDescription, &t.CreatedAt, &t.UpdatedAt, &t.UserID); err != nil {
			log.Println("Cannot scan the row:", err.Error())
			continue
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (tr *taskRepo) UpdateTask(ctx context.Context, task *model.TaskRequest) error {
	query := `
		UPDATE tasks SET task_title = $1, task_description = $2, user_id = $3 WHERE id = $4;
	`
	if _, err := tr.db.ExecContext(ctx, query,
		task.TaskTitle,
		task.TaskDescription,
		task.UserID,
		task.ID); err == sql.ErrNoRows {
		return nil
	} else {
		return err
	}
}

func (tr *taskRepo) DeleteTask(ctx context.Context, tID uuid.UUID) error {
	query := `
		DELETE FROM tasks WHERE id = $1;
	`
	if _, err := tr.db.ExecContext(ctx, query, tID); err == sql.ErrNoRows {
		return nil
	} else {
		return err
	}
}
