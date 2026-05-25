package gallery_test

import (
	"bytes"
	"context"
	"errors"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/quoteyouros/backend/internal/domain"
	"github.com/quoteyouros/backend/internal/usecase/gallery"
	apperrors "github.com/quoteyouros/backend/pkg/errors"
	"github.com/quoteyouros/backend/pkg/fileupload"
)

type mockGalleryRepository struct {
	createFunc func(ctx context.Context, item *domain.GalleryItem) error
	getFunc    func(ctx context.Context, id string) (*domain.GalleryItem, error)
	getAllFunc func(ctx context.Context) ([]*domain.GalleryItem, error)
	updateFunc func(ctx context.Context, id string, item *domain.GalleryItem) error
	deleteFunc func(ctx context.Context, id string) error
}

func (m *mockGalleryRepository) CreateGalleryItem(ctx context.Context, item *domain.GalleryItem) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, item)
	}
	return nil
}

func (m *mockGalleryRepository) GetGalleryItem(ctx context.Context, id string) (*domain.GalleryItem, error) {
	if m.getFunc != nil {
		return m.getFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockGalleryRepository) GetAllGalleryItems(ctx context.Context) ([]*domain.GalleryItem, error) {
	if m.getAllFunc != nil {
		return m.getAllFunc(ctx)
	}
	return nil, nil
}

func (m *mockGalleryRepository) UpdateGalleryItem(ctx context.Context, id string, item *domain.GalleryItem) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, id, item)
	}
	return nil
}

func (m *mockGalleryRepository) DeleteGalleryItem(ctx context.Context, id string) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return nil
}

// newTestImageFileHeader helper to construct a valid multipart.FileHeader with a real backing file
func newTestImageFileHeader(t *testing.T, filename string, content []byte) *multipart.FileHeader {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile("image", filename)
	if err != nil {
		t.Fatalf("failed to create form file: %v", err)
	}
	_, err = part.Write(content)
	if err != nil {
		t.Fatalf("failed to write part: %v", err)
	}
	err = writer.Close()
	if err != nil {
		t.Fatalf("failed to close writer: %v", err)
	}

	reader := multipart.NewReader(&buf, writer.Boundary())
	form, err := reader.ReadForm(10 * 1024 * 1024)
	if err != nil {
		t.Fatalf("failed to read form: %v", err)
	}

	files := form.File["image"]
	if len(files) == 0 {
		t.Fatalf("no files found in form")
	}
	return files[0]
}

func TestCreateGalleryItem_Success(t *testing.T) {
	tempDir := t.TempDir()
	fileUpload := fileupload.New(tempDir)

	var savedItem *domain.GalleryItem
	repo := &mockGalleryRepository{
		createFunc: func(ctx context.Context, item *domain.GalleryItem) error {
			savedItem = item
			return nil
		},
	}

	uc := gallery.New(repo, fileUpload)

	// Valid PNG file content
	pngHeader := []byte("\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR\x00\x00\x00\x01\x00\x00\x00\x01\x08\x06\x00\x00\x00\x1f\x15\xc4\x89")
	file := newTestImageFileHeader(t, "test.png", pngHeader)

	item, err := uc.CreateGalleryItem(context.Background(), "Awesome NFT", "NFT Description", file)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if item == nil {
		t.Fatal("expected item to be returned")
	}

	if item.Title != "Awesome NFT" {
		t.Errorf("expected title Awesome NFT, got %s", item.Title)
	}

	if item.Description != "NFT Description" {
		t.Errorf("expected description NFT Description, got %s", item.Description)
	}

	if item.ImageFileName == "" {
		t.Error("expected generated filename to be set")
	}

	if !strings.HasSuffix(item.ImageFilePath, item.ImageFileName) {
		t.Errorf("expected file path to end with filename, got path %s filename %s", item.ImageFilePath, item.ImageFileName)
	}

	// Verify file was written to disk
	if _, err := os.Stat(item.ImageFilePath); os.IsNotExist(err) {
		t.Errorf("expected file to exist at path %s", item.ImageFilePath)
	}

	if savedItem == nil {
		t.Fatal("expected item to be saved in repository")
	}
}

func TestCreateGalleryItem_ValidationFailures(t *testing.T) {
	tempDir := t.TempDir()
	fileUpload := fileupload.New(tempDir)
	repo := &mockGalleryRepository{}
	uc := gallery.New(repo, fileUpload)

	pngHeader := []byte("\x89PNG\r\n\x1a\n")
	file := newTestImageFileHeader(t, "test.png", pngHeader)

	tests := []struct {
		name        string
		title       string
		description string
		file        *multipart.FileHeader
		wantErr     string
	}{
		{
			name:        "empty title",
			title:       "",
			description: "Some description",
			file:        file,
			wantErr:     "title is required",
		},
		{
			name:        "empty description",
			title:       "Valid Title",
			description: "",
			file:        file,
			wantErr:     "description is required",
		},
		{
			name:        "nil file",
			title:       "Valid Title",
			description: "Some description",
			file:        nil,
			wantErr:     "image file is required",
		},
		{
			name:        "invalid file extension",
			title:       "Valid Title",
			description: "Some description",
			file:        newTestImageFileHeader(t, "test.pdf", []byte("%PDF-1.4")),
			wantErr:     "only JPEG, PNG, GIF, and WebP images are allowed",
		},
		{
			name:        "invalid file content",
			title:       "Valid Title",
			description: "Some description",
			file:        newTestImageFileHeader(t, "test.png", []byte("plain text content")),
			wantErr:     "file content is not a valid image",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := uc.CreateGalleryItem(context.Background(), tt.title, tt.description, tt.file)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			appErr, ok := err.(*apperrors.AppError)
			if !ok {
				t.Fatalf("expected AppError, got %T: %v", err, err)
			}
			if !strings.Contains(appErr.Details, tt.wantErr) {
				t.Errorf("expected error details containing %q, got %q", tt.wantErr, appErr.Details)
			}
		})
	}
}

func TestCreateGalleryItem_DatabaseFailureCleanup(t *testing.T) {
	tempDir := t.TempDir()
	fileUpload := fileupload.New(tempDir)

	repo := &mockGalleryRepository{
		createFunc: func(ctx context.Context, item *domain.GalleryItem) error {
			return errors.New("database connection lost")
		},
	}

	uc := gallery.New(repo, fileUpload)

	pngHeader := []byte("\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR\x00\x00\x00\x01\x00\x00\x00\x01\x08\x06\x00\x00\x00\x1f\x15\xc4\x89")
	file := newTestImageFileHeader(t, "test.png", pngHeader)

	_, err := uc.CreateGalleryItem(context.Background(), "Awesome NFT", "NFT Description", file)
	if err == nil {
		t.Fatal("expected database error, got nil")
	}

	// Verify no files are left in the temp directory (should be cleaned up)
	entries, err := os.ReadDir(tempDir)
	if err != nil {
		t.Fatalf("failed to read temp dir: %v", err)
	}
	if len(entries) > 0 {
		t.Errorf("expected temp directory to be clean, but found %d files", len(entries))
	}
}

func TestGetGalleryItem_SuccessAndNotFound(t *testing.T) {
	tempDir := t.TempDir()
	fileUpload := fileupload.New(tempDir)

	mockItem := &domain.GalleryItem{
		ID:            "item-123",
		Title:         "Sample NFT",
		Description:   "Description",
		ImageFileName: "img.png",
		ImageFilePath: "/some/path/img.png",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	repo := &mockGalleryRepository{
		getFunc: func(ctx context.Context, id string) (*domain.GalleryItem, error) {
			if id == "item-123" {
				return mockItem, nil
			}
			return nil, errors.New("no rows in result set")
		},
	}

	uc := gallery.New(repo, fileUpload)

	t.Run("success", func(t *testing.T) {
		item, err := uc.GetGalleryItem(context.Background(), "item-123")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if item.ID != "item-123" {
			t.Errorf("expected ID item-123, got %s", item.ID)
		}
	})

	t.Run("not found", func(t *testing.T) {
		_, err := uc.GetGalleryItem(context.Background(), "non-existent")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		appErr, ok := err.(*apperrors.AppError)
		if !ok {
			t.Fatalf("expected AppError, got %T: %v", err, err)
		}
		if appErr.Code != 404 {
			t.Errorf("expected status code 404, got %d", appErr.Code)
		}
	})
}

func TestUpdateGalleryItem_SuccessWithoutFile(t *testing.T) {
	tempDir := t.TempDir()
	fileUpload := fileupload.New(tempDir)

	existingItem := &domain.GalleryItem{
		ID:            "item-123",
		Title:         "Old Title",
		Description:   "Old Description",
		ImageFileName: "old.png",
		ImageFilePath: "/some/path/old.png",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	var updatedItem *domain.GalleryItem
	repo := &mockGalleryRepository{
		getFunc: func(ctx context.Context, id string) (*domain.GalleryItem, error) {
			if id == "item-123" {
				return existingItem, nil
			}
			return nil, errors.New("no rows in result set")
		},
		updateFunc: func(ctx context.Context, id string, item *domain.GalleryItem) error {
			updatedItem = item
			return nil
		},
	}

	uc := gallery.New(repo, fileUpload)

	item, err := uc.UpdateGalleryItem(context.Background(), "item-123", "New Title", "New Description", nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if item.Title != "New Title" || item.Description != "New Description" {
		t.Errorf("fields did not update properly: %+v", item)
	}

	if item.ImageFileName != "old.png" {
		t.Error("image file name should not have changed")
	}

	if updatedItem == nil {
		t.Fatal("expected update in repo to be called")
	}
}

func TestUpdateGalleryItem_SuccessWithFileReplacement(t *testing.T) {
	tempDir := t.TempDir()
	fileUpload := fileupload.New(tempDir)

	// Create a physical old file to be deleted
	oldFileName := "old_image.png"
	oldFilePath := filepath.Join(tempDir, oldFileName)
	err := os.WriteFile(oldFilePath, []byte("old content"), 0644)
	if err != nil {
		t.Fatalf("failed to create old file: %v", err)
	}

	existingItem := &domain.GalleryItem{
		ID:            "item-123",
		Title:         "Old Title",
		Description:   "Old Description",
		ImageFileName: oldFileName,
		ImageFilePath: oldFilePath,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	repo := &mockGalleryRepository{
		getFunc: func(ctx context.Context, id string) (*domain.GalleryItem, error) {
			return existingItem, nil
		},
		updateFunc: func(ctx context.Context, id string, item *domain.GalleryItem) error {
			return nil
		},
	}

	uc := gallery.New(repo, fileUpload)

	pngHeader := []byte("\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR\x00\x00\x00\x01\x00\x00\x00\x01\x08\x06\x00\x00\x00\x1f\x15\xc4\x89")
	newFile := newTestImageFileHeader(t, "new_image.png", pngHeader)

	item, err := uc.UpdateGalleryItem(context.Background(), "item-123", "New Title", "New Description", newFile)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify old file was deleted from disk
	if _, err := os.Stat(oldFilePath); !os.IsNotExist(err) {
		t.Errorf("expected old file %s to be deleted from disk", oldFilePath)
	}

	// Verify new file exists on disk
	if _, err := os.Stat(item.ImageFilePath); os.IsNotExist(err) {
		t.Errorf("expected new file %s to exist on disk", item.ImageFilePath)
	}
}

func TestDeleteGalleryItem_Success(t *testing.T) {
	tempDir := t.TempDir()
	fileUpload := fileupload.New(tempDir)

	// Create physical image file
	fileName := "to_delete.png"
	filePath := filepath.Join(tempDir, fileName)
	err := os.WriteFile(filePath, []byte("some image content"), 0644)
	if err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	existingItem := &domain.GalleryItem{
		ID:            "item-123",
		Title:         "Delete Me",
		Description:   "Desc",
		ImageFileName: fileName,
		ImageFilePath: filePath,
	}

	var deletedID string
	repo := &mockGalleryRepository{
		getFunc: func(ctx context.Context, id string) (*domain.GalleryItem, error) {
			if id == "item-123" {
				return existingItem, nil
			}
			return nil, errors.New("no rows in result set")
		},
		deleteFunc: func(ctx context.Context, id string) error {
			deletedID = id
			return nil
		},
	}

	uc := gallery.New(repo, fileUpload)

	err = uc.DeleteGalleryItem(context.Background(), "item-123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if deletedID != "item-123" {
		t.Errorf("expected deleted ID item-123, got %s", deletedID)
	}

	// Verify physical file was deleted from disk
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		t.Errorf("expected image file %s to be deleted from disk", filePath)
	}
}
