package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/quoteyouros/backend/internal/domain"
	"github.com/quoteyouros/backend/pkg/logger"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	logger.Debug("CreateUser: inserting user into database", "email", user.Email, "user_id", user.ID)

	query := `
		INSERT INTO users (id, email, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.pool.Exec(ctx, query, user.ID, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)

	if err != nil {
		logger.Error("CreateUser: failed to insert user", "email", user.Email, "user_id", user.ID, "error", err.Error())
		return err
	}

	logger.Info("CreateUser: user inserted successfully", "email", user.Email, "user_id", user.ID)
	return nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	logger.Debug("GetUserByEmail: retrieving user from database", "email", email)

	query := `
		SELECT id, email, password, created_at, updated_at
		FROM users WHERE email = $1
	`
	user := &domain.User{}
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Debug("GetUserByEmail: user not found", "email", email)
			return nil, nil
		}
		logger.Error("GetUserByEmail: database error", "email", email, "error", err.Error())
		return nil, err
	}
	logger.Debug("GetUserByEmail: user found", "email", email, "user_id", user.ID)
	return user, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	logger.Debug("GetUserByID: retrieving user from database", "user_id", id)

	query := `
		SELECT id, email, password, created_at, updated_at
		FROM users WHERE id = $1
	`
	user := &domain.User{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Debug("GetUserByID: user not found", "user_id", id)
			return nil, nil
		}
		logger.Error("GetUserByID: database error", "user_id", id, "error", err.Error())
		return nil, err
	}
	logger.Debug("GetUserByID: user found", "user_id", id, "email", user.Email)
	return user, nil
}
