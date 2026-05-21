package fileupload

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/quoteyouros/backend/pkg/logger"
)

const (
	// MaxResumeSize is the maximum allowed file size (10MB)
	MaxResumeSize = 10 * 1024 * 1024 // 10MB
	// ResumeStoragePath is the directory where resumes are stored
	ResumeStoragePath = "./storage/resumes"
	// AllowedMimeType is the only allowed MIME type
	AllowedMimeType = "application/pdf"
)

type FileUploadService struct {
	storagePath string
}

// New creates a new file upload service
func New(storagePath string) *FileUploadService {
	return &FileUploadService{
		storagePath: storagePath,
	}
}

// UploadResume handles resume file upload
func (s *FileUploadService) UploadResume(file *multipart.FileHeader) (fileName string, filePath string, fileSize int64, mimeType string, err error) {
	logger.Debug("uploadResume: starting resume upload", "filename", file.Filename, "size", file.Size)

	// Validate file size
	if file.Size > MaxResumeSize {
		logger.Warn("uploadResume: file too large", "filename", file.Filename, "size", file.Size, "max", MaxResumeSize)
		return "", "", 0, "", fmt.Errorf("file size exceeds maximum limit of %d bytes", MaxResumeSize)
	}

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".pdf" {
		logger.Warn("uploadResume: invalid file type", "filename", file.Filename, "extension", ext)
		return "", "", 0, "", fmt.Errorf("only PDF files are allowed")
	}

	// Open uploaded file
	src, err := file.Open()
	if err != nil {
		logger.Error("uploadResume: failed to open uploaded file", "filename", file.Filename, "error", err.Error())
		return "", "", 0, "", fmt.Errorf("failed to open file: %v", err)
	}
	defer src.Close()

	// Read first 1KB to validate MIME type
	header := make([]byte, 1024)
	n, err := src.Read(header)
	if err != nil && err != io.EOF {
		logger.Error("uploadResume: failed to read file header", "filename", file.Filename, "error", err.Error())
		return "", "", 0, "", fmt.Errorf("failed to read file: %v", err)
	}

	// Check PDF magic number (PDF files start with %PDF)
	if !strings.HasPrefix(string(header[:n]), "%PDF") {
		logger.Warn("uploadResume: file is not a valid PDF", "filename", file.Filename)
		return "", "", 0, "", fmt.Errorf("file is not a valid PDF")
	}

	// Reset file pointer to beginning
	src.Seek(0, 0)

	// Ensure storage directory exists
	if err := os.MkdirAll(s.storagePath, 0755); err != nil {
		logger.Error("uploadResume: failed to create storage directory", "path", s.storagePath, "error", err.Error())
		return "", "", 0, "", fmt.Errorf("failed to create storage directory: %v", err)
	}

	// Generate unique filename: timestamp_originalname
	timestamp := time.Now().Unix()
	originalName := strings.TrimSuffix(filepath.Base(file.Filename), ext)
	generatedFileName := fmt.Sprintf("%d_%s%s", timestamp, originalName, ext)
	fullPath := filepath.Join(s.storagePath, generatedFileName)

	// Create destination file
	dst, err := os.Create(fullPath)
	if err != nil {
		logger.Error("uploadResume: failed to create destination file", "path", fullPath, "error", err.Error())
		return "", "", 0, "", fmt.Errorf("failed to create file: %v", err)
	}
	defer dst.Close()

	// Copy file content
	fileSize, err = io.Copy(dst, src)
	if err != nil {
		logger.Error("uploadResume: failed to copy file content", "filename", file.Filename, "error", err.Error())
		os.Remove(fullPath) // Clean up on error
		return "", "", 0, "", fmt.Errorf("failed to save file: %v", err)
	}

	logger.Info("uploadResume: file uploaded successfully", "filename", generatedFileName, "size", fileSize, "path", fullPath)
	return generatedFileName, fullPath, fileSize, AllowedMimeType, nil
}

// DeleteResume deletes resume file from storage
func (s *FileUploadService) DeleteResume(filePath string) error {
	if filePath == "" {
		logger.Debug("deleteResume: no file path provided")
		return nil
	}

	logger.Debug("deleteResume: deleting file", "path", filePath)
	if err := os.Remove(filePath); err != nil {
		if os.IsNotExist(err) {
			logger.Warn("deleteResume: file not found", "path", filePath)
			return nil // Not an error if file doesn't exist
		}
		logger.Error("deleteResume: failed to delete file", "path", filePath, "error", err.Error())
		return err
	}

	logger.Info("deleteResume: file deleted successfully", "path", filePath)
	return nil
}

// ValidateResume validates a resume file without saving it
func (s *FileUploadService) ValidateResume(file *multipart.FileHeader) error {
	if file.Size > MaxResumeSize {
		return fmt.Errorf("file size exceeds maximum limit of %d bytes", MaxResumeSize)
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".pdf" {
		return fmt.Errorf("only PDF files are allowed")
	}

	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer src.Close()

	header := make([]byte, 1024)
	n, err := src.Read(header)
	if err != nil && err != io.EOF {
		return fmt.Errorf("failed to read file: %v", err)
	}

	if !strings.HasPrefix(string(header[:n]), "%PDF") {
		return fmt.Errorf("file is not a valid PDF")
	}

	return nil
}
