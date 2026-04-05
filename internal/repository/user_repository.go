package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	apperr "github.com/nullrish/task-manager-go/internal/errors"
	"github.com/nullrish/task-manager-go/internal/model"
)

type UserRepository interface {
	CreateUser(ctx context.Context, req *model.UserRequest) (*model.User, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	UpdateUser(ctx context.Context, userID uuid.UUID, req *model.UserRequest) (*model.User, error)
	DeleteUser(ctx context.Context, userID uuid.UUID) error
	VerifyUser(ctx context.Context, userID uuid.UUID) error
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) CreateUser(ctx context.Context, req *model.UserRequest) (*model.User, error) {
	query := `
		INSERT INTO users (username, email, password)
		VALUES ($1, $2, $3)
		RETURNING id, username, email, password, created_at, updated_at, verified;
	`
	var user model.User
	err := r.db.QueryRowContext(ctx, query, req.Username, req.Email, req.Password).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Verified,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			log.Printf("(user_repository) - [CreateUser] failed! code=%s constraint=%s detail=%s table=%s",
				pgErr.Code,
				pgErr.ConstraintName,
				pgErr.Detail,
				pgErr.TableName,
			)
			switch pgErr.ConstraintName {
			case "users_username_key":
				return nil, &apperr.ConflictError{Message: "username already exists"}
			case "users_email_key":
				return nil, &apperr.ConflictError{Message: "email already exists"}
			case "users_pkey":
				return nil, &apperr.ConflictError{Message: "user already exist"}
			default:
				return nil, &apperr.ConflictError{Message: "duplicate input field"}
			}
		}
		log.Printf("(user_repository) - [CreateUser] failed to create user %s: %v", req.Username, err)
		return nil, &apperr.DatabaseError{Message: err.Error()}
	}
	return &user, nil
}

func (r *userRepo) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	query := `
		SELECT id, username, email, password, created_at, updated_at, verified
		FROM users WHERE username = $1
	`
	var user model.User
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Verified,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &apperr.NotFoundError{Resource: "user", ID: username}
		}
		log.Printf("(user_repository) - [GetUserByUsername] failed to get username %s: %v", username, err)
		return nil, &apperr.DatabaseError{Message: err.Error()}
	}
	return &user, nil
}

func (r *userRepo) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, username, email, password, created_at, updated_at, verified
		FROM users WHERE email = $1
	`
	var user model.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Verified,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &apperr.NotFoundError{Resource: "user", ID: email}
		}
		log.Printf("(user_repository) - [GetUserByEmail] failed to get email %s: %v", email, err)
		return nil, &apperr.DatabaseError{Message: err.Error()}
	}
	return &user, nil
}

func (r *userRepo) UpdateUser(ctx context.Context, userID uuid.UUID, req *model.UserRequest) (*model.User, error) {
	query := `
		UPDATE users
		SET username = $1, email = $2, password = $3
		WHERE id = $4
		RETURNING id, username, email, password, created_at, updated_at, verified;
	`

	var user model.User
	err := r.db.QueryRowContext(ctx, query, req.Username, req.Email, req.Password, userID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Verified,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &apperr.NotFoundError{Resource: "user", ID: userID.String()}
		}
		log.Printf("(user_repository) - [UpdateUser] failed to update user %s: %v", userID, err)
		return nil, &apperr.DatabaseError{Message: err.Error()}
	}
	return &user, nil
}

func (r *userRepo) VerifyUser(ctx context.Context, userID uuid.UUID) error {
	query := `
		UPDATE users
		SET verified = true
		WHERE id = $1
		RETURNING id, username, email, password, created_at, updated_at, verified
	`
	var user model.User
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Verified,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return &apperr.NotFoundError{Resource: "user", ID: userID.String()}
		}
		log.Printf("(user_repository) - [VerifyUser] failed to verify user %s: %v", userID, err)
		return &apperr.DatabaseError{Message: err.Error()}
	}
	return nil
}

func (r *userRepo) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	query := `
		DELETE FROM users
		WHERE id = $1
	`
	result, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		log.Printf("(user_repository) - [DeleteUser] failed to delete user %s: %v", userID, err)
		return &apperr.DatabaseError{Message: err.Error()}
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return &apperr.DatabaseError{Message: "something went wrong while fetching rows affected"}
	}
	if rowsAffected == 0 {
		return &apperr.NotFoundError{Resource: "user", ID: userID.String()}
	}
	return nil
}
