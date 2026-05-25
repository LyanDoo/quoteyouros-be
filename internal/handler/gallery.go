package handler

import (
	"mime/multipart"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/quoteyouros/backend/internal/domain"
	galleryusecase "github.com/quoteyouros/backend/internal/usecase/gallery"
	apperrors "github.com/quoteyouros/backend/pkg/errors"
	"github.com/quoteyouros/backend/pkg/logger"
	"github.com/quoteyouros/backend/pkg/response"
)

type GalleryHandler struct {
	usecase *galleryusecase.GalleryUseCase
}

// NewGalleryHandler creates a new gallery handler
func NewGalleryHandler(usecase *galleryusecase.GalleryUseCase) *GalleryHandler {
	return &GalleryHandler{usecase: usecase}
}

// GetAllGalleryItems retrieves all gallery items
// GET /api/gallery
func (h *GalleryHandler) GetAllGalleryItems(c *fiber.Ctx) error {
	logger.Debug("getAllGalleryItems: retrieving all items")
	items, err := h.usecase.GetAllGalleryItems(c.Context())
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			logger.Warn("getAllGalleryItems: failed to retrieve items", "error", appErr.Details)
			return response.ErrorResponseJSON(c, appErr.Code, appErr.Details, appErr.Code)
		}
		logger.Error("getAllGalleryItems: unexpected error", "error", err.Error())
		return response.ErrorResponseJSON(c, fiber.StatusInternalServerError, "failed to retrieve gallery items", fiber.StatusInternalServerError)
	}

	// Map to responses
	resp := make([]domain.GalleryItemResponse, len(items))
	for i, item := range items {
		resp[i] = domain.GalleryItemResponse{
			ID:          item.ID,
			Title:       item.Title,
			Description: item.Description,
			Author:      item.Author,
			Image:       "/api/gallery/images/" + item.ImageFileName,
			CreatedAt:   item.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   item.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	logger.Info("getAllGalleryItems: items retrieved successfully", "count", len(resp))
	return response.SuccessResponse(c, fiber.StatusOK, resp, "Gallery items retrieved successfully")
}

// GetGalleryItem retrieves a single gallery item by ID
// GET /api/gallery/:id
func (h *GalleryHandler) GetGalleryItem(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		logger.Warn("getGalleryItem: missing item ID")
		return response.ErrorResponseJSON(c, fiber.StatusBadRequest, "gallery item ID is required", fiber.StatusBadRequest)
	}

	logger.Debug("getGalleryItem: retrieving item", "item_id", id)
	item, err := h.usecase.GetGalleryItem(c.Context(), id)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			logger.Warn("getGalleryItem: failed to retrieve item", "item_id", id, "error", appErr.Details)
			return response.ErrorResponseJSON(c, appErr.Code, appErr.Details, appErr.Code)
		}
		logger.Error("getGalleryItem: unexpected error", "item_id", id, "error", err.Error())
		return response.ErrorResponseJSON(c, fiber.StatusInternalServerError, "failed to retrieve gallery item", fiber.StatusInternalServerError)
	}

	resp := domain.GalleryItemResponse{
		ID:          item.ID,
		Title:       item.Title,
		Description: item.Description,
		Author:      item.Author,
		Image:       "/api/gallery/images/" + item.ImageFileName,
		CreatedAt:   item.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   item.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	logger.Debug("getGalleryItem: item retrieved successfully", "item_id", id)
	return response.SuccessResponse(c, fiber.StatusOK, resp, "Gallery item retrieved successfully")
}

// CreateGalleryItem creates a new gallery item (multipart/form-data)
// POST /api/gallery
func (h *GalleryHandler) CreateGalleryItem(c *fiber.Ctx) error {
	logger.Info("createGalleryItem: attempting to create gallery item")

	title := c.FormValue("title")
	if title == "" {
		title = c.FormValue("Title")
	}

	description := c.FormValue("description")
	if description == "" {
		description = c.FormValue("Description")
	}

	author := c.FormValue("author")
	if author == "" {
		author = c.FormValue("Author")
	}

	// Parse form file
	file, err := c.FormFile("image")
	if err != nil {
		file, err = c.FormFile("Image")
	}
	if err != nil {
		logger.Error("createGalleryItem: failed to parse image file from form", "error", err.Error())
		return response.ErrorResponseJSON(c, fiber.StatusBadRequest, "image file is required", fiber.StatusBadRequest)
	}

	item, err := h.usecase.CreateGalleryItem(c.Context(), title, description, author, file)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			logger.Warn("createGalleryItem: failed to create item", "error", appErr.Details)
			return response.ErrorResponseJSON(c, appErr.Code, appErr.Details, appErr.Code)
		}
		logger.Error("createGalleryItem: unexpected error", "error", err.Error())
		return response.ErrorResponseJSON(c, fiber.StatusInternalServerError, "failed to create gallery item", fiber.StatusInternalServerError)
	}

	resp := domain.GalleryItemResponse{
		ID:          item.ID,
		Title:       item.Title,
		Description: item.Description,
		Author:      item.Author,
		Image:       "/api/gallery/images/" + item.ImageFileName,
		CreatedAt:   item.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   item.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	logger.Info("createGalleryItem: gallery item created successfully", "item_id", item.ID)
	return response.SuccessResponse(c, fiber.StatusCreated, resp, "Gallery item created successfully")
}

// UpdateGalleryItem updates an existing gallery item metadata and optional image replacement (multipart/form-data)
// PUT /api/gallery/:id
func (h *GalleryHandler) UpdateGalleryItem(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		logger.Warn("updateGalleryItem: missing item ID")
		return response.ErrorResponseJSON(c, fiber.StatusBadRequest, "gallery item ID is required", fiber.StatusBadRequest)
	}

	title := c.FormValue("title")
	if title == "" {
		title = c.FormValue("Title")
	}

	description := c.FormValue("description")
	if description == "" {
		description = c.FormValue("Description")
	}

	author := c.FormValue("author")
	if author == "" {
		author = c.FormValue("Author")
	}

	// Parse optional form file
	var file *multipart.FileHeader
	form, err := c.MultipartForm()
	if err == nil && form != nil {
		if headers, ok := form.File["image"]; ok && len(headers) > 0 {
			file = headers[0]
		} else if headers, ok := form.File["Image"]; ok && len(headers) > 0 {
			file = headers[0]
		}
	}

	logger.Info("updateGalleryItem: attempting to update gallery item", "item_id", id)
	item, err := h.usecase.UpdateGalleryItem(c.Context(), id, title, description, author, file)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			logger.Warn("updateGalleryItem: failed to update item", "item_id", id, "error", appErr.Details)
			return response.ErrorResponseJSON(c, appErr.Code, appErr.Details, appErr.Code)
		}
		logger.Error("updateGalleryItem: unexpected error", "item_id", id, "error", err.Error())
		return response.ErrorResponseJSON(c, fiber.StatusInternalServerError, "failed to update gallery item", fiber.StatusInternalServerError)
	}

	resp := domain.GalleryItemResponse{
		ID:          item.ID,
		Title:       item.Title,
		Description: item.Description,
		Author:      item.Author,
		Image:       "/api/gallery/images/" + item.ImageFileName,
		CreatedAt:   item.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   item.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	logger.Info("updateGalleryItem: gallery item updated successfully", "item_id", id)
	return response.SuccessResponse(c, fiber.StatusOK, resp, "Gallery item updated successfully")
}

// DeleteGalleryItem deletes a gallery item and its image file
// DELETE /api/gallery/:id
func (h *GalleryHandler) DeleteGalleryItem(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		logger.Warn("deleteGalleryItem: missing item ID")
		return response.ErrorResponseJSON(c, fiber.StatusBadRequest, "gallery item ID is required", fiber.StatusBadRequest)
	}

	logger.Info("deleteGalleryItem: attempting to delete gallery item", "item_id", id)
	err := h.usecase.DeleteGalleryItem(c.Context(), id)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			logger.Warn("deleteGalleryItem: failed to delete item", "item_id", id, "error", appErr.Details)
			return response.ErrorResponseJSON(c, appErr.Code, appErr.Details, appErr.Code)
		}
		logger.Error("deleteGalleryItem: unexpected error", "item_id", id, "error", err.Error())
		return response.ErrorResponseJSON(c, fiber.StatusInternalServerError, "failed to delete gallery item", fiber.StatusInternalServerError)
	}

	logger.Info("deleteGalleryItem: gallery item deleted successfully", "item_id", id)
	return response.SuccessResponse(c, fiber.StatusOK, fiber.Map{}, "Gallery item deleted successfully")
}

// GetImage serves the raw image file for viewing in web browser
// GET /api/gallery/images/:filename
func (h *GalleryHandler) GetImage(c *fiber.Ctx) error {
	filename := c.Params("filename")
	if filename == "" {
		logger.Warn("getImage: missing filename")
		return response.ErrorResponseJSON(c, fiber.StatusBadRequest, "filename is required", fiber.StatusBadRequest)
	}

	// We look for files in GalleryStoragePath
	filePath := "./storage/gallery/" + filename

	logger.Debug("getImage: attempting to serve image file", "path", filePath)

	// Check if file exists
	if _, err := os.Stat(filePath); err != nil {
		logger.Warn("getImage: image file not found", "path", filePath)
		return response.ErrorResponseJSON(c, fiber.StatusNotFound, "image file not found", fiber.StatusNotFound)
	}

	logger.Info("getImage: sending image file", "path", filePath)
	return c.SendFile(filePath)
}
