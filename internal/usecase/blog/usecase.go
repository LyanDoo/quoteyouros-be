package blog

import (
	"context"
	"time"

	"github.com/quoteyouros/backend/internal/domain"
	apperrors "github.com/quoteyouros/backend/pkg/errors"
	"github.com/quoteyouros/backend/pkg/logger"
)

type BlogUseCase struct {
	blogRepo domain.BlogRepository
}

// New creates a new blog use case
func New(blogRepo domain.BlogRepository) *BlogUseCase {
	return &BlogUseCase{
		blogRepo: blogRepo,
	}
}

// CreateBlogPost creates a new blog post
func (u *BlogUseCase) CreateBlogPost(ctx context.Context, req *domain.CreateBlogPostRequest) (*domain.BlogPost, error) {
	logger.Debug("createBlogPost: validating request", "title", req.Title)

	// Parse date (default to today if empty)
	var date time.Time
	if req.Date == "" {
		date = time.Now()
		logger.Debug("createBlogPost: date not provided, using today", "date", date.Format("2006-01-02"))
	} else {
		var err error
		date, err = time.Parse("2006-01-02", req.Date)
		if err != nil {
			logger.Warn("createBlogPost: invalid date format", "date", req.Date, "error", err.Error())
			return nil, apperrors.BadRequest("invalid date format, use YYYY-MM-DD")
		}
	}

	// Create blog post
	post := domain.NewBlogPost(req.Title, req.Excerpt, req.Content, date)

	logger.Info("createBlogPost: creating blog post", "post_id", post.ID, "title", post.Title)
	if err := u.blogRepo.CreateBlogPost(ctx, post); err != nil {
		logger.Error("createBlogPost: failed to create blog post", "post_id", post.ID, "error", err.Error())
		return nil, apperrors.InternalServerError("failed to create blog post: " + err.Error())
	}

	logger.Info("createBlogPost: blog post created successfully", "post_id", post.ID, "title", post.Title)
	return post, nil
}

// GetBlogPost retrieves a single blog post
func (u *BlogUseCase) GetBlogPost(ctx context.Context, id string) (*domain.BlogPost, error) {
	logger.Debug("getBlogPost: retrieving blog post", "post_id", id)

	post, err := u.blogRepo.GetBlogPost(ctx, id)
	if err != nil {
		// Check if it's a "not found" error
		if err.Error() == "no rows in result set" {
			logger.Warn("getBlogPost: blog post not found", "post_id", id)
			return nil, apperrors.NotFound("blog post not found")
		}
		logger.Error("getBlogPost: failed to retrieve blog post", "post_id", id, "error", err.Error())
		return nil, apperrors.InternalServerError("failed to retrieve blog post: " + err.Error())
	}

	logger.Debug("getBlogPost: blog post retrieved successfully", "post_id", id, "title", post.Title)
	return post, nil
}

// GetAllBlogPosts retrieves all blog posts with pagination
func (u *BlogUseCase) GetAllBlogPosts(ctx context.Context, page, limit int) ([]*domain.BlogPost, int64, error) {
	logger.Debug("getAllBlogPosts: retrieving blog posts", "page", page, "limit", limit)

	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	posts, total, err := u.blogRepo.GetAllBlogPosts(ctx, limit, offset)
	if err != nil {
		logger.Error("getAllBlogPosts: failed to retrieve blog posts", "page", page, "limit", limit, "error", err.Error())
		return nil, 0, apperrors.InternalServerError("failed to retrieve blog posts: " + err.Error())
	}

	logger.Info("getAllBlogPosts: blog posts retrieved successfully", "count", len(posts), "total", total, "page", page)
	return posts, total, nil
}

// UpdateBlogPost updates an existing blog post
func (u *BlogUseCase) UpdateBlogPost(ctx context.Context, id string, req *domain.UpdateBlogPostRequest) (*domain.BlogPost, error) {
	logger.Debug("updateBlogPost: validating request", "post_id", id, "title", req.Title)

	// Check if blog post exists
	existingPost, err := u.blogRepo.GetBlogPost(ctx, id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			logger.Warn("updateBlogPost: blog post not found", "post_id", id)
			return nil, apperrors.NotFound("blog post not found")
		}
		logger.Error("updateBlogPost: failed to check existing blog post", "post_id", id, "error", err.Error())
		return nil, apperrors.InternalServerError("failed to update blog post: " + err.Error())
	}

	// Parse date (keep existing if empty)
	if req.Date != "" {
		date, err := time.Parse("2006-01-02", req.Date)
		if err != nil {
			logger.Warn("updateBlogPost: invalid date format", "date", req.Date, "error", err.Error())
			return nil, apperrors.BadRequest("invalid date format, use YYYY-MM-DD")
		}
		existingPost.Date = date
	}

	// Update fields
	existingPost.Title = req.Title
	existingPost.Excerpt = req.Excerpt
	existingPost.Content = req.Content
	existingPost.UpdatedAt = time.Now()

	logger.Info("updateBlogPost: updating blog post", "post_id", id, "title", req.Title)
	if err := u.blogRepo.UpdateBlogPost(ctx, id, existingPost); err != nil {
		logger.Error("updateBlogPost: failed to update blog post", "post_id", id, "error", err.Error())
		return nil, apperrors.InternalServerError("failed to update blog post: " + err.Error())
	}

	logger.Info("updateBlogPost: blog post updated successfully", "post_id", id, "title", req.Title)
	return existingPost, nil
}

// DeleteBlogPost deletes a blog post
func (u *BlogUseCase) DeleteBlogPost(ctx context.Context, id string) error {
	logger.Debug("deleteBlogPost: validating blog post existence", "post_id", id)

	// Check if blog post exists
	_, err := u.blogRepo.GetBlogPost(ctx, id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			logger.Warn("deleteBlogPost: blog post not found", "post_id", id)
			return apperrors.NotFound("blog post not found")
		}
		logger.Error("deleteBlogPost: failed to check blog post", "post_id", id, "error", err.Error())
		return apperrors.InternalServerError("failed to delete blog post: " + err.Error())
	}

	logger.Info("deleteBlogPost: deleting blog post", "post_id", id)
	if err := u.blogRepo.DeleteBlogPost(ctx, id); err != nil {
		logger.Error("deleteBlogPost: failed to delete blog post", "post_id", id, "error", err.Error())
		return apperrors.InternalServerError("failed to delete blog post: " + err.Error())
	}

	logger.Info("deleteBlogPost: blog post deleted successfully", "post_id", id)
	return nil
}
