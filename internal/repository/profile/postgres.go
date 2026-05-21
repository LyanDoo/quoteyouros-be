package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/quoteyouros/backend/internal/domain"
	"github.com/quoteyouros/backend/pkg/logger"
)

type ProfileRepository struct {
	pool *pgxpool.Pool
}

func NewProfileRepository(pool *pgxpool.Pool) *ProfileRepository {
	return &ProfileRepository{pool: pool}
}

// GetProfile retrieves the profile (assuming single profile for admin)
func (r *ProfileRepository) GetProfile(ctx context.Context) (*domain.Profile, error) {
	query := `
		SELECT id, about, resume_file_name, resume_file_size, resume_file_path, 
		       resume_mime_type, resume_uploaded_at, created_at, updated_at
		FROM profiles
		LIMIT 1
	`
	profile := &domain.Profile{}
	err := r.pool.QueryRow(ctx, query).Scan(
		&profile.ID, &profile.About, &profile.ResumeFileName, &profile.ResumeFileSize,
		&profile.ResumeFilePath, &profile.ResumeMimeType, &profile.ResumeUploadedAt,
		&profile.CreatedAt, &profile.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return profile, nil
}

// CreateProfile creates a new profile
func (r *ProfileRepository) CreateProfile(ctx context.Context, profile *domain.Profile) error {
	query := `
		INSERT INTO profiles (id, about, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.pool.Exec(ctx, query, profile.ID, profile.About, profile.CreatedAt, profile.UpdatedAt)
	if err != nil {
		logger.Error("createProfile: failed to create profile", "error", err.Error())
		return err
	}
	logger.Debug("createProfile: profile created successfully", "profile_id", profile.ID)
	return nil
}

// UpdateProfile updates profile information (about text)
func (r *ProfileRepository) UpdateProfile(ctx context.Context, profile *domain.Profile) error {
	query := `
		UPDATE profiles
		SET about = $1, updated_at = $2
		WHERE id = $3
	`
	result, err := r.pool.Exec(ctx, query, profile.About, time.Now(), profile.ID)
	if err != nil {
		logger.Error("updateProfile: failed to update profile", "profile_id", profile.ID, "error", err.Error())
		return err
	}
	if result.RowsAffected() == 0 {
		logger.Warn("updateProfile: profile not found", "profile_id", profile.ID)
		return err
	}
	logger.Debug("updateProfile: profile updated successfully", "profile_id", profile.ID)
	return nil
}

// SaveResume saves resume information
func (r *ProfileRepository) SaveResume(ctx context.Context, profile *domain.Profile) error {
	query := `
		UPDATE profiles
		SET resume_file_name = $1, resume_file_size = $2, resume_file_path = $3,
		    resume_mime_type = $4, resume_uploaded_at = $5, updated_at = $6
		WHERE id = $7
	`
	result, err := r.pool.Exec(ctx, query,
		profile.ResumeFileName, profile.ResumeFileSize, profile.ResumeFilePath,
		profile.ResumeMimeType, profile.ResumeUploadedAt, time.Now(), profile.ID,
	)
	if err != nil {
		logger.Error("saveResume: failed to save resume", "profile_id", profile.ID, "error", err.Error())
		return err
	}
	if result.RowsAffected() == 0 {
		logger.Warn("saveResume: profile not found", "profile_id", profile.ID)
		return err
	}
	logger.Debug("saveResume: resume saved successfully", "profile_id", profile.ID)
	return nil
}

// DeleteResume removes resume information from profile
func (r *ProfileRepository) DeleteResume(ctx context.Context) error {
	query := `
		UPDATE profiles
		SET resume_file_name = NULL, resume_file_size = NULL, resume_file_path = NULL,
		    resume_mime_type = NULL, resume_uploaded_at = NULL, updated_at = $1
		WHERE id = (SELECT id FROM profiles LIMIT 1)
	`
	_, err := r.pool.Exec(ctx, query, time.Now())
	if err != nil {
		logger.Error("deleteResume: failed to delete resume", "error", err.Error())
		return err
	}
	logger.Debug("deleteResume: resume deleted successfully")
	return nil
}
