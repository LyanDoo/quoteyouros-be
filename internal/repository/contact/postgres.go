package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/quoteyouros/backend/internal/domain"
)

type ContactRepository struct {
	pool *pgxpool.Pool
}

func NewContactRepository(pool *pgxpool.Pool) *ContactRepository {
	return &ContactRepository{pool: pool}
}

func (r *ContactRepository) CreateContactMessage(ctx context.Context, message *domain.ContactMessage) error {
	query := `
		INSERT INTO contact_messages (id, from_email, subject, message, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.pool.Exec(ctx, query, message.ID, message.FromEmail, message.Subject, message.Message, message.CreatedAt)
	return err
}
