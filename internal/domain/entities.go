package domain

import (
	"time"

	"github.com/google/uuid"
)

// User represents an admin user in the system
type User struct {
	ID        string    `db:"id"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// BlogPost represents a blog post
type BlogPost struct {
	ID        string    `db:"id"`
	Title     string    `db:"title"`
	Date      time.Time `db:"date"`
	Excerpt   string    `db:"excerpt"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// Project represents a portfolio project
type Project struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	Icon      string    `db:"icon"`
	Desc      string    `db:"desc"`
	Tech      string    `db:"tech"` // JSON string or comma-separated
	URL       string    `db:"url"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// ContactMessage represents a contact form submission
type ContactMessage struct {
	ID        string    `db:"id"`
	FromEmail string    `db:"from_email"`
	Subject   string    `db:"subject"`
	Message   string    `db:"message"`
	CreatedAt time.Time `db:"created_at"`
}

// Message is an alias for ContactMessage in the messages table
type Message struct {
	ID        string    `db:"id"`
	FromEmail string    `db:"from_email"`
	Subject   string    `db:"subject"`
	Message   string    `db:"message"`
	CreatedAt time.Time `db:"created_at"`
}

// DTOs for API requests/responses

// CreateBlogPostRequest DTO
type CreateBlogPostRequest struct {
	Title   string `json:"title" validate:"required,min=3,max=255"`
	Date    string `json:"date" validate:"required"`
	Excerpt string `json:"excerpt" validate:"required,min=10,max=500"`
	Content string `json:"content" validate:"required,min=10"`
}

// UpdateBlogPostRequest DTO
type UpdateBlogPostRequest struct {
	Title   string `json:"title" validate:"required,min=3,max=255"`
	Date    string `json:"date" validate:"required"`
	Excerpt string `json:"excerpt" validate:"required,min=10,max=500"`
	Content string `json:"content" validate:"required,min=10"`
}

// CreateProjectRequest DTO
type CreateProjectRequest struct {
	Name string `json:"name" validate:"required,min=3,max=255"`
	Icon string `json:"icon" validate:"required"`
	Desc string `json:"desc" validate:"required,min=10"`
	Tech string `json:"tech" validate:"required"`
	URL  string `json:"url" validate:"required,url"`
}

// UpdateProjectRequest DTO
type UpdateProjectRequest struct {
	Name string `json:"name" validate:"required,min=3,max=255"`
	Icon string `json:"icon" validate:"required"`
	Desc string `json:"desc" validate:"required,min=10"`
	Tech string `json:"tech" validate:"required"`
	URL  string `json:"url" validate:"required,url"`
}

// CreateContactRequest DTO
type CreateContactRequest struct {
	From    string `json:"from" validate:"required,email"`
	Subject string `json:"subject" validate:"required,min=3,max=200"`
	Message string `json:"message" validate:"required,min=10,max=5000"`
}

// LoginRequest DTO
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// LoginResponse DTO
type LoginResponse struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expires_in"`
	User      struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	} `json:"user"`
}

// RegisterRequest DTO
type RegisterRequest struct {
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

// NewUser creates a new user with ID and timestamp
func NewUser(email, passwordHash string) *User {
	now := time.Now()
	return &User{
		ID:        uuid.New().String(),
		Email:     email,
		Password:  passwordHash,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// NewBlogPost creates a new blog post with ID and timestamp
func NewBlogPost(title, excerpt, content string, date time.Time) *BlogPost {
	now := time.Now()
	return &BlogPost{
		ID:        uuid.New().String(),
		Title:     title,
		Date:      date,
		Excerpt:   excerpt,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// NewProject creates a new project with ID and timestamp
func NewProject(name, icon, desc, tech, url string) *Project {
	now := time.Now()
	return &Project{
		ID:        uuid.New().String(),
		Name:      name,
		Icon:      icon,
		Desc:      desc,
		Tech:      tech,
		URL:       url,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// NewContactMessage creates a new contact message
func NewContactMessage(fromEmail, subject, message string) *ContactMessage {
	return &ContactMessage{
		ID:        uuid.New().String(),
		FromEmail: fromEmail,
		Subject:   subject,
		Message:   message,
		CreatedAt: time.Now(),
	}
}
