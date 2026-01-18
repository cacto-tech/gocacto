package media

import "time"

// Media represents a media file entity
type Media struct {
	ID          int       `json:"id"`
	Filename    string    `json:"filename"`
	OriginalName string   `json:"original_name"`
	MimeType   string    `json:"mime_type"`
	Size        int64     `json:"size"`
	AltText     string    `json:"alt_text,omitempty"`
	Path        string    `json:"path"`
	URL         string    `json:"url"`
	CreatedAt   time.Time `json:"created_at"`
}

// IsImage checks if media is an image
func (m *Media) IsImage() bool {
	return m.MimeType != "" && 
		(m.MimeType == "image/jpeg" || 
		 m.MimeType == "image/png" || 
		 m.MimeType == "image/gif" || 
		 m.MimeType == "image/webp" ||
		 m.MimeType == "image/svg+xml")
}

// IsVideo checks if media is a video
func (m *Media) IsVideo() bool {
	return m.MimeType != "" && 
		(m.MimeType == "video/mp4" || 
		 m.MimeType == "video/webm" || 
		 m.MimeType == "video/ogg")
}
