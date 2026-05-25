package gallery

import (
	"context"
	"mime/multipart"
	"time"

	"github.com/quoteyouros/backend/internal/domain"
	apperrors "github.com/quoteyouros/backend/pkg/errors"
	"github.com/quoteyouros/backend/pkg/fileupload"
	"github.com/quoteyouros/backend/pkg/logger"
)

type GalleryUseCase struct {
	galleryRepo domain.GalleryRepository
	fileUpload  *fileupload.FileUploadService
}

// New creates a new gallery use case
func New(galleryRepo domain.GalleryRepository, fileUpload *fileupload.FileUploadService) *GalleryUseCase {
	return &GalleryUseCase{
		galleryRepo: galleryRepo,
		fileUpload:  fileUpload,
	}
}

// CreateGalleryItem uploads an image and saves gallery item metadata
func (u *GalleryUseCase) CreateGalleryItem(ctx context.Context, title, description string, file *multipart.FileHeader) (*domain.GalleryItem, error) {
	logger.Debug("createGalleryItem: validating request", "title", title)

	if title == "" {
		return nil, apperrors.BadRequest("title is required")
	}
	if description == "" {
		return nil, apperrors.BadRequest("description is required")
	}
	if file == nil {
		return nil, apperrors.BadRequest("image file is required")
	}

	// Validate file upload first
	_, err := u.fileUpload.ValidateImage(file)
	if err != nil {
		logger.Warn("createGalleryItem: file validation failed", "error", err.Error())
		return nil, apperrors.BadRequest(err.Error())
	}

	// Upload file
	fileName, filePath, fileSize, mimeType, err := u.fileUpload.UploadImage(file)
	if err != nil {
		logger.Error("createGalleryItem: failed to upload image", "error", err.Error())
		return nil, apperrors.BadRequest(err.Error())
	}

	// Build gallery item entity
	item := domain.NewGalleryItem(title, description, fileName, filePath, fileSize, mimeType)

	logger.Info("createGalleryItem: creating gallery item in database", "item_id", item.ID, "title", item.Title)
	if err := u.galleryRepo.CreateGalleryItem(ctx, item); err != nil {
		logger.Error("createGalleryItem: failed to save item to database", "item_id", item.ID, "error", err.Error())
		// Clean up uploaded file since db save failed
		if cleanupErr := u.fileUpload.DeleteImage(filePath); cleanupErr != nil {
			logger.Error("createGalleryItem: failed to clean up uploaded file", "path", filePath, "error", cleanupErr.Error())
		}
		return nil, apperrors.InternalServerError("failed to save gallery item: " + err.Error())
	}

	logger.Info("createGalleryItem: gallery item created successfully", "item_id", item.ID)
	return item, nil
}

// GetGalleryItem retrieves a single gallery item
func (u *GalleryUseCase) GetGalleryItem(ctx context.Context, id string) (*domain.GalleryItem, error) {
	logger.Debug("getGalleryItem: retrieving item", "item_id", id)

	item, err := u.galleryRepo.GetGalleryItem(ctx, id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			logger.Warn("getGalleryItem: item not found", "item_id", id)
			return nil, apperrors.NotFound("gallery item not found")
		}
		logger.Error("getGalleryItem: failed to retrieve item", "item_id", id, "error", err.Error())
		return nil, apperrors.InternalServerError("failed to retrieve gallery item: " + err.Error())
	}

	logger.Debug("getGalleryItem: item retrieved successfully", "item_id", id)
	return item, nil
}

// GetAllGalleryItems retrieves all gallery items
func (u *GalleryUseCase) GetAllGalleryItems(ctx context.Context) ([]*domain.GalleryItem, error) {
	logger.Debug("getAllGalleryItems: retrieving all items")

	items, err := u.galleryRepo.GetAllGalleryItems(ctx)
	if err != nil {
		logger.Error("getAllGalleryItems: failed to retrieve items", "error", err.Error())
		return nil, apperrors.InternalServerError("failed to retrieve gallery items: " + err.Error())
	}

	logger.Info("getAllGalleryItems: items retrieved successfully", "count", len(items))
	return items, nil
}

// UpdateGalleryItem updates an existing gallery item metadata and optional image replacement
func (u *GalleryUseCase) UpdateGalleryItem(ctx context.Context, id string, title, description string, file *multipart.FileHeader) (*domain.GalleryItem, error) {
	logger.Debug("updateGalleryItem: validating request", "item_id", id)

	// Retrieve existing item
	item, err := u.galleryRepo.GetGalleryItem(ctx, id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			logger.Warn("updateGalleryItem: item not found", "item_id", id)
			return nil, apperrors.NotFound("gallery item not found")
		}
		logger.Error("updateGalleryItem: failed to retrieve item", "item_id", id, "error", err.Error())
		return nil, apperrors.InternalServerError("failed to update gallery item: " + err.Error())
	}

	var oldFilePath string
	var newUploaded bool
	var newFileName, newFilePath, newMimeType string
	var newFileSize int64

	// If new file is uploaded
	if file != nil {
		logger.Debug("updateGalleryItem: new file uploaded, validating", "filename", file.Filename)
		_, err := u.fileUpload.ValidateImage(file)
		if err != nil {
			logger.Warn("updateGalleryItem: new file validation failed", "error", err.Error())
			return nil, apperrors.BadRequest(err.Error())
		}

		// Save old file path for cleanup after successful db update
		oldFilePath = item.ImageFilePath

		// Upload new image
		newFileName, newFilePath, newFileSize, newMimeType, err = u.fileUpload.UploadImage(file)
		if err != nil {
			logger.Error("updateGalleryItem: failed to upload new image", "error", err.Error())
			return nil, apperrors.BadRequest(err.Error())
		}
		newUploaded = true
	}

	// Update fields
	if title != "" {
		item.Title = title
	}
	if description != "" {
		item.Description = description
	}
	if newUploaded {
		item.ImageFileName = newFileName
		item.ImageFilePath = newFilePath
		item.ImageFileSize = newFileSize
		item.ImageMimeType = newMimeType
	}
	item.UpdatedAt = time.Now()

	logger.Info("updateGalleryItem: updating gallery item in database", "item_id", id)
	if err := u.galleryRepo.UpdateGalleryItem(ctx, id, item); err != nil {
		logger.Error("updateGalleryItem: failed to update item in database", "item_id", id, "error", err.Error())
		// Clean up newly uploaded file if db update failed
		if newUploaded {
			if cleanupErr := u.fileUpload.DeleteImage(newFilePath); cleanupErr != nil {
				logger.Error("updateGalleryItem: failed to clean up new file", "path", newFilePath, "error", cleanupErr.Error())
			}
		}
		return nil, apperrors.InternalServerError("failed to update gallery item: " + err.Error())
	}

	// Delete old file if a new file was successfully set and updated
	if newUploaded && oldFilePath != "" {
		logger.Debug("updateGalleryItem: deleting old file", "path", oldFilePath)
		if cleanupErr := u.fileUpload.DeleteImage(oldFilePath); cleanupErr != nil {
			logger.Warn("updateGalleryItem: failed to delete old file", "path", oldFilePath, "error", cleanupErr.Error())
		}
	}

	logger.Info("updateGalleryItem: gallery item updated successfully", "item_id", id)
	return item, nil
}

// DeleteGalleryItem deletes a gallery item and its associated image file
func (u *GalleryUseCase) DeleteGalleryItem(ctx context.Context, id string) error {
	logger.Debug("deleteGalleryItem: validating existence", "item_id", id)

	// Retrieve existing item to get file path
	item, err := u.galleryRepo.GetGalleryItem(ctx, id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			logger.Warn("deleteGalleryItem: item not found", "item_id", id)
			return apperrors.NotFound("gallery item not found")
		}
		logger.Error("deleteGalleryItem: failed to retrieve item", "item_id", id, "error", err.Error())
		return apperrors.InternalServerError("failed to delete gallery item: " + err.Error())
	}

	logger.Info("deleteGalleryItem: deleting item from database", "item_id", id)
	if err := u.galleryRepo.DeleteGalleryItem(ctx, id); err != nil {
		logger.Error("deleteGalleryItem: failed to delete item from database", "item_id", id, "error", err.Error())
		return apperrors.InternalServerError("failed to delete gallery item: " + err.Error())
	}

	// Delete associated image file
	if item.ImageFilePath != "" {
		logger.Debug("deleteGalleryItem: deleting image file", "path", item.ImageFilePath)
		if cleanupErr := u.fileUpload.DeleteImage(item.ImageFilePath); cleanupErr != nil {
			logger.Warn("deleteGalleryItem: failed to delete image file", "path", item.ImageFilePath, "error", cleanupErr.Error())
		}
	}

	logger.Info("deleteGalleryItem: gallery item deleted successfully", "item_id", id)
	return nil
}
