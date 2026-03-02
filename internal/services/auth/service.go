package auth

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/nullrish/task-manager-go/internal/middleware/jwt"
	userModel "github.com/nullrish/task-manager-go/internal/models/user"
	"github.com/nullrish/task-manager-go/internal/repositories/user"
	"github.com/nullrish/task-manager-go/internal/util/hashing"
	"github.com/nullrish/task-manager-go/internal/util/validator"
)

type Service struct {
	repo user.UserRepository
}

func NewService(repo user.UserRepository) *Service {
	return &Service{repo}
}

func (s *Service) RegisterUser(ctx context.Context, user *userModel.UserRequest) error {
	// If fields are empty then return error of missing field
	if user.Username == "" || user.Email == "" || user.Password == "" {
		return errors.New("missing fields required")
	}
	// Validate email address
	if !validator.ValidateEmail(user.Email) {
		return errors.New("invalid email address")
	}
	// Validate username
	if !validator.ValidateUsername(user.Username) {
		return errors.New("username can only contain letters, number and underscores (3-20 characters)")
	}
	// Validate password
	if !validator.ValidatePassword(user.Password) {
		return errors.New("password must be 8-32 chars, include uppercase, lowercase, number, and special char")
	}
	// Check if username already exists.
	if existing, err := s.repo.GetUserByUsername(ctx, user.Username); err != nil {
		return err
	} else if existing != nil {
		return errors.New("username already exists")
	}
	// Check if email is already registered.
	if existing, err := s.repo.GetUserByEmail(ctx, user.Email); err != nil {
		return err
	} else if existing != nil {
		return errors.New("email already exists")
	}
	var err error
	user.Password, err = hashing.HashPassword(user.Password)
	if err != nil {
		return errors.New("failed to hash the password")
	}
	return s.repo.CreateUser(ctx, user)
}

func (s *Service) LoginUser(ctx context.Context, user *userModel.UserRequest) (string, error) {
	// If fields are empty then return error of missing field
	if user.Username == "" && user.Email == "" {
		return "", errors.New("enter username and email")
	}

	var u *userModel.User
	var err error
	if user.Email != "" {
		u, err = s.repo.GetUserByEmail(ctx, user.Email)
		if err != nil {
			return "", err
		}
		if u == nil {
			return "", errors.New("invalid login or password")
		}
	} else {
		u, err = s.repo.GetUserByUsername(ctx, user.Username)
		if err != nil {
			return "", err
		}
		if u == nil {
			return "", errors.New("invalid login or password")
		}
	}
	matched := hashing.CheckHashedPassword(user.Password, u.Password)
	if matched {
		return jwt.GenerateNewUserToken(u.ID)
	} else {
		return "", errors.New("invalid login or password")
	}
}

func (s *Service) GenerateRefreshToken(ctx context.Context, userID uuid.UUID) (string, error) {
	if userID.String() == "" {
		return "", errors.New("invalid user id")
	}
	return jwt.GenerateNewUserToken(userID)
}
