package contact_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/quoteyouros/backend/internal/domain"
	"github.com/quoteyouros/backend/internal/usecase/contact"
	apperrors "github.com/quoteyouros/backend/pkg/errors"
)

type mockContactRepository struct {
	createFunc func(ctx context.Context, message *domain.ContactMessage) error
}

func (m *mockContactRepository) CreateContactMessage(ctx context.Context, message *domain.ContactMessage) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, message)
	}
	return nil
}

type mockEmailService struct {
	sendFunc func(to, subject, body string) error
}

func (m *mockEmailService) SendEmail(to, subject, body string) error {
	if m.sendFunc != nil {
		return m.sendFunc(to, subject, body)
	}
	return nil
}

func TestSubmitContact_Success(t *testing.T) {
	var savedMessage *domain.ContactMessage
	repo := &mockContactRepository{
		createFunc: func(ctx context.Context, message *domain.ContactMessage) error {
			savedMessage = message
			return nil
		},
	}

	emailChan := make(chan bool, 1)
	emailSvc := &mockEmailService{
		sendFunc: func(to, subject, body string) error {
			if to != "admin@example.com" {
				t.Errorf("expected to email to admin@example.com, got %s", to)
			}
			emailChan <- true
			return nil
		},
	}

	uc := contact.New(repo, emailSvc, "admin@example.com")

	req := &domain.CreateContactRequest{
		From:    "user@example.com",
		Subject: "Inquiry about QuoteYourOS",
		Message: "Hello, I am interested in your project and would like to connect.",
	}

	err := uc.SubmitContact(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Wait for email goroutine to execute
	<-emailChan

	if savedMessage == nil {
		t.Fatal("expected message to be saved in repository")
	}

	if savedMessage.FromEmail != req.From {
		t.Errorf("expected from %s, got %s", req.From, savedMessage.FromEmail)
	}

	if savedMessage.Subject != req.Subject {
		t.Errorf("expected subject %s, got %s", req.Subject, savedMessage.Subject)
	}

	if savedMessage.Message != req.Message {
		t.Errorf("expected message %s, got %s", req.Message, savedMessage.Message)
	}
}

func TestSubmitContact_ValidationFailures(t *testing.T) {
	repo := &mockContactRepository{}
	emailSvc := &mockEmailService{}
	uc := contact.New(repo, emailSvc, "admin@example.com")

	tests := []struct {
		name    string
		req     *domain.CreateContactRequest
		wantErr string
	}{
		{
			name: "invalid email",
			req: &domain.CreateContactRequest{
				From:    "invalid-email",
				Subject: "Valid Subject",
				Message: "Valid message content of correct length.",
			},
			wantErr: "invalid email address",
		},
		{
			name: "subject too short",
			req: &domain.CreateContactRequest{
				From:    "user@example.com",
				Subject: "Hi",
				Message: "Valid message content of correct length.",
			},
			wantErr: "subject must be between 3 and 200 characters",
		},
		{
			name: "subject too long",
			req: &domain.CreateContactRequest{
				From:    "user@example.com",
				Subject: strings.Repeat("s", 201),
				Message: "Valid message content of correct length.",
			},
			wantErr: "subject must be between 3 and 200 characters",
		},
		{
			name: "message too short",
			req: &domain.CreateContactRequest{
				From:    "user@example.com",
				Subject: "Valid Subject",
				Message: "Short",
			},
			wantErr: "message must be between 10 and 5000 characters",
		},
		{
			name: "message too long",
			req: &domain.CreateContactRequest{
				From:    "user@example.com",
				Subject: "Valid Subject",
				Message: strings.Repeat("m", 5001),
			},
			wantErr: "message must be between 10 and 5000 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := uc.SubmitContact(context.Background(), tt.req)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			appErr, ok := err.(*apperrors.AppError)
			if !ok {
				t.Fatalf("expected AppError, got %T: %v", err, err)
			}
			if !strings.Contains(appErr.Details, tt.wantErr) {
				t.Errorf("expected error details containing %q, got %q", tt.wantErr, appErr.Details)
			}
		})
	}
}

func TestSubmitContact_RepoError(t *testing.T) {
	repo := &mockContactRepository{
		createFunc: func(ctx context.Context, message *domain.ContactMessage) error {
			return errors.New("db error")
		},
	}
	emailSvc := &mockEmailService{}
	uc := contact.New(repo, emailSvc, "admin@example.com")

	req := &domain.CreateContactRequest{
		From:    "user@example.com",
		Subject: "Valid Subject",
		Message: "Valid message content of correct length.",
	}

	err := uc.SubmitContact(context.Background(), req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	appErr, ok := err.(*apperrors.AppError)
	if !ok {
		t.Fatalf("expected AppError, got %T: %v", err, err)
	}
	if !strings.Contains(appErr.Details, "failed to store contact message") {
		t.Errorf("expected database storage failure error details containing 'failed to store contact message', got %q", appErr.Details)
	}
}
