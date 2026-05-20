package response

import "github.com/gofiber/fiber/v2"

type Response[T any] struct {
	Success bool   `json:"success"`
	Data    T      `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

type PaginatedResponse[T any] struct {
	Success bool   `json:"success"`
	Data    []T    `json:"data"`
	Page    int    `json:"page"`
	Limit   int    `json:"limit"`
	Total   int64  `json:"total"`
	Message string `json:"message,omitempty"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    int    `json:"code"`
}

type AppError struct {
	Code    int
	Message string
	Details string
}

func (e *AppError) Error() string {
	return e.Message
}

// SuccessResponse returns a successful response
func SuccessResponse[T any](ctx *fiber.Ctx, statusCode int, data T, message string) error {
	return ctx.Status(statusCode).JSON(Response[T]{
		Success: true,
		Data:    data,
		Message: message,
	})
}

// PaginatedSuccessResponse returns a paginated successful response
func PaginatedSuccessResponse[T any](
	ctx *fiber.Ctx,
	statusCode int,
	data []T,
	page int,
	limit int,
	total int64,
	message string,
) error {
	return ctx.Status(statusCode).JSON(PaginatedResponse[T]{
		Success: true,
		Data:    data,
		Page:    page,
		Limit:   limit,
		Total:   total,
		Message: message,
	})
}

// ErrorResponseJSON returns an error response
func ErrorResponseJSON(ctx *fiber.Ctx, statusCode int, errMsg string, code int) error {
	return ctx.Status(statusCode).JSON(ErrorResponse{
		Success: false,
		Error:   errMsg,
		Code:    code,
	})
}
