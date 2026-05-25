package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/quoteyouros/backend/internal/domain"
	contactusecase "github.com/quoteyouros/backend/internal/usecase/contact"
	apperrors "github.com/quoteyouros/backend/pkg/errors"
	"github.com/quoteyouros/backend/pkg/logger"
	"github.com/quoteyouros/backend/pkg/response"
)

type ContactHandler struct {
	usecase *contactusecase.ContactUseCase
}

// NewContactHandler creates a new contact handler
func NewContactHandler(usecase *contactusecase.ContactUseCase) *ContactHandler {
	return &ContactHandler{usecase: usecase}
}

// SubmitContact receives a contact form submission and forwards it to use case
// POST /api/contact
func (h *ContactHandler) SubmitContact(c *fiber.Ctx) error {
	var req domain.CreateContactRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Error("submitContact: failed to parse request body", "error", err.Error())
		return response.ErrorResponseJSON(c, fiber.StatusBadRequest, "invalid request body", fiber.StatusBadRequest)
	}

	logger.Info("submitContact: attempting contact submission", "from", req.From)
	err := h.usecase.SubmitContact(c.Context(), &req)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			logger.Warn("submitContact: submission failed", "from", req.From, "error", appErr.Details)
			return response.ErrorResponseJSON(c, appErr.Code, appErr.Details, appErr.Code)
		}
		logger.Error("submitContact: unexpected error", "from", req.From, "error", err.Error())
		return response.ErrorResponseJSON(c, fiber.StatusInternalServerError, "failed to submit contact message", fiber.StatusInternalServerError)
	}

	logger.Info("submitContact: contact submission successful", "from", req.From)
	return response.SuccessResponse(c, fiber.StatusCreated, fiber.Map{}, "Contact message sent successfully")
}
