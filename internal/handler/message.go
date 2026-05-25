package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	messageusecase "github.com/quoteyouros/backend/internal/usecase/message"
	apperrors "github.com/quoteyouros/backend/pkg/errors"
	"github.com/quoteyouros/backend/pkg/logger"
	"github.com/quoteyouros/backend/pkg/response"
)

type MessageHandler struct {
	usecase *messageusecase.MessageUseCase
}

// NewMessageHandler creates a new message handler
func NewMessageHandler(usecase *messageusecase.MessageUseCase) *MessageHandler {
	return &MessageHandler{usecase: usecase}
}

// GetAllMessages retrieves contact submissions (admin only)
// GET /api/messages
func (h *MessageHandler) GetAllMessages(c *fiber.Ctx) error {
	// Extract pagination params from query
	pageStr := c.Query("page", "1")
	limitStr := c.Query("limit", "20")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		logger.Warn("getAllMessages: invalid page parameter", "page", pageStr)
		return response.ErrorResponseJSON(c, fiber.StatusBadRequest, "invalid page parameter", fiber.StatusBadRequest)
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		logger.Warn("getAllMessages: invalid limit parameter", "limit", limitStr)
		return response.ErrorResponseJSON(c, fiber.StatusBadRequest, "invalid limit parameter", fiber.StatusBadRequest)
	}

	logger.Debug("getAllMessages: retrieving messages", "page", page, "limit", limit)
	messages, total, err := h.usecase.GetAllMessages(c.Context(), page, limit)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			logger.Warn("getAllMessages: failed to retrieve messages", "page", page, "error", appErr.Details)
			return response.ErrorResponseJSON(c, appErr.Code, appErr.Details, appErr.Code)
		}
		logger.Error("getAllMessages: unexpected error", "page", page, "error", err.Error())
		return response.ErrorResponseJSON(c, fiber.StatusInternalServerError, "failed to retrieve messages", fiber.StatusInternalServerError)
	}

	logger.Info("getAllMessages: messages retrieved successfully", "count", len(messages), "total", total)
	return response.PaginatedSuccessResponse(c, fiber.StatusOK, messages, page, limit, total, "Messages retrieved successfully")
}

// DeleteMessage deletes a contact submission (admin only)
// DELETE /api/messages/:id
func (h *MessageHandler) DeleteMessage(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		logger.Warn("deleteMessage: missing message ID")
		return response.ErrorResponseJSON(c, fiber.StatusBadRequest, "message ID is required", fiber.StatusBadRequest)
	}

	logger.Info("deleteMessage: attempting to delete message", "message_id", id)
	err := h.usecase.DeleteMessage(c.Context(), id)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			logger.Warn("deleteMessage: failed to delete message", "message_id", id, "error", appErr.Details)
			return response.ErrorResponseJSON(c, appErr.Code, appErr.Details, appErr.Code)
		}
		logger.Error("deleteMessage: unexpected error", "message_id", id, "error", err.Error())
		return response.ErrorResponseJSON(c, fiber.StatusInternalServerError, "failed to delete message", fiber.StatusInternalServerError)
	}

	logger.Info("deleteMessage: message deleted successfully", "message_id", id)
	return response.SuccessResponse(c, fiber.StatusOK, fiber.Map{}, "Message deleted successfully")
}
