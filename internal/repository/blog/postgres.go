package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/quoteyouros/backend/internal/domain"
)

type BlogRepository struct {
	pool *pgxpool.Pool
}

func NewBlogRepository(pool *pgxpool.Pool) *BlogRepository {
	return &BlogRepository{pool: pool}
}

func (r *BlogRepository) CreateBlogPost(ctx context.Context, post *domain.BlogPost) error {
	query := `
		INSERT INTO blog_posts (id, title, date, excerpt, content, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.pool.Exec(ctx, query, post.ID, post.Title, post.Date, post.Excerpt, post.Content, post.CreatedAt, post.UpdatedAt)
	return err
}

func (r *BlogRepository) GetBlogPost(ctx context.Context, id string) (*domain.BlogPost, error) {
	query := `
		SELECT id, title, date, excerpt, content, created_at, updated_at
		FROM blog_posts WHERE id = $1
	`
	post := &domain.BlogPost{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&post.ID, &post.Title, &post.Date, &post.Excerpt, &post.Content, &post.CreatedAt, &post.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (r *BlogRepository) GetAllBlogPosts(ctx context.Context, limit, offset int) ([]*domain.BlogPost, int64, error) {
	// Get total count
	countQuery := `SELECT COUNT(*) FROM blog_posts`
	var total int64
	err := r.pool.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, title, date, excerpt, content, created_at, updated_at
		FROM blog_posts
		ORDER BY date DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	posts := make([]*domain.BlogPost, 0)
	for rows.Next() {
		post := &domain.BlogPost{}
		err := rows.Scan(
			&post.ID, &post.Title, &post.Date, &post.Excerpt, &post.Content, &post.CreatedAt, &post.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		posts = append(posts, post)
	}

	return posts, total, nil
}

func (r *BlogRepository) UpdateBlogPost(ctx context.Context, id string, post *domain.BlogPost) error {
	query := `
		UPDATE blog_posts
		SET title = $1, date = $2, excerpt = $3, content = $4, updated_at = $5
		WHERE id = $6
	`
	_, err := r.pool.Exec(ctx, query, post.Title, post.Date, post.Excerpt, post.Content, post.UpdatedAt, id)
	return err
}

func (r *BlogRepository) DeleteBlogPost(ctx context.Context, id string) error {
	query := `DELETE FROM blog_posts WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}
