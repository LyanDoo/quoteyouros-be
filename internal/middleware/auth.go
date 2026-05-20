package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/quoteyouros/backend/pkg/jwt"
	"github.com/quoteyouros/backend/pkg/logger"
	"github.com/quoteyouros/backend/pkg/response"
)

// JWTAuth validates JWT tokens using the jwt package
func JWTAuth(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			logger.Warn("JWTAuth: missing authorization header", "path", c.Path())
			return response.ErrorResponseJSON(c, fiber.StatusUnauthorized, "missing authorization header", fiber.StatusUnauthorized)
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			logger.Warn("JWTAuth: invalid authorization format", "path", c.Path())
			return response.ErrorResponseJSON(c, fiber.StatusUnauthorized, "invalid authorization format", fiber.StatusUnauthorized)
		}

		tokenString := parts[1]

		// Validate token using jwt package
		claims, err := jwt.ValidateToken(tokenString, secret)
		if err != nil {
			logger.Warn("JWTAuth: invalid or expired token", "path", c.Path(), "error", err.Error())
			return response.ErrorResponseJSON(c, fiber.StatusUnauthorized, "invalid or expired token", fiber.StatusUnauthorized)
		}

		// Store user info in context
		c.Locals("user_id", claims.UserID)
		c.Locals("email", claims.Email)

		logger.Debug("JWTAuth: token validated successfully", "user_id", claims.UserID, "path", c.Path())

		return c.Next()
	}
}
