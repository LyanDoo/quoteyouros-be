package domain

import (
	"context"
	"mime/multipart"
)

// BlogUseCase defines blog business logic operations
type BlogUseCase interface {
	CreateBlogPost(ctx context.Context, req *CreateBlogPostRequest) (*BlogPost, error)
	GetBlogPost(ctx context.Context, id string) (*BlogPost, error)
	GetAllBlogPosts(ctx context.Context, page, limit int) ([]*BlogPost, int64, error)
	UpdateBlogPost(ctx context.Context, id string, req *UpdateBlogPostRequest) (*BlogPost, error)
	DeleteBlogPost(ctx context.Context, id string) error
}

// ProjectUseCase defines project business logic operations
type ProjectUseCase interface {
	CreateProject(ctx context.Context, req *CreateProjectRequest) (*Project, error)
	GetProject(ctx context.Context, id string) (*Project, error)
	GetAllProjects(ctx context.Context) ([]*Project, error)
	UpdateProject(ctx context.Context, id string, req *UpdateProjectRequest) (*Project, error)
	DeleteProject(ctx context.Context, id string) error
}

// ContactUseCase defines contact business logic operations
type ContactUseCase interface {
	SubmitContact(ctx context.Context, req *CreateContactRequest) error
}

// AuthUseCase defines authentication business logic operations
type AuthUseCase interface {
	Register(ctx context.Context, req *RegisterRequest) (*User, error)
	Login(ctx context.Context, req *LoginRequest) (string, int64, error)
	GetUserByID(ctx context.Context, id string) (*User, error)
}

// MessageUseCase defines message retrieval business logic operations
type MessageUseCase interface {
	GetAllMessages(ctx context.Context, page, limit int) ([]*Message, int64, error)
	DeleteMessage(ctx context.Context, id string) error
}

// ProfileUseCase defines profile business logic operations
type ProfileUseCase interface {
	GetAbout(ctx context.Context) (string, error)
	GetResume(ctx context.Context) (interface{}, error)
	GetResumeDownloadURL(ctx context.Context) (string, error)
}

// GalleryUseCase defines gallery business logic operations
type GalleryUseCase interface {
	CreateGalleryItem(ctx context.Context, title, description string, file *multipart.FileHeader) (*GalleryItem, error)
	GetGalleryItem(ctx context.Context, id string) (*GalleryItem, error)
	GetAllGalleryItems(ctx context.Context) ([]*GalleryItem, error)
	UpdateGalleryItem(ctx context.Context, id string, title, description string, file *multipart.FileHeader) (*GalleryItem, error)
	DeleteGalleryItem(ctx context.Context, id string) error
}

