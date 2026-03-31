package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	apperr "github.com/nullrish/task-manager-go/internal/errors"
	"github.com/nullrish/task-manager-go/internal/middleware"
	"github.com/nullrish/task-manager-go/internal/model"
	"github.com/nullrish/task-manager-go/internal/repository"
	"github.com/nullrish/task-manager-go/internal/util/hashing"
)

type AuthService struct {
	repo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) *AuthService {
	return &AuthService{repo}
}

func (s *AuthService) RegisterUser(ctx context.Context, req *model.UserRequest) (*model.User, error) {
	// If fields are empty then return error of missing field
	// if user.Username == "" || user.Email == "" || user.Password == "" {
	// 	return nil, errors.New("missing fields required")
	// }
	// // Validate email address
	// if !validator.ValidateEmail(user.Email) {
	// 	return errors.New("invalid email address")
	// }
	// // Validate username
	// if !validator.ValidateUsername(user.Username) {
	// 	return errors.New("username can only contain letters, number and underscores (3-20 characters)")
	// }
	// // Validate password
	// if !validator.ValidatePassword(user.Password) {
	// 	return errors.New("password must be 8-32 chars, include uppercase, lowercase, number, and special char")
	// }
	// Check if username already exists.
	var err error
	req.Password, err = hashing.HashPassword(req.Password)
	if err != nil {
		return nil, &apperr.InternalServerError{Message: "failed to hash password"}
	}
	return s.repo.CreateUser(ctx, req)
}

func (s *AuthService) LoginUser(ctx context.Context, req *model.UserRequest) (*model.UserLoginResponse, error) {
	// If fields are empty then return error of missing field
	var user *model.User
	var field string
	var err error
	if req.Email != "" {
		user, err = s.repo.GetUserByEmail(ctx, req.Email)
		field = "email"
		if user == nil {
			return nil, err
		}
	} else {
		user, err = s.repo.GetUserByUsername(ctx, req.Username)
		field = "username"
		if user == nil {
			return nil, err
		}
	}
	matched := hashing.CheckHashedPassword(req.Password, user.Password)
	if matched {
		token, err := middleware.GenerateNewUserToken(user.ID)
		if err != nil {
			return nil, &apperr.InternalServerError{Message: "failed to generate login token"}
		}
		return &model.UserLoginResponse{
			User:  user,
			Token: token,
		}, nil

	} else {
		message := fmt.Sprintf("Invalid %s or password", field)
		return nil, &apperr.AuthenticationError{Message: message}
	}
}

func (s *AuthService) GenerateRefreshToken(ctx context.Context, userID uuid.UUID) (string, error) {
	token, err := middleware.GenerateNewUserToken(userID)
	if err != nil {
		return "", &apperr.InternalServerError{Message: "failed to generate login token"}
	}
	return token, nil
}
