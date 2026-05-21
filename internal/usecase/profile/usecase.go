package profile

import (
	"context"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
	"github.com/quoteyouros/backend/internal/domain"
	apperrors "github.com/quoteyouros/backend/pkg/errors"
	"github.com/quoteyouros/backend/pkg/fileupload"
	"github.com/quoteyouros/backend/pkg/logger"
)

type ProfileUseCase struct {
	profileRepo domain.ProfileRepository
	fileUpload  *fileupload.FileUploadService
}

// New creates a new profile use case
func New(profileRepo domain.ProfileRepository, fileUpload *fileupload.FileUploadService) *ProfileUseCase {
	return &ProfileUseCase{
		profileRepo: profileRepo,
		fileUpload:  fileUpload,
	}
}

// GetProfile retrieves the profile
func (u *ProfileUseCase) GetProfile(ctx context.Context) (*domain.Profile, error) {
	logger.Debug("getProfile: retrieving profile")

	profile, err := u.profileRepo.GetProfile(ctx)
	if err != nil {
		if err.Error() == "no rows in result set" {
			logger.Warn("getProfile: profile not found")
			return nil, apperrors.NotFound("profile not found")
		}
		logger.Error("getProfile: failed to retrieve profile", "error", err.Error())
		return nil, apperrors.InternalServerError("failed to retrieve profile: " + err.Error())
	}

	logger.Debug("getProfile: profile retrieved successfully")
	return profile, nil
}

// UpdateProfile updates profile information
func (u *ProfileUseCase) UpdateProfile(ctx context.Context, req *domain.UpdateProfileRequest) (*domain.Profile, error) {
	logger.Debug("updateProfile: updating profile", "about_length", len(req.About))

	// Get existing profile or create new one
	profile, err := u.profileRepo.GetProfile(ctx)
	if err != nil {
		if err.Error() == "no rows in result set" {
			// Create new profile if it doesn't exist
			profile = &domain.Profile{
				ID:        uuid.New().String(),
				About:     req.About,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			logger.Info("updateProfile: creating new profile", "profile_id", profile.ID)
			if err := u.profileRepo.CreateProfile(ctx, profile); err != nil {
				logger.Error("updateProfile: failed to create profile", "error", err.Error())
				return nil, apperrors.InternalServerError("failed to create profile: " + err.Error())
			}
			return profile, nil
		}
		logger.Error("updateProfile: failed to retrieve profile", "error", err.Error())
		return nil, apperrors.InternalServerError("failed to retrieve profile: " + err.Error())
	}

	// Update profile
	profile.About = req.About
	logger.Info("updateProfile: updating profile", "profile_id", profile.ID)
	if err := u.profileRepo.UpdateProfile(ctx, profile); err != nil {
		logger.Error("updateProfile: failed to update profile", "error", err.Error())
		return nil, apperrors.InternalServerError("failed to update profile: " + err.Error())
	}

	logger.Info("updateProfile: profile updated successfully", "profile_id", profile.ID)
	return profile, nil
}

// UploadResume uploads a new resume file
func (u *ProfileUseCase) UploadResume(ctx context.Context, file *multipart.FileHeader) (*domain.Profile, error) {
	logger.Info("uploadResume: starting resume upload", "filename", file.Filename, "size", file.Size)

	// Validate file
	if err := u.fileUpload.ValidateResume(file); err != nil {
		logger.Warn("uploadResume: validation failed", "filename", file.Filename, "error", err.Error())
		return nil, apperrors.BadRequest(err.Error())
	}

	// Get existing profile or create new one
	profile, err := u.profileRepo.GetProfile(ctx)
	if err != nil {
		if err.Error() == "no rows in result set" {
			// Create new profile
			profile = &domain.Profile{
				ID:        uuid.New().String(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			logger.Info("uploadResume: creating new profile", "profile_id", profile.ID)
			if err := u.profileRepo.CreateProfile(ctx, profile); err != nil {
				logger.Error("uploadResume: failed to create profile", "error", err.Error())
				return nil, apperrors.InternalServerError("failed to create profile: " + err.Error())
			}
		} else {
			logger.Error("uploadResume: failed to retrieve profile", "error", err.Error())
			return nil, apperrors.InternalServerError("failed to retrieve profile: " + err.Error())
		}
	}

	// Delete old resume if it exists
	if profile.ResumeFilePath != nil && *profile.ResumeFilePath != "" {
		logger.Debug("uploadResume: deleting old resume", "path", *profile.ResumeFilePath)
		if err := u.fileUpload.DeleteResume(*profile.ResumeFilePath); err != nil {
			logger.Warn("uploadResume: failed to delete old resume", "error", err.Error())
			// Don't fail the upload if old file deletion fails
		}
	}

	// Upload new file
	fileName, filePath, fileSize, mimeType, err := u.fileUpload.UploadResume(file)
	if err != nil {
		logger.Error("uploadResume: failed to upload file", "filename", file.Filename, "error", err.Error())
		return nil, apperrors.BadRequest(err.Error())
	}

	// Update profile with resume information
	now := time.Now()
	profile.ResumeFileName = &fileName
	profile.ResumeFilePath = &filePath
	profile.ResumeFileSize = &fileSize
	profile.ResumeMimeType = &mimeType
	profile.ResumeUploadedAt = &now

	logger.Info("uploadResume: saving resume metadata to database", "profile_id", profile.ID, "filename", fileName)
	if err := u.profileRepo.SaveResume(ctx, profile); err != nil {
		logger.Error("uploadResume: failed to save resume metadata", "profile_id", profile.ID, "error", err.Error())
		// Clean up uploaded file
		u.fileUpload.DeleteResume(filePath)
		return nil, apperrors.InternalServerError("failed to save resume: " + err.Error())
	}

	logger.Info("uploadResume: resume uploaded successfully", "profile_id", profile.ID, "filename", fileName)
	return profile, nil
}

// GetResumeFilePath retrieves the resume file path for download
func (u *ProfileUseCase) GetResumeFilePath(ctx context.Context) (string, error) {
	logger.Debug("getResumeFilePath: retrieving resume file path")

	profile, err := u.profileRepo.GetProfile(ctx)
	if err != nil {
		logger.Error("getResumeFilePath: failed to retrieve profile", "error", err.Error())
		return "", apperrors.InternalServerError("failed to retrieve profile")
	}

	if profile.ResumeFilePath == nil || *profile.ResumeFilePath == "" {
		logger.Warn("getResumeFilePath: no resume found")
		return "", apperrors.NotFound("resume not found")
	}

	logger.Debug("getResumeFilePath: resume file path retrieved", "path", *profile.ResumeFilePath)
	return *profile.ResumeFilePath, nil
}
