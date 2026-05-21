package comment

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/quoteyouros/backend/internal/domain"
	"github.com/quoteyouros/backend/pkg/logger"
)

type CommentRepository struct {
	pool *pgxpool.Pool
}

func NewCommentRepository(pool *pgxpool.Pool) *CommentRepository {
	return &CommentRepository{pool: pool}
}

// CreateComment creates a new comment
func (r *CommentRepository) CreateComment(ctx context.Context, comment *domain.Comment) error {
	query := `
		INSERT INTO comments (id, blog_post_id, reply_to_comment_id, author_name, author_email, 
		                       content, rating, is_spam, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.pool.Exec(ctx, query,
		comment.ID, comment.BlogPostID, comment.ReplyToCommentID, comment.AuthorName,
		comment.AuthorEmail, comment.Content, comment.Rating, comment.IsSpam,
		comment.CreatedAt, comment.UpdatedAt,
	)
	if err != nil {
		logger.Error("createComment: failed to create comment", "blog_post_id", comment.BlogPostID, "error", err.Error())
		return err
	}
	logger.Debug("createComment: comment created successfully", "comment_id", comment.ID)
	return nil
}

// GetComment retrieves a single comment by ID
func (r *CommentRepository) GetComment(ctx context.Context, id string) (*domain.Comment, error) {
	query := `
		SELECT id, blog_post_id, reply_to_comment_id, author_name, author_email,
		       content, rating, is_spam, created_at, updated_at
		FROM comments
		WHERE id = $1
	`
	comment := &domain.Comment{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&comment.ID, &comment.BlogPostID, &comment.ReplyToCommentID, &comment.AuthorName,
		&comment.AuthorEmail, &comment.Content, &comment.Rating, &comment.IsSpam,
		&comment.CreatedAt, &comment.UpdatedAt,
	)
	if err != nil {
		logger.Debug("getComment: comment not found", "comment_id", id)
		return nil, err
	}
	logger.Debug("getComment: comment retrieved successfully", "comment_id", id)
	return comment, nil
}

// GetCommentsByBlogPost retrieves all comments for a blog post (excluding replies)
func (r *CommentRepository) GetCommentsByBlogPost(ctx context.Context, blogPostID string, limit, offset int) ([]*domain.Comment, int64, error) {
	query := `
		SELECT id, blog_post_id, reply_to_comment_id, author_name, author_email,
		       content, rating, is_spam, created_at, updated_at
		FROM comments
		WHERE blog_post_id = $1 AND reply_to_comment_id IS NULL AND is_spam = false
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.pool.Query(ctx, query, blogPostID, limit, offset)
	if err != nil {
		logger.Error("getCommentsByBlogPost: failed to query comments", "blog_post_id", blogPostID, "error", err.Error())
		return nil, 0, err
	}
	defer rows.Close()

	comments := []*domain.Comment{}
	for rows.Next() {
		comment := &domain.Comment{}
		err := rows.Scan(
			&comment.ID, &comment.BlogPostID, &comment.ReplyToCommentID, &comment.AuthorName,
			&comment.AuthorEmail, &comment.Content, &comment.Rating, &comment.IsSpam,
			&comment.CreatedAt, &comment.UpdatedAt,
		)
		if err != nil {
			logger.Error("getCommentsByBlogPost: failed to scan row", "error", err.Error())
			return nil, 0, err
		}
		comments = append(comments, comment)
	}

	// Get total count
	countQuery := `
		SELECT COUNT(*) FROM comments
		WHERE blog_post_id = $1 AND reply_to_comment_id IS NULL AND is_spam = false
	`
	var total int64
	err = r.pool.QueryRow(ctx, countQuery, blogPostID).Scan(&total)
	if err != nil {
		logger.Error("getCommentsByBlogPost: failed to get count", "blog_post_id", blogPostID, "error", err.Error())
		return comments, 0, err
	}

	logger.Debug("getCommentsByBlogPost: comments retrieved", "blog_post_id", blogPostID, "count", len(comments))
	return comments, total, nil
}

// GetReplies retrieves all replies to a specific comment
func (r *CommentRepository) GetReplies(ctx context.Context, commentID string, limit, offset int) ([]*domain.Comment, int64, error) {
	query := `
		SELECT id, blog_post_id, reply_to_comment_id, author_name, author_email,
		       content, rating, is_spam, created_at, updated_at
		FROM comments
		WHERE reply_to_comment_id = $1 AND is_spam = false
		ORDER BY created_at ASC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.pool.Query(ctx, query, commentID, limit, offset)
	if err != nil {
		logger.Error("getReplies: failed to query replies", "comment_id", commentID, "error", err.Error())
		return nil, 0, err
	}
	defer rows.Close()

	replies := []*domain.Comment{}
	for rows.Next() {
		reply := &domain.Comment{}
		err := rows.Scan(
			&reply.ID, &reply.BlogPostID, &reply.ReplyToCommentID, &reply.AuthorName,
			&reply.AuthorEmail, &reply.Content, &reply.Rating, &reply.IsSpam,
			&reply.CreatedAt, &reply.UpdatedAt,
		)
		if err != nil {
			logger.Error("getReplies: failed to scan row", "error", err.Error())
			return nil, 0, err
		}
		replies = append(replies, reply)
	}

	// Get total count
	countQuery := `
		SELECT COUNT(*) FROM comments
		WHERE reply_to_comment_id = $1 AND is_spam = false
	`
	var total int64
	err = r.pool.QueryRow(ctx, countQuery, commentID).Scan(&total)
	if err != nil {
		logger.Error("getReplies: failed to get count", "comment_id", commentID, "error", err.Error())
		return replies, 0, err
	}

	logger.Debug("getReplies: replies retrieved", "comment_id", commentID, "count", len(replies))
	return replies, total, nil
}

// UpdateComment updates a comment
func (r *CommentRepository) UpdateComment(ctx context.Context, id string, comment *domain.Comment) error {
	query := `
		UPDATE comments
		SET content = $1, rating = $2, is_spam = $3, updated_at = $4
		WHERE id = $5
	`
	result, err := r.pool.Exec(ctx, query, comment.Content, comment.Rating, comment.IsSpam, time.Now(), id)
	if err != nil {
		logger.Error("updateComment: failed to update comment", "comment_id", id, "error", err.Error())
		return err
	}
	if result.RowsAffected() == 0 {
		logger.Warn("updateComment: comment not found", "comment_id", id)
		return err
	}
	logger.Debug("updateComment: comment updated successfully", "comment_id", id)
	return nil
}

// DeleteComment deletes a comment (cascade will remove replies)
func (r *CommentRepository) DeleteComment(ctx context.Context, id string) error {
	query := `
		DELETE FROM comments
		WHERE id = $1
	`
	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		logger.Error("deleteComment: failed to delete comment", "comment_id", id, "error", err.Error())
		return err
	}
	if result.RowsAffected() == 0 {
		logger.Warn("deleteComment: comment not found", "comment_id", id)
		return err
	}
	logger.Info("deleteComment: comment deleted successfully", "comment_id", id)
	return nil
}
