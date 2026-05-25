package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/quoteyouros/backend/internal/domain"
)

type GalleryRepository struct {
	pool *pgxpool.Pool
}

func NewGalleryRepository(pool *pgxpool.Pool) *GalleryRepository {
	return &GalleryRepository{pool: pool}
}

func (r *GalleryRepository) CreateGalleryItem(ctx context.Context, item *domain.GalleryItem) error {
	query := `
		INSERT INTO gallery_items (id, title, description, image_file_name, image_file_path, image_file_size, image_mime_type, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.pool.Exec(ctx, query,
		item.ID, item.Title, item.Description,
		item.ImageFileName, item.ImageFilePath, item.ImageFileSize, item.ImageMimeType,
		item.CreatedAt, item.UpdatedAt,
	)
	return err
}

func (r *GalleryRepository) GetGalleryItem(ctx context.Context, id string) (*domain.GalleryItem, error) {
	query := `
		SELECT id, title, description, image_file_name, image_file_path, image_file_size, image_mime_type, created_at, updated_at
		FROM gallery_items WHERE id = $1
	`
	item := &domain.GalleryItem{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&item.ID, &item.Title, &item.Description,
		&item.ImageFileName, &item.ImageFilePath, &item.ImageFileSize, &item.ImageMimeType,
		&item.CreatedAt, &item.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (r *GalleryRepository) GetAllGalleryItems(ctx context.Context) ([]*domain.GalleryItem, error) {
	query := `
		SELECT id, title, description, image_file_name, image_file_path, image_file_size, image_mime_type, created_at, updated_at
		FROM gallery_items
		ORDER BY created_at DESC
	`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*domain.GalleryItem, 0)
	for rows.Next() {
		item := &domain.GalleryItem{}
		err := rows.Scan(
			&item.ID, &item.Title, &item.Description,
			&item.ImageFileName, &item.ImageFilePath, &item.ImageFileSize, &item.ImageMimeType,
			&item.CreatedAt, &item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func (r *GalleryRepository) UpdateGalleryItem(ctx context.Context, id string, item *domain.GalleryItem) error {
	query := `
		UPDATE gallery_items
		SET title = $1, description = $2, image_file_name = $3, image_file_path = $4, image_file_size = $5, image_mime_type = $6, updated_at = $7
		WHERE id = $8
	`
	_, err := r.pool.Exec(ctx, query,
		item.Title, item.Description,
		item.ImageFileName, item.ImageFilePath, item.ImageFileSize, item.ImageMimeType,
		item.UpdatedAt, id,
	)
	return err
}

func (r *GalleryRepository) DeleteGalleryItem(ctx context.Context, id string) error {
	query := `DELETE FROM gallery_items WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}
