package errors

import "github.com/gofiber/fiber/v2"

const (
	ErrInvalidInput   = "invalid input"
	ErrUnauthorized   = "unauthorized"
	ErrForbidden      = "forbidden"
	ErrNotFound       = "not found"
	ErrConflict       = "conflict"
	ErrInternalServer = "internal server error"
)

type AppError struct {
	Code    int
	Message string
	Details string
}

func NewAppError(code int, message, details string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

func (e *AppError) Error() string {
	return e.Message
}

// Common errors
func BadRequest(details string) *AppError {
	return NewAppError(fiber.StatusBadRequest, ErrInvalidInput, details)
}

func Unauthorized(details string) *AppError {
	return NewAppError(fiber.StatusUnauthorized, ErrUnauthorized, details)
}

func Forbidden(details string) *AppError {
	return NewAppError(fiber.StatusForbidden, ErrForbidden, details)
}

func NotFound(details string) *AppError {
	return NewAppError(fiber.StatusNotFound, ErrNotFound, details)
}

func Conflict(details string) *AppError {
	return NewAppError(fiber.StatusConflict, ErrConflict, details)
}

func InternalServerError(details string) *AppError {
	return NewAppError(fiber.StatusInternalServerError, ErrInternalServer, details)
}
