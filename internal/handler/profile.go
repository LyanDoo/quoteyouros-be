package handler

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/quoteyouros/backend/internal/domain"
	profileusecase "github.com/quoteyouros/backend/internal/usecase/profile"
	apperrors "github.com/quoteyouros/backend/pkg/errors"
	"github.com/quoteyouros/backend/pkg/logger"
	"github.com/quoteyouros/backend/pkg/response"
)

type ProfileHandler struct {
	usecase *profileusecase.ProfileUseCase
}

// NewProfileHandler creates a new profile handler
func NewProfileHandler(usecase *profileusecase.ProfileUseCase) *ProfileHandler {
	return &ProfileHandler{usecase: usecase}
}

// GetProfile retrieves profile information
// GET /api/profile/about
func (h *ProfileHandler) GetProfile(c *fiber.Ctx) error {
	logger.Debug("getProfile: retrieving profile")
	profile, err := h.usecase.GetProfile(c.Context())
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			logger.Warn("getProfile: failed to retrieve profile", "error", appErr.Details)
			return response.ErrorResponseJSON(c, appErr.Code, appErr.Details, appErr.Code)
		}
		logger.Error("getProfile: unexpected error", "error", err.Error())
		return response.ErrorResponseJSON(c, fiber.StatusInternalServerError, "failed to retrieve profile", fiber.StatusInternalServerError)
	}

	// Build response
	profileResp := domain.ProfileResponse{}
	profileResp.About = profile.About
	profileResp.Resume.HasResume = profile.ResumeFilePath != ""
	if profileResp.Resume.HasResume {
		profileResp.Resume.FileName = profile.ResumeFileName
		profileResp.Resume.FileSize = profile.ResumeFileSize
		if profile.ResumeUploadedAt != nil {
			profileResp.Resume.UploadedAt = profile.ResumeUploadedAt.Format("2006-01-02 15:04:05")
		}
	}

	logger.Info("getProfile: profile retrieved successfully")
	return response.SuccessResponse(c, fiber.StatusOK, profileResp, "Profile retrieved successfully")
}

// UpdateProfile updates profile information
// PUT /api/profile/about
func (h *ProfileHandler) UpdateProfile(c *fiber.Ctx) error {
	var req domain.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Error("updateProfile: failed to parse request body", "error", err.Error())
		return response.ErrorResponseJSON(c, fiber.StatusBadRequest, "invalid request body", fiber.StatusBadRequest)
	}

	logger.Info("updateProfile: attempting to update profile", "about_length", len(req.About))
	profile, err := h.usecase.UpdateProfile(c.Context(), &req)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			logger.Warn("updateProfile: failed to update profile", "error", appErr.Details)
			return response.ErrorResponseJSON(c, appErr.Code, appErr.Details, appErr.Code)
		}
		logger.Error("updateProfile: unexpected error", "error", err.Error())
		return response.ErrorResponseJSON(c, fiber.StatusInternalServerError, "failed to update profile", fiber.StatusInternalServerError)
	}

	// Build response
	profileResp := domain.ProfileResponse{}
	profileResp.About = profile.About
	profileResp.Resume.HasResume = profile.ResumeFilePath != ""
	if profileResp.Resume.HasResume {
		profileResp.Resume.FileName = profile.ResumeFileName
		profileResp.Resume.FileSize = profile.ResumeFileSize
		if profile.ResumeUploadedAt != nil {
			profileResp.Resume.UploadedAt = profile.ResumeUploadedAt.Format("2006-01-02 15:04:05")
		}
	}

	logger.Info("updateProfile: profile updated successfully")
	return response.SuccessResponse(c, fiber.StatusOK, profileResp, "Profile updated successfully")
}

// UploadResume uploads or replaces resume file
// POST /api/profile/resume
func (h *ProfileHandler) UploadResume(c *fiber.Ctx) error {
	logger.Info("uploadResume: attempting to upload resume")

	// Parse multipart form
	file, err := c.FormFile("file")
	if err != nil {
		logger.Error("uploadResume: failed to parse file from form", "error", err.Error())
		return response.ErrorResponseJSON(c, fiber.StatusBadRequest, "file is required", fiber.StatusBadRequest)
	}

	logger.Debug("uploadResume: file received", "filename", file.Filename, "size", file.Size)
	profile, err := h.usecase.UploadResume(c.Context(), file)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			logger.Warn("uploadResume: failed to upload resume", "error", appErr.Details)
			return response.ErrorResponseJSON(c, appErr.Code, appErr.Details, appErr.Code)
		}
		logger.Error("uploadResume: unexpected error", "error", err.Error())
		return response.ErrorResponseJSON(c, fiber.StatusInternalServerError, "failed to upload resume", fiber.StatusInternalServerError)
	}

	// Build response
	uploadResp := domain.ResumeUploadResponse{
		FileName:    profile.ResumeFileName,
		FileSize:    profile.ResumeFileSize,
		DownloadURL: "/api/profile/resume/download",
	}
	if profile.ResumeUploadedAt != nil {
		uploadResp.UploadedAt = profile.ResumeUploadedAt.Format("2006-01-02 15:04:05")
	}

	logger.Info("uploadResume: resume uploaded successfully", "filename", profile.ResumeFileName)
	return response.SuccessResponse(c, fiber.StatusCreated, uploadResp, "Resume uploaded successfully")
}

// DownloadResume downloads the resume file
// GET /api/profile/resume/download
func (h *ProfileHandler) DownloadResume(c *fiber.Ctx) error {
	logger.Debug("downloadResume: attempting to download resume")

	filePath, err := h.usecase.GetResumeFilePath(c.Context())
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			logger.Warn("downloadResume: failed to get resume", "error", appErr.Details)
			return response.ErrorResponseJSON(c, appErr.Code, appErr.Details, appErr.Code)
		}
		logger.Error("downloadResume: unexpected error", "error", err.Error())
		return response.ErrorResponseJSON(c, fiber.StatusInternalServerError, "failed to download resume", fiber.StatusInternalServerError)
	}

	// Check if file exists
	if _, err := os.Stat(filePath); err != nil {
		logger.Warn("downloadResume: resume file not found", "path", filePath)
		return response.ErrorResponseJSON(c, fiber.StatusNotFound, "resume file not found", fiber.StatusNotFound)
	}

	logger.Info("downloadResume: sending resume file", "path", filePath)
	return c.Download(filePath, "resume.pdf")
}
