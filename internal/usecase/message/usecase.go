package message

import (
	"context"

	"github.com/quoteyouros/backend/internal/domain"
	apperrors "github.com/quoteyouros/backend/pkg/errors"
	"github.com/quoteyouros/backend/pkg/logger"
)

type MessageUseCase struct {
	messageRepo domain.MessageRepository
}

// New creates a new message use case
func New(messageRepo domain.MessageRepository) *MessageUseCase {
	return &MessageUseCase{
		messageRepo: messageRepo,
	}
}

// GetAllMessages retrieves messages with pagination
func (u *MessageUseCase) GetAllMessages(ctx context.Context, page, limit int) ([]*domain.Message, int64, error) {
	logger.Debug("getAllMessages: retrieving messages", "page", page, "limit", limit)

	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	messages, total, err := u.messageRepo.GetAllMessages(ctx, limit, offset)
	if err != nil {
		logger.Error("getAllMessages: failed to retrieve messages", "page", page, "limit", limit, "error", err.Error())
		return nil, 0, apperrors.InternalServerError("failed to retrieve messages: " + err.Error())
	}

	logger.Info("getAllMessages: messages retrieved successfully", "count", len(messages), "total", total, "page", page)
	return messages, total, nil
}

// DeleteMessage deletes a message by ID
func (u *MessageUseCase) DeleteMessage(ctx context.Context, id string) error {
	logger.Debug("deleteMessage: deleting message", "message_id", id)

	if id == "" {
		logger.Warn("deleteMessage: message ID is empty")
		return apperrors.BadRequest("message ID is required")
	}

	if err := u.messageRepo.DeleteMessage(ctx, id); err != nil {
		logger.Error("deleteMessage: failed to delete message", "message_id", id, "error", err.Error())
		return apperrors.InternalServerError("failed to delete message: " + err.Error())
	}

	logger.Info("deleteMessage: message deleted successfully", "message_id", id)
	return nil
}
