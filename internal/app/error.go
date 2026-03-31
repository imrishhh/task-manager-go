package app

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v3"

	apperr "github.com/nullrish/task-manager-go/internal/errors"
	"github.com/nullrish/task-manager-go/internal/model"
)

func errorHandler(c fiber.Ctx, err error) error {
	var (
		database       *apperr.DatabaseError
		business       *apperr.BusinessError
		conflict       *apperr.ConflictError
		notFound       *apperr.NotFoundError
		validation     *apperr.ValidationError
		unknown        *apperr.UnknownError
		internalServer *apperr.InternalServerError
		authentication *apperr.AuthenticationError
	)
	switch {
	case errors.As(err, &database):
		return c.Status(fiber.StatusInternalServerError).JSON(&model.Response{
			Message: "Something went wrong!",
			Error:   fmt.Errorf("failed to query database"),
		})
	case errors.As(err, &business):
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&model.Response{
			Message: "Something went wrong!",
			Error:   err,
		})
	case errors.As(err, &conflict):
		return c.Status(fiber.StatusConflict).JSON(&model.Response{
			Message: "Conflicting input!",
			Error:   err,
		})
	case errors.As(err, &notFound):
		return c.Status(fiber.StatusNotFound).JSON(&model.Response{
			Message: "Resource Not Found!",
			Error:   err,
		})
	case errors.As(err, &validation):
		return c.Status(fiber.StatusBadRequest).JSON(&model.Response{
			Message: "Invalid Input!",
			Error:   err,
		})
	case errors.As(err, &internalServer):
		return c.Status(fiber.StatusInternalServerError).JSON(&model.Response{
			Message: "Something went wrong while processing the request.",
			Error:   err,
		})
	case errors.As(err, &unknown):
		return c.Status(fiber.StatusInternalServerError).JSON(&model.Response{
			Message: "Something went wrong!",
			Error:   fmt.Errorf("unknown Error"),
		})
	case errors.As(err, &authentication):
		return c.Status(fiber.StatusBadRequest).JSON(&model.Response{
			Message: "Invalid username or password",
			Error:   err,
		})
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(&model.Response{
			Message: "",
		})
	}
}
