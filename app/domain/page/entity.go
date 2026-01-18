package page

import "time"

// Status represents the publication status of a page
type Status string

const (
	StatusDraft     Status = "draft"
	StatusPublished Status = "published"
	StatusArchived  Status = "archived"
)

// Page represents a page entity in the domain
type Page struct {
	ID              int        `json:"id"`
	Slug            string    `json:"slug"`
	Title           string    `json:"title"`
	Content         string    `json:"content"`
	MetaTitle       string    `json:"meta_title"`
	MetaDescription string    `json:"meta_description"`
	MetaKeywords    string    `json:"meta_keywords"`
	OGImage         string    `json:"og_image"`
	Status          Status    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Components      []Component `json:"components,omitempty"`
}

// Component represents a page component entity
type Component struct {
	ID       int    `json:"id"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Content  string `json:"content"`
	ImageURL string `json:"image_url"`
	LinkURL  string `json:"link_url"`
	LinkText string `json:"link_text"`
	DataJSON string `json:"data_json"`
	Position int    `json:"position"`
}
