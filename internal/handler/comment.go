package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/quoteyouros/backend/internal/domain"
	commentusecase "github.com/quoteyouros/backend/internal/usecase/comment"
	apperrors "github.com/quoteyouros/backend/pkg/errors"
	"github.com/quoteyouros/backend/pkg/logger"
	"github.com/quoteyouros/backend/pkg/response"
)

type CommentHandler struct {
	usecase *commentusecase.CommentUseCase
}

// NewCommentHandler creates a new comment handler
func NewCommentHandler(usecase *commentusecase.CommentUseCase) *CommentHandler {
	return &CommentHandler{usecase: usecase}
}

// CreateComment creates a new comment on a blog post
// POST /api/blog/:id/comments
func (h *CommentHandler) CreateComment(c *fiber.Ctx) error {
	blogPostID := c.Params("id")
	if blogPostID == "" {
		logger.Error("createComment: blog_post_id is required")
		return response.ErrorResponseJSON(c, fiber.StatusBadRequest, "blog_post_id is required", fiber.StatusBadRequest)
	}

	var req domain.CreateCommentRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Error("createComment: failed to parse request body", "error", err.Error())
		return response.ErrorResponseJSON(c, fiber.StatusBadRequest, "invalid request body", fiber.StatusBadRequest)
	}

	logger.Info("createComment: attempting to create comment", "blog_post_id", blogPostID, "author", req.AuthorName)
	comment, err := h.usecase.CreateComment(c.Context(), &req, blogPostID)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			logger.Warn("createComment: failed to create comment", "error", appErr.Details)
			return response.ErrorResponseJSON(c, appErr.Code, appErr.Details, appErr.Code)
		}
		logger.Error("createComment: unexpected error", "error", err.Error())
		return response.ErrorResponseJSON(c, fiber.StatusInternalServerError, "failed to create comment", fiber.StatusInternalServerError)
	}

	// Build response
	commentResp := buildCommentResponse(comment)

	logger.Info("createComment: comment created successfully", "comment_id", comment.ID)
	return response.SuccessResponse(c, fiber.StatusCreated, commentResp, "Comment created successfully")
}

// GetCommentsByBlogPost retrieves all comments for a blog post (paginated)
// GET /api/blog/:id/comments?page=1&limit=10
func (h *CommentHandler) GetCommentsByBlogPost(c *fiber.Ctx) error {
	blogPostID := c.Params("id")
	if blogPostID == "" {
		logger.Error("getCommentsByBlogPost: blog_post_id is required")
		return response.ErrorResponseJSON(c, fiber.StatusBadRequest, "blog_post_id is required", fiber.StatusBadRequest)
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	logger.Debug("getCommentsByBlogPost: retrieving comments", "blog_post_id", blogPostID, "page", page, "limit", limit)
	comments, total, err := h.usecase.GetCommentsByBlogPost(c.Context(), blogPostID, page, limit)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			logger.Warn("getCommentsByBlogPost: failed to retrieve comments", "error", appErr.Details)
			return response.ErrorResponseJSON(c, appErr.Code, appErr.Details, appErr.Code)
		}
		logger.Error("getCommentsByBlogPost: unexpected error", "error", err.Error())
		return response.ErrorResponseJSON(c, fiber.StatusInternalServerError, "failed to retrieve comments", fiber.StatusInternalServerError)
	}

	// Build paginated response
	totalPages := (int(total) + limit - 1) / limit
	commentResps := []domain.CommentResponse{}
	for _, comment := range comments {
		commentResps = append(commentResps, buildCommentResponse(comment))
	}

	listResp := domain.CommentListResponse{
		Comments:   commentResps,
		Total:      int(total),
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}

	logger.Info("getCommentsByBlogPost: comments retrieved", "blog_post_id", blogPostID, "count", len(comments))
	return response.SuccessResponse(c, fiber.StatusOK, listResp, "Comments retrieved successfully")
}

// GetReplies retrieves all replies to a comment
// GET /api/comments/:id/replies?page=1&limit=10
func (h *CommentHandler) GetReplies(c *fiber.Ctx) error {
	commentID := c.Params("id")
	if commentID == "" {
		logger.Error("getReplies: comment_id is required")
		return response.ErrorResponseJSON(c, fiber.StatusBadRequest, "comment_id is required", fiber.StatusBadRequest)
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	logger.Debug("getReplies: retrieving replies", "comment_id", commentID, "page", page, "limit", limit)
	replies, total, err := h.usecase.GetReplies(c.Context(), commentID, page, limit)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			logger.Warn("getReplies: failed to retrieve replies", "error", appErr.Details)
			return response.ErrorResponseJSON(c, appErr.Code, appErr.Details, appErr.Code)
		}
		logger.Error("getReplies: unexpected error", "error", err.Error())
		return response.ErrorResponseJSON(c, fiber.StatusInternalServerError, "failed to retrieve replies", fiber.StatusInternalServerError)
	}

	// Build paginated response
	totalPages := (int(total) + limit - 1) / limit
	replyResps := []domain.CommentResponse{}
	for _, reply := range replies {
		replyResps = append(replyResps, buildCommentResponse(reply))
	}

	listResp := domain.CommentListResponse{
		Comments:   replyResps,
		Total:      int(total),
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}

	logger.Info("getReplies: replies retrieved", "comment_id", commentID, "count", len(replies))
	return response.SuccessResponse(c, fiber.StatusOK, listResp, "Replies retrieved successfully")
}

// UpdateComment updates a comment (admin only)
// PUT /api/comments/:id
func (h *CommentHandler) UpdateComment(c *fiber.Ctx) error {
	commentID := c.Params("id")
	if commentID == "" {
		logger.Error("updateComment: comment_id is required")
		return response.ErrorResponseJSON(c, fiber.StatusBadRequest, "comment_id is required", fiber.StatusBadRequest)
	}

	var req domain.UpdateCommentRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Error("updateComment: failed to parse request body", "error", err.Error())
		return response.ErrorResponseJSON(c, fiber.StatusBadRequest, "invalid request body", fiber.StatusBadRequest)
	}

	logger.Info("updateComment: attempting to update comment", "comment_id", commentID)
	comment, err := h.usecase.UpdateComment(c.Context(), commentID, &req)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			logger.Warn("updateComment: failed to update comment", "error", appErr.Details)
			return response.ErrorResponseJSON(c, appErr.Code, appErr.Details, appErr.Code)
		}
		logger.Error("updateComment: unexpected error", "error", err.Error())
		return response.ErrorResponseJSON(c, fiber.StatusInternalServerError, "failed to update comment", fiber.StatusInternalServerError)
	}

	// Build response
	commentResp := buildCommentResponse(comment)

	logger.Info("updateComment: comment updated successfully", "comment_id", commentID)
	return response.SuccessResponse(c, fiber.StatusOK, commentResp, "Comment updated successfully")
}

// DeleteComment deletes a comment (admin only)
// DELETE /api/comments/:id
func (h *CommentHandler) DeleteComment(c *fiber.Ctx) error {
	commentID := c.Params("id")
	if commentID == "" {
		logger.Error("deleteComment: comment_id is required")
		return response.ErrorResponseJSON(c, fiber.StatusBadRequest, "comment_id is required", fiber.StatusBadRequest)
	}

	logger.Info("deleteComment: attempting to delete comment", "comment_id", commentID)
	err := h.usecase.DeleteComment(c.Context(), commentID)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			logger.Warn("deleteComment: failed to delete comment", "error", appErr.Details)
			return response.ErrorResponseJSON(c, appErr.Code, appErr.Details, appErr.Code)
		}
		logger.Error("deleteComment: unexpected error", "error", err.Error())
		return response.ErrorResponseJSON(c, fiber.StatusInternalServerError, "failed to delete comment", fiber.StatusInternalServerError)
	}

	logger.Info("deleteComment: comment deleted successfully", "comment_id", commentID)
	return response.SuccessResponse(c, fiber.StatusOK, fiber.Map{"message": "Comment deleted successfully"}, "Comment deleted successfully")
}

// Helper function to build CommentResponse
func buildCommentResponse(comment *domain.Comment) domain.CommentResponse {
	return domain.CommentResponse{
		ID:               comment.ID,
		BlogPostID:       comment.BlogPostID,
		ReplyToCommentID: comment.ReplyToCommentID,
		AuthorName:       comment.AuthorName,
		AuthorEmail:      comment.AuthorEmail,
		Content:          comment.Content,
		Rating:           comment.Rating,
		IsSpam:           comment.IsSpam,
		CreatedAt:        comment.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:        comment.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
