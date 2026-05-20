package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/quoteyouros/backend/internal/domain"
	apperrors "github.com/quoteyouros/backend/pkg/errors"
	"github.com/quoteyouros/backend/pkg/logger"
	"github.com/quoteyouros/backend/pkg/response"
)

type AuthHandler struct {
	usecase domain.AuthUseCase
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(usecase domain.AuthUseCase) *AuthHandler {
	return &AuthHandler{usecase: usecase}
}

// Register creates a new user account
// POST /api/auth/register
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req domain.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Error("register: failed to parse request body", "error", err.Error())
		return response.ErrorResponseJSON(c, fiber.StatusBadRequest, "invalid request body", fiber.StatusBadRequest)
	}

	logger.Info("register: attempting to register user", "email", req.Email)
	user, err := h.usecase.Register(c.Context(), &req)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			logger.Warn("register: registration failed", "email", req.Email, "error", appErr.Details)
			return response.ErrorResponseJSON(c, appErr.Code, appErr.Details, appErr.Code)
		}
		logger.Error("register: unexpected error", "email", req.Email, "error", err.Error())
		return response.ErrorResponseJSON(c, fiber.StatusInternalServerError, "registration failed", fiber.StatusInternalServerError)
	}

	logger.Info("register: user registered successfully", "user_id", user.ID, "email", user.Email)
	return response.SuccessResponse(c, fiber.StatusCreated, fiber.Map{
		"id":    user.ID,
		"email": user.Email,
	}, "User registered successfully")
}

// Login authenticates user and returns JWT token
// POST /api/auth/login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req domain.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Error("login: failed to parse request body", "error", err.Error())
		return response.ErrorResponseJSON(c, fiber.StatusBadRequest, "invalid request body", fiber.StatusBadRequest)
	}

	logger.Info("login: attempting login", "email", req.Email)
	token, expiresIn, err := h.usecase.Login(c.Context(), &req)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			logger.Warn("login: authentication failed", "email", req.Email, "error", appErr.Details)
			return response.ErrorResponseJSON(c, appErr.Code, appErr.Details, appErr.Code)
		}
		logger.Error("login: unexpected error", "email", req.Email, "error", err.Error())
		return response.ErrorResponseJSON(c, fiber.StatusInternalServerError, "authentication failed", fiber.StatusInternalServerError)
	}

	logger.Info("login: successful login", "email", req.Email)
	return response.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"token":      token,
		"expires_in": expiresIn,
	}, "Login successful")
}

// GetCurrentUser returns the current authenticated user
// GET /api/auth/me
func (h *AuthHandler) GetCurrentUser(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		logger.Warn("me: unauthorized access attempt")
		return response.ErrorResponseJSON(c, fiber.StatusUnauthorized, "unauthorized", fiber.StatusUnauthorized)
	}

	logger.Debug("me: retrieving user", "user_id", userID)
	user, err := h.usecase.GetUserByID(c.Context(), userID)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			logger.Warn("me: failed to retrieve user", "user_id", userID, "error", appErr.Details)
			return response.ErrorResponseJSON(c, appErr.Code, appErr.Details, appErr.Code)
		}
		logger.Error("me: unexpected error", "user_id", userID, "error", err.Error())
		return response.ErrorResponseJSON(c, fiber.StatusInternalServerError, "failed to fetch user", fiber.StatusInternalServerError)
	}

	logger.Debug("me: user retrieved successfully", "user_id", userID)
	return response.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"id":         user.ID,
		"email":      user.Email,
		"created_at": user.CreatedAt,
	}, "User retrieved successfully")
}

// Logout invalidates the current session
// POST /api/auth/logout
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// JWT tokens are stateless, so logout is just a client-side operation
	// Client should delete the token from storage
	return response.SuccessResponse(c, fiber.StatusOK, fiber.Map{}, "Logged out successfully")
}
