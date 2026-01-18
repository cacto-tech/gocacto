package media

import (
	"database/sql"
	"fmt"

	"cacto-cms/app/domain/media"
)

// Repository implements media.Repository interface
type Repository struct {
	db *sql.DB
}

// NewRepository creates a new media repository
func NewRepository(db *sql.DB) media.Repository {
	return &Repository{db: db}
}

// FindByID retrieves a media by ID
func (r *Repository) FindByID(id int) (*media.Media, error) {
	query := `
		SELECT id, filename, original_name, mime_type, size, alt_text, created_at
		FROM media WHERE id = ?
	`

	m := &media.Media{}
	err := r.db.QueryRow(query, id).Scan(
		&m.ID, &m.Filename, &m.OriginalName, &m.MimeType,
		&m.Size, &m.AltText, &m.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("media not found")
	}
	if err != nil {
		return nil, err
	}

	// Set path and URL (these are computed, not stored)
	m.Path = fmt.Sprintf("/uploads/%s", m.Filename)
	m.URL = fmt.Sprintf("/uploads/%s", m.Filename)

	return m, nil
}

// FindAll retrieves all media with pagination
func (r *Repository) FindAll(limit, offset int) ([]*media.Media, error) {
	query := `
		SELECT id, filename, original_name, mime_type, size, alt_text, created_at
		FROM media ORDER BY created_at DESC LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	mediaList := make([]*media.Media, 0)
	for rows.Next() {
		m := &media.Media{}
		err := rows.Scan(
			&m.ID, &m.Filename, &m.OriginalName, &m.MimeType,
			&m.Size, &m.AltText, &m.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		m.Path = fmt.Sprintf("/uploads/%s", m.Filename)
		m.URL = fmt.Sprintf("/uploads/%s", m.Filename)

		mediaList = append(mediaList, m)
	}

	return mediaList, rows.Err()
}

// FindByFilename retrieves a media by filename
func (r *Repository) FindByFilename(filename string) (*media.Media, error) {
	query := `
		SELECT id, filename, original_name, mime_type, size, alt_text, created_at
		FROM media WHERE filename = ?
	`

	m := &media.Media{}
	err := r.db.QueryRow(query, filename).Scan(
		&m.ID, &m.Filename, &m.OriginalName, &m.MimeType,
		&m.Size, &m.AltText, &m.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("media not found")
	}
	if err != nil {
		return nil, err
	}

	m.Path = fmt.Sprintf("/uploads/%s", m.Filename)
	m.URL = fmt.Sprintf("/uploads/%s", m.Filename)

	return m, nil
}

// Create creates a new media record
func (r *Repository) Create(m *media.Media) error {
	query := `
		INSERT INTO media (filename, original_name, mime_type, size, alt_text, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(query,
		m.Filename, m.OriginalName, m.MimeType, m.Size, m.AltText, m.CreatedAt,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	m.ID = int(id)
	return nil
}

// Update updates media metadata
func (r *Repository) Update(m *media.Media) error {
	query := `
		UPDATE media 
		SET filename = ?, original_name = ?, mime_type = ?, size = ?, alt_text = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(query,
		m.Filename, m.OriginalName, m.MimeType, m.Size, m.AltText, m.ID,
	)
	return err
}

// Delete deletes a media by ID
func (r *Repository) Delete(id int) error {
	query := `DELETE FROM media WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

// Count returns total number of media files
func (r *Repository) Count() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM media").Scan(&count)
	return count, err
}
