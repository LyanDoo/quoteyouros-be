package comment

import (
	"context"

	"github.com/quoteyouros/backend/internal/domain"
	apperrors "github.com/quoteyouros/backend/pkg/errors"
	"github.com/quoteyouros/backend/pkg/logger"
)

type CommentUseCase struct {
	commentRepo domain.CommentRepository
	blogRepo    domain.BlogRepository
}

// New creates a new comment use case
func New(commentRepo domain.CommentRepository, blogRepo domain.BlogRepository) *CommentUseCase {
	return &CommentUseCase{
		commentRepo: commentRepo,
		blogRepo:    blogRepo,
	}
}

// CreateComment creates a new comment
func (u *CommentUseCase) CreateComment(ctx context.Context, req *domain.CreateCommentRequest, blogPostID string) (*domain.Comment, error) {
	logger.Debug("createComment: creating comment", "blog_post_id", blogPostID, "author", req.AuthorName)

	// Verify blog post exists
	_, err := u.blogRepo.GetBlogPost(ctx, blogPostID)
	if err != nil {
		logger.Warn("createComment: blog post not found", "blog_post_id", blogPostID)
		return nil, apperrors.NotFound("blog post not found")
	}

	// Verify reply_to_comment exists if provided
	if req.ReplyToCommentID != nil && *req.ReplyToCommentID != "" {
		_, err := u.commentRepo.GetComment(ctx, *req.ReplyToCommentID)
		if err != nil {
			logger.Warn("createComment: reply_to_comment not found", "reply_to_comment_id", *req.ReplyToCommentID)
			return nil, apperrors.NotFound("comment to reply to not found")
		}
	}

	// Create comment
	comment := domain.NewComment(blogPostID, req.AuthorName, req.AuthorEmail, req.Content, req.ReplyToCommentID, req.Rating)

	logger.Info("createComment: saving comment to database", "blog_post_id", blogPostID, "comment_id", comment.ID)
	if err := u.commentRepo.CreateComment(ctx, comment); err != nil {
		logger.Error("createComment: failed to create comment", "error", err.Error())
		return nil, apperrors.InternalServerError("failed to create comment: " + err.Error())
	}

	logger.Info("createComment: comment created successfully", "comment_id", comment.ID)
	return comment, nil
}

// GetComment retrieves a single comment
func (u *CommentUseCase) GetComment(ctx context.Context, commentID string) (*domain.Comment, error) {
	logger.Debug("getComment: retrieving comment", "comment_id", commentID)

	comment, err := u.commentRepo.GetComment(ctx, commentID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			logger.Warn("getComment: comment not found", "comment_id", commentID)
			return nil, apperrors.NotFound("comment not found")
		}
		logger.Error("getComment: failed to retrieve comment", "error", err.Error())
		return nil, apperrors.InternalServerError("failed to retrieve comment: " + err.Error())
	}

	logger.Debug("getComment: comment retrieved successfully", "comment_id", commentID)
	return comment, nil
}

// GetCommentsByBlogPost retrieves all comments for a blog post with pagination
func (u *CommentUseCase) GetCommentsByBlogPost(ctx context.Context, blogPostID string, page, limit int) ([]*domain.Comment, int64, error) {
	logger.Debug("getCommentsByBlogPost: retrieving comments", "blog_post_id", blogPostID, "page", page, "limit", limit)

	// Verify blog post exists
	_, err := u.blogRepo.GetBlogPost(ctx, blogPostID)
	if err != nil {
		logger.Warn("getCommentsByBlogPost: blog post not found", "blog_post_id", blogPostID)
		return nil, 0, apperrors.NotFound("blog post not found")
	}

	// Validate pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	comments, total, err := u.commentRepo.GetCommentsByBlogPost(ctx, blogPostID, limit, offset)
	if err != nil {
		logger.Error("getCommentsByBlogPost: failed to retrieve comments", "error", err.Error())
		return nil, 0, apperrors.InternalServerError("failed to retrieve comments: " + err.Error())
	}

	logger.Debug("getCommentsByBlogPost: comments retrieved", "blog_post_id", blogPostID, "count", len(comments), "total", total)
	return comments, total, nil
}

// GetReplies retrieves all replies to a comment with pagination
func (u *CommentUseCase) GetReplies(ctx context.Context, commentID string, page, limit int) ([]*domain.Comment, int64, error) {
	logger.Debug("getReplies: retrieving replies", "comment_id", commentID, "page", page, "limit", limit)

	// Verify parent comment exists
	_, err := u.commentRepo.GetComment(ctx, commentID)
	if err != nil {
		logger.Warn("getReplies: comment not found", "comment_id", commentID)
		return nil, 0, apperrors.NotFound("comment not found")
	}

	// Validate pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	replies, total, err := u.commentRepo.GetReplies(ctx, commentID, limit, offset)
	if err != nil {
		logger.Error("getReplies: failed to retrieve replies", "error", err.Error())
		return nil, 0, apperrors.InternalServerError("failed to retrieve replies: " + err.Error())
	}

	logger.Debug("getReplies: replies retrieved", "comment_id", commentID, "count", len(replies), "total", total)
	return replies, total, nil
}

// UpdateComment updates a comment (admin only)
func (u *CommentUseCase) UpdateComment(ctx context.Context, commentID string, req *domain.UpdateCommentRequest) (*domain.Comment, error) {
	logger.Debug("updateComment: updating comment", "comment_id", commentID)

	// Get existing comment
	comment, err := u.commentRepo.GetComment(ctx, commentID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			logger.Warn("updateComment: comment not found", "comment_id", commentID)
			return nil, apperrors.NotFound("comment not found")
		}
		logger.Error("updateComment: failed to retrieve comment", "error", err.Error())
		return nil, apperrors.InternalServerError("failed to retrieve comment: " + err.Error())
	}

	// Update fields
	comment.Content = req.Content
	if req.Rating != nil {
		comment.Rating = req.Rating
	}
	if req.IsSpam != nil {
		comment.IsSpam = *req.IsSpam
	}

	logger.Info("updateComment: updating comment in database", "comment_id", commentID)
	if err := u.commentRepo.UpdateComment(ctx, commentID, comment); err != nil {
		logger.Error("updateComment: failed to update comment", "error", err.Error())
		return nil, apperrors.InternalServerError("failed to update comment: " + err.Error())
	}

	// Fetch updated comment
	updatedComment, err := u.commentRepo.GetComment(ctx, commentID)
	if err != nil {
		logger.Error("updateComment: failed to retrieve updated comment", "error", err.Error())
		return nil, apperrors.InternalServerError("failed to retrieve updated comment: " + err.Error())
	}

	logger.Info("updateComment: comment updated successfully", "comment_id", commentID)
	return updatedComment, nil
}

// DeleteComment deletes a comment (admin only)
func (u *CommentUseCase) DeleteComment(ctx context.Context, commentID string) error {
	logger.Debug("deleteComment: deleting comment", "comment_id", commentID)

	// Verify comment exists
	_, err := u.commentRepo.GetComment(ctx, commentID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			logger.Warn("deleteComment: comment not found", "comment_id", commentID)
			return apperrors.NotFound("comment not found")
		}
		logger.Error("deleteComment: failed to retrieve comment", "error", err.Error())
		return apperrors.InternalServerError("failed to retrieve comment: " + err.Error())
	}

	logger.Info("deleteComment: deleting comment from database", "comment_id", commentID)
	if err := u.commentRepo.DeleteComment(ctx, commentID); err != nil {
		logger.Error("deleteComment: failed to delete comment", "error", err.Error())
		return apperrors.InternalServerError("failed to delete comment: " + err.Error())
	}

	logger.Info("deleteComment: comment deleted successfully", "comment_id", commentID)
	return nil
}
