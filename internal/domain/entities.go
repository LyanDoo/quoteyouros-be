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

// Profile represents user profile information including resume
type Profile struct {
	ID               string     `db:"id"`
	About            string     `db:"about"`
	ResumeFileName   *string    `db:"resume_file_name"`
	ResumeFileSize   *int64     `db:"resume_file_size"`
	ResumeFilePath   *string    `db:"resume_file_path"`
	ResumeMimeType   *string    `db:"resume_mime_type"`
	ResumeUploadedAt *time.Time `db:"resume_uploaded_at"`
	CreatedAt        time.Time  `db:"created_at"`
	UpdatedAt        time.Time  `db:"updated_at"`
}

// Comment represents a blog post comment
type Comment struct {
	ID               string    `db:"id"`
	BlogPostID       string    `db:"blog_post_id"`
	ReplyToCommentID *string   `db:"reply_to_comment_id"`
	AuthorName       string    `db:"author_name"`
	AuthorEmail      *string   `db:"author_email"`
	Content          string    `db:"content"`
	Rating           *int      `db:"rating"`
	IsSpam           bool      `db:"is_spam"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}

// DTOs for API requests/responses

// CreateBlogPostRequest DTO
type CreateBlogPostRequest struct {
	Title   string `json:"title" validate:"required,min=3,max=255"`
	Date    string `json:"date"`
	Excerpt string `json:"excerpt" validate:"required,min=10,max=500"`
	Content string `json:"content" validate:"required,min=10"`
}

// UpdateBlogPostRequest DTO
type UpdateBlogPostRequest struct {
	Title   string `json:"title" validate:"required,min=3,max=255"`
	Date    string `json:"date"`
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

// UpdateProfileRequest DTO
type UpdateProfileRequest struct {
	About string `json:"about"`
}

// ProfileResponse DTO for API responses
type ProfileResponse struct {
	About  string `json:"about"`
	Resume struct {
		FileName   string `json:"file_name"`
		FileSize   int64  `json:"file_size"`
		UploadedAt string `json:"uploaded_at"`
		HasResume  bool   `json:"has_resume"`
	} `json:"resume"`
}

// ResumeUploadResponse DTO for successful upload
type ResumeUploadResponse struct {
	FileName    string `json:"file_name"`
	FileSize    int64  `json:"file_size"`
	UploadedAt  string `json:"uploaded_at"`
	DownloadURL string `json:"download_url"`
}

// CreateCommentRequest DTO for creating a comment
type CreateCommentRequest struct {
	AuthorName       string  `json:"author_name" validate:"required,min=1,max=255"`
	AuthorEmail      *string `json:"author_email" validate:"omitempty,email"`
	Content          string  `json:"content" validate:"required,min=1,max=5000"`
	Rating           *int    `json:"rating" validate:"omitempty,gte=1,lte=5"`
	ReplyToCommentID *string `json:"reply_to_comment_id"`
}

// UpdateCommentRequest DTO for updating a comment
type UpdateCommentRequest struct {
	Content string `json:"content" validate:"required,min=1,max=5000"`
	Rating  *int   `json:"rating" validate:"omitempty,gte=1,lte=5"`
	IsSpam  *bool  `json:"is_spam"`
}

// CommentResponse DTO for API responses
type CommentResponse struct {
	ID               string  `json:"id"`
	BlogPostID       string  `json:"blog_post_id"`
	ReplyToCommentID *string `json:"reply_to_comment_id,omitempty"`
	AuthorName       string  `json:"author_name"`
	AuthorEmail      *string `json:"author_email,omitempty"`
	Content          string  `json:"content"`
	Rating           *int    `json:"rating,omitempty"`
	IsSpam           bool    `json:"is_spam"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
}

// CommentListResponse DTO for paginated comments
type CommentListResponse struct {
	Comments   []CommentResponse `json:"comments"`
	Total      int               `json:"total"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	TotalPages int               `json:"total_pages"`
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

// NewComment creates a new comment with ID and timestamp
func NewComment(blogPostID, authorName string, authorEmail *string, content string, replyToCommentID *string, rating *int) *Comment {
	now := time.Now()
	return &Comment{
		ID:               uuid.New().String(),
		BlogPostID:       blogPostID,
		AuthorName:       authorName,
		AuthorEmail:      authorEmail,
		Content:          content,
		ReplyToCommentID: replyToCommentID,
		Rating:           rating,
		IsSpam:           false,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}
