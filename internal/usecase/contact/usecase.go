package contact

import (
	"context"
	"fmt"
	"regexp"

	"github.com/quoteyouros/backend/internal/domain"
	"github.com/quoteyouros/backend/internal/infrastructure/email"
	apperrors "github.com/quoteyouros/backend/pkg/errors"
	"github.com/quoteyouros/backend/pkg/logger"
)

type ContactUseCase struct {
	contactRepo  domain.ContactRepository
	emailService email.EmailService
	adminEmail   string
}

// New creates a new contact use case
func New(contactRepo domain.ContactRepository, emailService email.EmailService, adminEmail string) *ContactUseCase {
	return &ContactUseCase{
		contactRepo:  contactRepo,
		emailService: emailService,
		adminEmail:   adminEmail,
	}
}

// SubmitContact validates request, stores in DB, and optionally sends notification email
func (u *ContactUseCase) SubmitContact(ctx context.Context, req *domain.CreateContactRequest) error {
	logger.Debug("submitContact: validating request", "from", req.From, "subject", req.Subject)

	// Validate email format
	if !isValidEmail(req.From) {
		logger.Warn("submitContact: invalid email format", "from", req.From)
		return apperrors.BadRequest("invalid email address")
	}

	// Validate subject length
	if len(req.Subject) < 3 || len(req.Subject) > 200 {
		logger.Warn("submitContact: invalid subject length", "len", len(req.Subject))
		return apperrors.BadRequest("subject must be between 3 and 200 characters")
	}

	// Validate message length
	if len(req.Message) < 10 || len(req.Message) > 5000 {
		logger.Warn("submitContact: invalid message length", "len", len(req.Message))
		return apperrors.BadRequest("message must be between 10 and 5000 characters")
	}

	// Create contact message entity
	message := domain.NewContactMessage(req.From, req.Subject, req.Message)

	logger.Info("submitContact: storing message in database", "message_id", message.ID, "from", message.FromEmail)
	if err := u.contactRepo.CreateContactMessage(ctx, message); err != nil {
		logger.Error("submitContact: failed to store message", "message_id", message.ID, "error", err.Error())
		return apperrors.InternalServerError("failed to store contact message: " + err.Error())
	}

	logger.Info("submitContact: contact message saved successfully", "message_id", message.ID)

	// Asynchronously trigger email notification if configured
	if u.emailService != nil && u.adminEmail != "" {
		go func() {
			subject := fmt.Sprintf("[QuoteYourOS Contact] %s", message.Subject)
			body := fmt.Sprintf(
				"You have received a new contact message from Outlook Express form on QuoteYourOS.\n\nFrom: %s\nSubject: %s\nReceived: %s\n\nMessage:\n%s\n",
				message.FromEmail, message.Subject, message.CreatedAt.Format("2006-01-02 15:04:05 MST"), message.Message,
			)

			logger.Debug("submitContact: sending email notification asynchronously", "to", u.adminEmail)
			if err := u.emailService.SendEmail(u.adminEmail, subject, body); err != nil {
				logger.Error("submitContact: failed to send email notification", "error", err.Error())
			} else {
				logger.Info("submitContact: email notification sent successfully", "to", u.adminEmail)
			}
		}()
	}

	return nil
}

// isValidEmail checks if email format is valid
func isValidEmail(email string) bool {
	const emailPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailPattern)
	return re.MatchString(email)
}
