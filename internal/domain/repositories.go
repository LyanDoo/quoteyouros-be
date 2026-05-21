package domain

import "context"

// Repository interfaces follow the Interface Segregation Principle

// BlogRepository defines blog data access operations
type BlogRepository interface {
	CreateBlogPost(ctx context.Context, post *BlogPost) error
	GetBlogPost(ctx context.Context, id string) (*BlogPost, error)
	GetAllBlogPosts(ctx context.Context, limit, offset int) ([]*BlogPost, int64, error)
	UpdateBlogPost(ctx context.Context, id string, post *BlogPost) error
	DeleteBlogPost(ctx context.Context, id string) error
}

// ProjectRepository defines project data access operations
type ProjectRepository interface {
	CreateProject(ctx context.Context, project *Project) error
	GetProject(ctx context.Context, id string) (*Project, error)
	GetAllProjects(ctx context.Context) ([]*Project, error)
	UpdateProject(ctx context.Context, id string, project *Project) error
	DeleteProject(ctx context.Context, id string) error
}

// ContactRepository defines contact message data access operations
type ContactRepository interface {
	CreateContactMessage(ctx context.Context, message *ContactMessage) error
}

// MessageRepository defines message retrieval operations
type MessageRepository interface {
	GetAllMessages(ctx context.Context, limit, offset int) ([]*Message, int64, error)
	DeleteMessage(ctx context.Context, id string) error
}

// UserRepository defines user data access operations
type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id string) (*User, error)
}

// ProfileRepository defines profile data access operations
type ProfileRepository interface {
	GetProfile(ctx context.Context) (*Profile, error)
	CreateProfile(ctx context.Context, profile *Profile) error
	UpdateProfile(ctx context.Context, profile *Profile) error
	SaveResume(ctx context.Context, profile *Profile) error
	DeleteResume(ctx context.Context) error
}

// CommentRepository defines comment data access operations
type CommentRepository interface {
	CreateComment(ctx context.Context, comment *Comment) error
	GetComment(ctx context.Context, id string) (*Comment, error)
	GetCommentsByBlogPost(ctx context.Context, blogPostID string, limit, offset int) ([]*Comment, int64, error)
	GetReplies(ctx context.Context, commentID string, limit, offset int) ([]*Comment, int64, error)
	UpdateComment(ctx context.Context, id string, comment *Comment) error
	DeleteComment(ctx context.Context, id string) error
}
