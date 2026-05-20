package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/quoteyouros/backend/internal/domain"
)

type MessageRepository struct {
	pool *pgxpool.Pool
}

func NewMessageRepository(pool *pgxpool.Pool) *MessageRepository {
	return &MessageRepository{pool: pool}
}

func (r *MessageRepository) GetAllMessages(ctx context.Context, limit, offset int) ([]*domain.Message, int64, error) {
	// Get total count
	countQuery := `SELECT COUNT(*) FROM messages`
	var total int64
	err := r.pool.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, from_email, subject, message, created_at
		FROM messages
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	messages := make([]*domain.Message, 0)
	for rows.Next() {
		msg := &domain.Message{}
		err := rows.Scan(
			&msg.ID, &msg.FromEmail, &msg.Subject, &msg.Message, &msg.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		messages = append(messages, msg)
	}

	return messages, total, nil
}

func (r *MessageRepository) DeleteMessage(ctx context.Context, id string) error {
	query := `DELETE FROM messages WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}
