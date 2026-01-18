package page

import (
	"fmt"
	"strings"
	"time"
	"unicode"

	"cacto-cms/app/domain/page"
)

// Service handles business logic for pages
type Service struct {
	repo page.Repository
}

// NewService creates a new page service
func NewService(repo page.Repository) *Service {
	return &Service{repo: repo}
}

// GetPageBySlug retrieves a page by its slug
func (s *Service) GetPageBySlug(slug string) (*page.Page, error) {
	p, err := s.repo.FindBySlug(slug)
	if err != nil {
		return nil, err
	}

	// Load components
	components, err := s.repo.GetComponents(p.ID)
	if err == nil {
		p.Components = components
	}

	return p, nil
}

// GetAllPages retrieves all pages
func (s *Service) GetAllPages() ([]*page.Page, error) {
	return s.repo.FindAll()
}

// GetPublishedPages retrieves only published pages
func (s *Service) GetPublishedPages() ([]*page.Page, error) {
	return s.repo.FindPublished()
}

// CreatePage creates a new page
func (s *Service) CreatePage(p *page.Page) error {
	// Auto-generate slug if empty
	if p.Slug == "" {
		p.Slug = GenerateSlug(p.Title)
	}

	// Set timestamps
	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now

	// Set default status
	if p.Status == "" {
		p.Status = page.StatusDraft
	}

	return s.repo.Create(p)
}

// UpdatePage updates an existing page
func (s *Service) UpdatePage(p *page.Page) error {
	p.UpdatedAt = time.Now()
	return s.repo.Update(p)
}

// DeletePage deletes a page by ID
func (s *Service) DeletePage(id int) error {
	return s.repo.Delete(id)
}

// GenerateSlug creates a URL-friendly slug from text
func GenerateSlug(text string) string {
	// Convert to lowercase
	slug := strings.ToLower(text)

	// Turkish character replacements
	replacements := map[rune]string{
		'ç': "c", 'ğ': "g", 'ı': "i", 'ö': "o", 'ş': "s", 'ü': "u",
		'Ç': "c", 'Ğ': "g", 'İ': "i", 'Ö': "o", 'Ş': "s", 'Ü': "u",
	}

	var result strings.Builder
	for _, r := range slug {
		if replacement, ok := replacements[r]; ok {
			result.WriteString(replacement)
		} else if unicode.IsLetter(r) || unicode.IsNumber(r) {
			result.WriteRune(r)
		} else if unicode.IsSpace(r) || r == '-' || r == '_' {
			result.WriteRune('-')
		}
	}

	// Clean up multiple dashes
	slug = result.String()
	slug = strings.Trim(slug, "-")
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}

	return slug
}

// ValidateSlug checks if slug is valid and unique
func (s *Service) ValidateSlug(slug string, excludeID int) error {
	if slug == "" {
		return fmt.Errorf("slug cannot be empty")
	}

	if slug != GenerateSlug(slug) {
		return fmt.Errorf("slug contains invalid characters")
	}

	// Check uniqueness
	existing, err := s.repo.FindBySlug(slug)
	if err == nil && existing.ID != excludeID {
		return fmt.Errorf("slug already exists")
	}

	return nil
}
