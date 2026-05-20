package auth

import (
	"context"
	"regexp"
	"time"

	"github.com/quoteyouros/backend/internal/domain"
	"github.com/quoteyouros/backend/pkg/crypto"
	apperrors "github.com/quoteyouros/backend/pkg/errors"
	jwtpkg "github.com/quoteyouros/backend/pkg/jwt"
	"github.com/quoteyouros/backend/pkg/logger"
)

type AuthUseCase struct {
	userRepo      domain.UserRepository
	jwtSecret     string
	jwtExpiration time.Duration
}

// New creates a new auth use case
func New(userRepo domain.UserRepository, jwtSecret string, jwtExpiration time.Duration) *AuthUseCase {
	return &AuthUseCase{
		userRepo:      userRepo,
		jwtSecret:     jwtSecret,
		jwtExpiration: jwtExpiration,
	}
}

// Register creates a new user account
func (u *AuthUseCase) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.User, error) {
	logger.Debug("register: validating request", "email", req.Email)

	// Validate email format
	if !isValidEmail(req.Email) {
		logger.Warn("register: invalid email format", "email", req.Email)
		return nil, apperrors.BadRequest("invalid email format")
	}

	// Validate password strength (minimum 8 characters)
	if len(req.Password) < 8 {
		logger.Warn("register: password too short", "email", req.Email)
		return nil, apperrors.BadRequest("password must be at least 8 characters")
	}

	// Verify passwords match
	if req.Password != req.ConfirmPassword {
		logger.Warn("register: passwords do not match", "email", req.Email)
		return nil, apperrors.BadRequest("passwords do not match")
	}

	// Check if email already exists
	logger.Debug("register: checking if email exists", "email", req.Email)
	existingUser, _ := u.userRepo.GetUserByEmail(ctx, req.Email)
	if existingUser != nil {
		logger.Warn("register: email already registered", "email", req.Email)
		return nil, apperrors.Conflict("email already registered")
	}

	// Hash password
	logger.Debug("register: hashing password", "email", req.Email)
	hashedPassword, err := crypto.HashPassword(req.Password)
	if err != nil {
		logger.Error("register: failed to hash password", "email", req.Email, "error", err.Error())
		return nil, apperrors.InternalServerError("failed to process password")
	}

	// Create new user
	newUser := domain.NewUser(req.Email, hashedPassword)

	// Save to database
	logger.Info("register: creating user in database", "email", req.Email, "user_id", newUser.ID)
	if err := u.userRepo.CreateUser(ctx, newUser); err != nil {
		logger.Error("register: failed to create user", "email", req.Email, "error", err.Error())
		return nil, apperrors.InternalServerError("failed to create user: " + err.Error())
	}

	logger.Info("register: user created successfully", "email", req.Email, "user_id", newUser.ID)
	return newUser, nil
}

// Login authenticates user and returns JWT token
func (u *AuthUseCase) Login(ctx context.Context, req *domain.LoginRequest) (string, int64, error) {
	logger.Debug("login: attempting authentication", "email", req.Email)

	// Find user by email
	user, err := u.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil || user == nil {
		logger.Warn("login: user not found", "email", req.Email)
		return "", 0, apperrors.Unauthorized("invalid email or password")
	}

	// Verify password
	if !crypto.VerifyPassword(user.Password, req.Password) {
		logger.Warn("login: invalid password", "email", req.Email)
		return "", 0, apperrors.Unauthorized("invalid email or password")
	}

	// Generate JWT token
	logger.Debug("login: generating JWT token", "email", req.Email, "user_id", user.ID)
	token, err := jwtpkg.GenerateToken(user.ID, user.Email, u.jwtSecret, u.jwtExpiration)
	if err != nil {
		logger.Error("login: failed to generate token", "email", req.Email, "error", err.Error())
		return "", 0, apperrors.InternalServerError("failed to generate token")
	}

	// Return token and expiration time in seconds
	expiresIn := int64(u.jwtExpiration.Seconds())

	logger.Info("login: authentication successful", "email", req.Email, "user_id", user.ID)
	return token, expiresIn, nil
}

// GetUserByID retrieves a user by ID (for protected routes)
func (u *AuthUseCase) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	logger.Debug("getUserByID: retrieving user", "user_id", id)

	user, err := u.userRepo.GetUserByID(ctx, id)
	if err != nil {
		logger.Warn("getUserByID: user not found", "user_id", id, "error", err.Error())
		return nil, apperrors.NotFound("user not found")
	}

	// Don't return password
	user.Password = ""
	return user, nil
}

// ValidateToken validates a JWT token and returns claims
func (u *AuthUseCase) ValidateToken(token string) (*jwtpkg.Claims, error) {
	claims, err := jwtpkg.ValidateToken(token, u.jwtSecret)
	if err != nil {
		return nil, apperrors.Unauthorized("invalid or expired token")
	}
	return claims, nil
}

// isValidEmail checks if email format is valid
func isValidEmail(email string) bool {
	// Simple email regex pattern
	const emailPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailPattern)
	return re.MatchString(email)
}
