package message_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/quoteyouros/backend/internal/domain"
	"github.com/quoteyouros/backend/internal/usecase/message"
	apperrors "github.com/quoteyouros/backend/pkg/errors"
)

type mockMessageRepository struct {
	getAllFunc func(ctx context.Context, limit, offset int) ([]*domain.Message, int64, error)
	deleteFunc func(ctx context.Context, id string) error
}

func (m *mockMessageRepository) GetAllMessages(ctx context.Context, limit, offset int) ([]*domain.Message, int64, error) {
	if m.getAllFunc != nil {
		return m.getAllFunc(ctx, limit, offset)
	}
	return nil, 0, nil
}

func (m *mockMessageRepository) DeleteMessage(ctx context.Context, id string) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return nil
}

func TestGetAllMessages_Success(t *testing.T) {
	mockMessages := []*domain.Message{
		{ID: "1", FromEmail: "a@b.com", Subject: "Hi", Message: "Hello", CreatedAt: time.Now()},
		{ID: "2", FromEmail: "c@d.com", Subject: "Bye", Message: "Goodbye", CreatedAt: time.Now()},
	}

	var passedLimit, passedOffset int
	repo := &mockMessageRepository{
		getAllFunc: func(ctx context.Context, limit, offset int) ([]*domain.Message, int64, error) {
			passedLimit = limit
			passedOffset = offset
			return mockMessages, 2, nil
		},
	}

	uc := message.New(repo)

	msgs, total, err := uc.GetAllMessages(context.Background(), 2, 20)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if total != 2 {
		t.Errorf("expected total 2, got %d", total)
	}

	if len(msgs) != 2 {
		t.Errorf("expected 2 messages, got %d", len(msgs))
	}

	// offset = (page - 1) * limit = (2 - 1) * 20 = 20
	if passedLimit != 20 {
		t.Errorf("expected limit 20, got %d", passedLimit)
	}

	if passedOffset != 20 {
		t.Errorf("expected offset 20, got %d", passedOffset)
	}
}

func TestGetAllMessages_PaginationDefaults(t *testing.T) {
	repo := &mockMessageRepository{
		getAllFunc: func(ctx context.Context, limit, offset int) ([]*domain.Message, int64, error) {
			if limit != 10 {
				t.Errorf("expected default limit 10, got %d", limit)
			}
			if offset != 0 {
				t.Errorf("expected offset 0, got %d", offset)
			}
			return nil, 0, nil
		},
	}

	uc := message.New(repo)

	// page 0 and limit 0 should trigger defaults (page=1, limit=10)
	_, _, err := uc.GetAllMessages(context.Background(), 0, 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestGetAllMessages_RepoError(t *testing.T) {
	repo := &mockMessageRepository{
		getAllFunc: func(ctx context.Context, limit, offset int) ([]*domain.Message, int64, error) {
			return nil, 0, errors.New("db query failed")
		},
	}

	uc := message.New(repo)

	_, _, err := uc.GetAllMessages(context.Background(), 1, 10)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	appErr, ok := err.(*apperrors.AppError)
	if !ok {
		t.Fatalf("expected AppError, got %T: %v", err, err)
	}

	if !strings.Contains(appErr.Details, "failed to retrieve messages") {
		t.Errorf("expected retrieve messages error, got details: %q", appErr.Details)
	}
}

func TestDeleteMessage_Success(t *testing.T) {
	var deletedID string
	repo := &mockMessageRepository{
		deleteFunc: func(ctx context.Context, id string) error {
			deletedID = id
			return nil
		},
	}

	uc := message.New(repo)

	err := uc.DeleteMessage(context.Background(), "test-id")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if deletedID != "test-id" {
		t.Errorf("expected to delete test-id, got %q", deletedID)
	}
}

func TestDeleteMessage_ValidationError(t *testing.T) {
	repo := &mockMessageRepository{}
	uc := message.New(repo)

	err := uc.DeleteMessage(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty ID, got nil")
	}

	appErr, ok := err.(*apperrors.AppError)
	if !ok {
		t.Fatalf("expected AppError, got %T: %v", err, err)
	}

	if !strings.Contains(appErr.Details, "message ID is required") {
		t.Errorf("expected empty ID error message, got %q", appErr.Details)
	}
}

func TestDeleteMessage_RepoError(t *testing.T) {
	repo := &mockMessageRepository{
		deleteFunc: func(ctx context.Context, id string) error {
			return errors.New("db delete failed")
		},
	}

	uc := message.New(repo)

	err := uc.DeleteMessage(context.Background(), "test-id")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	appErr, ok := err.(*apperrors.AppError)
	if !ok {
		t.Fatalf("expected AppError, got %T: %v", err, err)
	}

	if !strings.Contains(appErr.Details, "failed to delete message") {
		t.Errorf("expected delete failure error, got details: %q", appErr.Details)
	}
}
