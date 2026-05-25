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
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query1 := `
		INSERT INTO contact_messages (id, from_email, subject, message, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err = tx.Exec(ctx, query1, message.ID, message.FromEmail, message.Subject, message.Message, message.CreatedAt)
	if err != nil {
		return err
	}

	query2 := `
		INSERT INTO messages (id, from_email, subject, message, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err = tx.Exec(ctx, query2, message.ID, message.FromEmail, message.Subject, message.Message, message.CreatedAt)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
