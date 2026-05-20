package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/quoteyouros/backend/internal/domain"
	blogusecase "github.com/quoteyouros/backend/internal/usecase/blog"
	apperrors "github.com/quoteyouros/backend/pkg/errors"
	"github.com/quoteyouros/backend/pkg/logger"
	"github.com/quoteyouros/backend/pkg/response"
)

type BlogHandler struct {
	usecase *blogusecase.BlogUseCase
}

// NewBlogHandler creates a new blog handler
func NewBlogHandler(usecase *blogusecase.BlogUseCase) *BlogHandler {
	return &BlogHandler{usecase: usecase}
}

// GetAllBlogPosts retrieves all blog posts with pagination
// GET /api/blog
func (h *BlogHandler) GetAllBlogPosts(c *fiber.Ctx) error {
	// Get pagination parameters from query
	pageStr := c.Query("page", "1")
	limitStr := c.Query("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		logger.Warn("getAllBlogPosts: invalid page parameter", "page", pageStr)
		return response.ErrorResponseJSON(c, fiber.StatusBadRequest, "invalid page parameter", fiber.StatusBadRequest)
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		logger.Warn("getAllBlogPosts: invalid limit parameter", "limit", limitStr)
		return response.ErrorResponseJSON(c, fiber.StatusBadRequest, "invalid limit parameter", fiber.StatusBadRequest)
	}

	logger.Debug("getAllBlogPosts: retrieving blog posts", "page", page, "limit", limit)
	posts, total, err := h.usecase.GetAllBlogPosts(c.Context(), page, limit)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			logger.Warn("getAllBlogPosts: failed to retrieve posts", "page", page, "error", appErr.Details)
			return response.ErrorResponseJSON(c, appErr.Code, appErr.Details, appErr.Code)
		}
		logger.Error("getAllBlogPosts: unexpected error", "page", page, "error", err.Error())
		return response.ErrorResponseJSON(c, fiber.StatusInternalServerError, "failed to retrieve blog posts", fiber.StatusInternalServerError)
	}

	logger.Info("getAllBlogPosts: posts retrieved successfully", "count", len(posts), "total", total)
	return response.PaginatedSuccessResponse(c, fiber.StatusOK, posts, page, limit, total, "Blog posts retrieved successfully")
}

// GetBlogPost retrieves a single blog post by ID
// GET /api/blog/:id
func (h *BlogHandler) GetBlogPost(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		logger.Warn("getBlogPost: missing blog post ID")
		return response.ErrorResponseJSON(c, fiber.StatusBadRequest, "blog post ID is required", fiber.StatusBadRequest)
	}

	logger.Debug("getBlogPost: retrieving blog post", "post_id", id)
	post, err := h.usecase.GetBlogPost(c.Context(), id)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			logger.Warn("getBlogPost: failed to retrieve post", "post_id", id, "error", appErr.Details)
			return response.ErrorResponseJSON(c, appErr.Code, appErr.Details, appErr.Code)
		}
		logger.Error("getBlogPost: unexpected error", "post_id", id, "error", err.Error())
		return response.ErrorResponseJSON(c, fiber.StatusInternalServerError, "failed to retrieve blog post", fiber.StatusInternalServerError)
	}

	logger.Debug("getBlogPost: blog post retrieved successfully", "post_id", id)
	return response.SuccessResponse(c, fiber.StatusOK, post, "Blog post retrieved successfully")
}

// CreateBlogPost creates a new blog post
// POST /api/blog
func (h *BlogHandler) CreateBlogPost(c *fiber.Ctx) error {
	var req domain.CreateBlogPostRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Error("createBlogPost: failed to parse request body", "error", err.Error())
		return response.ErrorResponseJSON(c, fiber.StatusBadRequest, "invalid request body", fiber.StatusBadRequest)
	}

	logger.Info("createBlogPost: attempting to create blog post", "title", req.Title)
	post, err := h.usecase.CreateBlogPost(c.Context(), &req)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			logger.Warn("createBlogPost: failed to create post", "title", req.Title, "error", appErr.Details)
			return response.ErrorResponseJSON(c, appErr.Code, appErr.Details, appErr.Code)
		}
		logger.Error("createBlogPost: unexpected error", "title", req.Title, "error", err.Error())
		return response.ErrorResponseJSON(c, fiber.StatusInternalServerError, "failed to create blog post", fiber.StatusInternalServerError)
	}

	logger.Info("createBlogPost: blog post created successfully", "post_id", post.ID, "title", post.Title)
	return response.SuccessResponse(c, fiber.StatusCreated, post, "Blog post created successfully")
}

// UpdateBlogPost updates an existing blog post
// PUT /api/blog/:id
func (h *BlogHandler) UpdateBlogPost(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		logger.Warn("updateBlogPost: missing blog post ID")
		return response.ErrorResponseJSON(c, fiber.StatusBadRequest, "blog post ID is required", fiber.StatusBadRequest)
	}

	var req domain.UpdateBlogPostRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Error("updateBlogPost: failed to parse request body", "post_id", id, "error", err.Error())
		return response.ErrorResponseJSON(c, fiber.StatusBadRequest, "invalid request body", fiber.StatusBadRequest)
	}

	logger.Info("updateBlogPost: attempting to update blog post", "post_id", id, "title", req.Title)
	post, err := h.usecase.UpdateBlogPost(c.Context(), id, &req)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			logger.Warn("updateBlogPost: failed to update post", "post_id", id, "error", appErr.Details)
			return response.ErrorResponseJSON(c, appErr.Code, appErr.Details, appErr.Code)
		}
		logger.Error("updateBlogPost: unexpected error", "post_id", id, "error", err.Error())
		return response.ErrorResponseJSON(c, fiber.StatusInternalServerError, "failed to update blog post", fiber.StatusInternalServerError)
	}

	logger.Info("updateBlogPost: blog post updated successfully", "post_id", id, "title", post.Title)
	return response.SuccessResponse(c, fiber.StatusOK, post, "Blog post updated successfully")
}

// DeleteBlogPost deletes a blog post
// DELETE /api/blog/:id
func (h *BlogHandler) DeleteBlogPost(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		logger.Warn("deleteBlogPost: missing blog post ID")
		return response.ErrorResponseJSON(c, fiber.StatusBadRequest, "blog post ID is required", fiber.StatusBadRequest)
	}

	logger.Info("deleteBlogPost: attempting to delete blog post", "post_id", id)
	err := h.usecase.DeleteBlogPost(c.Context(), id)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			logger.Warn("deleteBlogPost: failed to delete post", "post_id", id, "error", appErr.Details)
			return response.ErrorResponseJSON(c, appErr.Code, appErr.Details, appErr.Code)
		}
		logger.Error("deleteBlogPost: unexpected error", "post_id", id, "error", err.Error())
		return response.ErrorResponseJSON(c, fiber.StatusInternalServerError, "failed to delete blog post", fiber.StatusInternalServerError)
	}

	logger.Info("deleteBlogPost: blog post deleted successfully", "post_id", id)
	return response.SuccessResponse(c, fiber.StatusOK, fiber.Map{}, "Blog post deleted successfully")
}
