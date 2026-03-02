package user

import (
	"context"
	"database/sql"

	"github.com/google/uuid"

	userModel "github.com/nullrish/task-manager-go/internal/models/user"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *userModel.UserRequest) error
	GetUserByUsername(ctx context.Context, username string) (*userModel.User, error)
	GetUserByEmail(ctx context.Context, email string) (*userModel.User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, user *userModel.UserRequest) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepo{db: db}
}

func (repo *userRepo) CreateUser(ctx context.Context, user *userModel.UserRequest) error {
	query := `
		INSERT INTO users (username, email, password)
		VALUES ($1, $2, $3)
	`

	_, err := repo.db.ExecContext(ctx, query, user.Username, user.Email, user.Password)
	return err
}

func (repo *userRepo) GetUserByUsername(ctx context.Context, username string) (*userModel.User, error) {
	query := `
		SELECT id, username, email, password, created_at, updated_at
		FROM users WHERE username = $1
	`
	var user userModel.User
	err := repo.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (repo *userRepo) GetUserByEmail(ctx context.Context, email string) (*userModel.User, error) {
	query := `
		SELECT id, username, email, password, created_at, updated_at
		FROM users WHERE email = $1
	`
	var user userModel.User
	err := repo.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (repo *userRepo) UpdateUser(ctx context.Context, id uuid.UUID, user *userModel.UserRequest) error {
	query := `
		UPDATE users
		SET username = $1, email = $2, password = $3
		WHERE id = $4
	`
	_, err := repo.db.ExecContext(ctx, query, user.Username, user.Email, user.Password, id)
	return err
}

func (repo *userRepo) DeleteUser(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM users
		WHERE id = $1
	`
	_, err := repo.db.ExecContext(ctx, query, id)
	return err
}
