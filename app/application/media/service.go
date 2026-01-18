package media

import (
	"cacto-cms/app/domain/media"
	"cacto-cms/app/shared/errors"
	"mime"
	"path/filepath"
	"strings"
	"time"
)

// Service handles business logic for media
type Service struct {
	repo media.Repository
}

// NewService creates a new media service
func NewService(repo media.Repository) *Service {
	return &Service{repo: repo}
}

// GetMediaByID retrieves a media by ID
func (s *Service) GetMediaByID(id int) (*media.Media, error) {
	m, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrCodeNotFound, "Media not found", 404)
	}
	return m, nil
}

// GetAllMedia retrieves all media with pagination
func (s *Service) GetAllMedia(limit, offset int) ([]*media.Media, error) {
	return s.repo.FindAll(limit, offset)
}

// CreateMedia creates a new media record
func (s *Service) CreateMedia(filename, originalName, mimeType string, size int64, path, url string) (*media.Media, error) {
	// Validate mime type
	if mimeType == "" {
		ext := filepath.Ext(originalName)
		mimeType = mime.TypeByExtension(ext)
		if mimeType == "" {
			mimeType = "application/octet-stream"
		}
	}

	// Create media entity
	m := &media.Media{
		Filename:     filename,
		OriginalName: originalName,
		MimeType:     mimeType,
		Size:         size,
		Path:         path,
		URL:          url,
		CreatedAt:    time.Now(),
	}

	if err := s.repo.Create(m); err != nil {
		return nil, errors.NewInternal("Failed to create media", err)
	}

	return m, nil
}

// UpdateMedia updates media metadata
func (s *Service) UpdateMedia(m *media.Media) error {
	return s.repo.Update(m)
}

// DeleteMedia deletes a media by ID
func (s *Service) DeleteMedia(id int) error {
	return s.repo.Delete(id)
}

// ValidateFileType validates if file type is allowed
func (s *Service) ValidateFileType(mimeType string) bool {
	allowedTypes := []string{
		"image/jpeg", "image/png", "image/gif", "image/webp", "image/svg+xml",
		"video/mp4", "video/webm", "video/ogg",
		"application/pdf",
		"text/plain", "text/csv",
	}

	mimeType = strings.ToLower(mimeType)
	for _, allowed := range allowedTypes {
		if mimeType == allowed {
			return true
		}
	}

	return false
}

// ValidateFileSize validates file size (max 10MB)
func (s *Service) ValidateFileSize(size int64) bool {
	const maxSize = 10 * 1024 * 1024 // 10MB
	return size > 0 && size <= maxSize
}
