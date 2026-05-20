# Development Guide - QuoteYourOS Backend

## 📚 Project Overview

This is a clean architecture REST API backend built with Go and PostgreSQL. The project has been initialized with all the necessary folder structure, configuration files, and basic domain layer.

## 🎯 What's Already Set Up

### ✅ Completed

1. **Project Structure** - Complete folder hierarchy following clean architecture
2. **Configuration** - Environment-based config using Viper
3. **Database** - PostgreSQL setup with Docker Compose
4. **Domain Layer** - Entities and interfaces defined
5. **Repository Layer** - Basic PostgreSQL implementations for:
   - Blog posts
   - Projects
   - Contact messages
   - Users
   - Messages
6. **Infrastructure** - Email and file service stubs
7. **Middleware** - JWT auth and CORS
8. **Docker** - Complete Docker and Docker Compose setup
9. **Documentation** - README and this guide

### ⚙️ Next Steps to Implement

1. **Use Cases** - Implement business logic layer
2. **HTTP Handlers** - Create request handlers for all endpoints
3. **Email Integration** - Complete email service implementation (Resend/SMTP)
4. **File Storage** - Implement PDF resume handling
5. **Tests** - Add comprehensive unit and integration tests
6. **Validation** - Add input validation middleware
7. **Error Handling** - Enhance error handling with proper status codes

## 🛠️ Implementation Guide

### 1. Implementing Use Cases

**Location**: `internal/usecase/<feature>/`

Example for Blog UseCase:

```go
// internal/usecase/blog/usecase.go
package blog

import (
	"context"
	"time"

	"github.com/quoteyouros/backend/internal/domain"
	"github.com/quoteyouros/backend/pkg/errors"
)

type BlogUseCase struct {
	repo domain.BlogRepository
}

func New(repo domain.BlogRepository) *BlogUseCase {
	return &BlogUseCase{repo: repo}
}

func (u *BlogUseCase) CreateBlogPost(ctx context.Context, req *domain.CreateBlogPostRequest) (*domain.BlogPost, error) {
	date, err := time.Parse(time.RFC3339, req.Date)
	if err != nil {
		return nil, errors.BadRequest("invalid date format")
	}

	post := domain.NewBlogPost(req.Title, req.Excerpt, req.Content, date)
	if err := u.repo.CreateBlogPost(ctx, post); err != nil {
		return nil, errors.InternalServerError(err.Error())
	}

	return post, nil
}

func (u *BlogUseCase) GetAllBlogPosts(ctx context.Context, page, limit int) ([]*domain.BlogPost, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit
	return u.repo.GetAllBlogPosts(ctx, limit, offset)
}

// Implement other methods...
}
```

### 2. Implementing HTTP Handlers

**Location**: `internal/handler/`

Example for Blog Handler:

```go
// internal/handler/blog.go
package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/quoteyouros/backend/internal/domain"
	"github.com/quoteyouros/backend/pkg/response"
)

type BlogHandler struct {
	usecase domain.BlogUseCase
}

func NewBlogHandler(usecase domain.BlogUseCase) *BlogHandler {
	return &BlogHandler{usecase: usecase}
}

func (h *BlogHandler) GetAllBlogPosts(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	posts, total, err := h.usecase.GetAllBlogPosts(c.Context(), page, limit)
	if err != nil {
		return response.ErrorResponseJSON(c, fiber.StatusInternalServerError, err.Error(), fiber.StatusInternalServerError)
	}

	return response.PaginatedSuccessResponse(c, fiber.StatusOK, posts, page, limit, total, "Blog posts retrieved successfully")
}

func (h *BlogHandler) GetBlogPost(c *fiber.Ctx) error {
	id := c.Params("id")
	post, err := h.usecase.GetBlogPost(c.Context(), id)
	if err != nil {
		return response.ErrorResponseJSON(c, fiber.StatusNotFound, "Blog post not found", fiber.StatusNotFound)
	}

	return response.SuccessResponse(c, fiber.StatusOK, post, "Blog post retrieved successfully")
}

// Implement other handlers...
}
```

### 3. Registering Routes

**Location**: `cmd/main.go`

```go
// Update setupPublicRoutes and setupProtectedRoutes
func setupPublicRoutes(app fiber.Router, handlers *AllHandlers) {
	blog := app.Group("/blog")
	blog.Get("", handlers.Blog.GetAllBlogPosts)
	blog.Get("/:id", handlers.Blog.GetBlogPost)
	
	// ... other routes
}

func setupProtectedRoutes(app fiber.Router, handlers *AllHandlers, authMiddleware fiber.Handler) {
	blog := app.Group("/blog", authMiddleware)
	blog.Post("", handlers.Blog.CreateBlogPost)
	blog.Put("/:id", handlers.Blog.UpdateBlogPost)
	blog.Delete("/:id", handlers.Blog.DeleteBlogPost)
	
	// ... other routes
}
```

### 4. Email Service Implementation

**For Resend**:
```go
// Send email via Resend API
func (s *ResendEmailService) SendEmail(to, subject, body string) error {
	// Use github.com/resend/resend-go
	// client := resend.NewClient(s.apiKey)
	// params := &resend.SendEmailRequest{...}
	// _, err := client.Emails.Send(params)
	// return err
}
```

**For SMTP**:
```go
// Send email via SMTP
func (s *SMTPEmailService) SendEmail(to, subject, body string) error {
	// Use github.com/go-mail/mail
	// m := mail.NewMessage()
	// m.SetHeader("From", s.fromEmail)
	// ... setup message
	// return dialer.DialAndSend(m)
}
```

### 5. Testing

Create test files alongside implementations:

```go
// internal/usecase/blog/usecase_test.go
package blog

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/quoteyouros/backend/internal/domain"
)

func TestCreateBlogPost(t *testing.T) {
	// Setup mock repository
	// Execute use case
	// Assert results
}
```

## 🚀 Quick Start Development

1. **Install dependencies**:
   ```bash
   go mod tidy
   ```

2. **Start database**:
   ```bash
   make db-up
   ```

3. **Run migrations**:
   ```bash
   make db-migrate
   ```

4. **Run application**:
   ```bash
   make run
   ```

5. **Test the API**:
   ```bash
   curl http://localhost:8000/health
   ```

## 📦 Dependencies to Add

When implementing features, add these as needed:

```bash
# Email
go get github.com/resend/resend-go
go get github.com/go-mail/mail/v2

# Validation
go get github.com/go-playground/validator/v10

# PDF handling
go get github.com/jung-kurt/gofpdf

# Testing
go get github.com/stretchr/testify

# Logging (already in go.mod)
go get go.uber.org/zap
```

## 🔒 Authentication Flow

1. User calls `POST /api/auth/login` with email and password
2. Backend validates credentials against database
3. Backend generates JWT token (using JWT secret from config)
4. Token returned to client
5. Client includes token in `Authorization: Bearer <token>` header
6. Middleware validates token on protected routes

## 📝 Key Files to Know

| File | Purpose |
|------|---------|
| `cmd/main.go` | Application entry point and routing |
| `internal/config/config.go` | Configuration management |
| `internal/domain/entities.go` | Domain entities and DTOs |
| `internal/domain/repositories.go` | Repository interfaces |
| `internal/domain/usecases.go` | UseCase interfaces |
| `internal/middleware/auth.go` | JWT authentication |
| `migrations/001_init_schema.sql` | Database schema |

## 🐛 Debugging Tips

- Use `go run ./cmd` to run directly (better error messages)
- Check database connection: `docker logs quoteyouros_db`
- View app logs: `docker logs quoteyouros_backend`
- Test endpoints with curl or Postman
- Enable debug logging in config

## 📋 Architecture Checklist

- [ ] All use cases implemented
- [ ] All handlers created
- [ ] Routes registered
- [ ] Email service working
- [ ] File upload working
- [ ] JWT authentication working
- [ ] Input validation added
- [ ] Error handling complete
- [ ] Tests written
- [ ] Docker builds successfully
- [ ] Environment variables documented

---

**Happy coding! 🚀**
