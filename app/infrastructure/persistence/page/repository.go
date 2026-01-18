package page

import (
	"database/sql"
	"fmt"
	"time"

	"cacto-cms/app/domain/page"
)

// Repository implements the page.Repository interface using SQLite
type Repository struct {
	db *sql.DB
}

// NewRepository creates a new page repository
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// FindByID retrieves a page by its ID
func (r *Repository) FindByID(id int) (*page.Page, error) {
	query := `
		SELECT id, slug, title, content, meta_title, meta_description, 
		       meta_keywords, og_image, status, created_at, updated_at
		FROM pages WHERE id = ?
	`

	p := &page.Page{}
	err := r.db.QueryRow(query, id).Scan(
		&p.ID, &p.Slug, &p.Title, &p.Content, &p.MetaTitle, &p.MetaDescription,
		&p.MetaKeywords, &p.OGImage, &p.Status, &p.CreatedAt, &p.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("page not found")
	}
	if err != nil {
		return nil, err
	}

	return p, nil
}

// FindBySlug retrieves a page by its slug
// Empty slug means home page
func (r *Repository) FindBySlug(slug string) (*page.Page, error) {
	var query string
	if slug == "" {
		// Home page has empty slug
		query = `
			SELECT id, slug, title, content, meta_title, meta_description, 
			       meta_keywords, og_image, status, created_at, updated_at
			FROM pages WHERE slug = '' OR slug IS NULL
		`
	} else {
		query = `
			SELECT id, slug, title, content, meta_title, meta_description, 
			       meta_keywords, og_image, status, created_at, updated_at
			FROM pages WHERE slug = ?
		`
	}

	p := &page.Page{}
	var err error
	if slug == "" {
		err = r.db.QueryRow(query).Scan(
			&p.ID, &p.Slug, &p.Title, &p.Content, &p.MetaTitle, &p.MetaDescription,
			&p.MetaKeywords, &p.OGImage, &p.Status, &p.CreatedAt, &p.UpdatedAt,
		)
	} else {
		err = r.db.QueryRow(query, slug).Scan(
			&p.ID, &p.Slug, &p.Title, &p.Content, &p.MetaTitle, &p.MetaDescription,
			&p.MetaKeywords, &p.OGImage, &p.Status, &p.CreatedAt, &p.UpdatedAt,
		)
	}

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("page not found")
	}
	if err != nil {
		return nil, err
	}

	return p, nil
}

// FindAll retrieves all pages
func (r *Repository) FindAll() ([]*page.Page, error) {
	query := `
		SELECT id, slug, title, content, meta_title, meta_description, 
		       meta_keywords, og_image, status, created_at, updated_at
		FROM pages ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanPages(rows)
}

// FindPublished retrieves only published pages
func (r *Repository) FindPublished() ([]*page.Page, error) {
	query := `
		SELECT id, slug, title, content, meta_title, meta_description, 
		       meta_keywords, og_image, status, created_at, updated_at
		FROM pages WHERE status = 'published' ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanPages(rows)
}

// Create creates a new page
func (r *Repository) Create(p *page.Page) error {
	query := `
		INSERT INTO pages (slug, title, content, meta_title, meta_description, 
		                   meta_keywords, og_image, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(query,
		p.Slug, p.Title, p.Content, p.MetaTitle, p.MetaDescription,
		p.MetaKeywords, p.OGImage, p.Status, p.CreatedAt, p.UpdatedAt,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	p.ID = int(id)
	return nil
}

// Update updates an existing page
func (r *Repository) Update(p *page.Page) error {
	query := `
		UPDATE pages 
		SET slug = ?, title = ?, content = ?, meta_title = ?, meta_description = ?,
		    meta_keywords = ?, og_image = ?, status = ?, updated_at = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(query,
		p.Slug, p.Title, p.Content, p.MetaTitle, p.MetaDescription,
		p.MetaKeywords, p.OGImage, p.Status, time.Now(), p.ID,
	)

	return err
}

// Delete deletes a page by ID
func (r *Repository) Delete(id int) error {
	query := "DELETE FROM pages WHERE id = ?"
	_, err := r.db.Exec(query, id)
	return err
}

// GetComponents retrieves all components for a page
func (r *Repository) GetComponents(pageID int) ([]page.Component, error) {
	query := `
		SELECT c.id, c.type, c.name, c.title, c.subtitle, c.content,
		       c.image_url, c.link_url, c.link_text, c.data_json, pc.position
		FROM components c
		JOIN page_components pc ON c.id = pc.component_id
		WHERE pc.page_id = ?
		ORDER BY pc.position ASC
	`

	rows, err := r.db.Query(query, pageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var components []page.Component
	for rows.Next() {
		var c page.Component
		err := rows.Scan(
			&c.ID, &c.Type, &c.Name, &c.Title, &c.Subtitle, &c.Content,
			&c.ImageURL, &c.LinkURL, &c.LinkText, &c.DataJSON, &c.Position,
		)
		if err != nil {
			return nil, err
		}
		components = append(components, c)
	}

	return components, nil
}

// scanPages is a helper method to scan multiple pages from rows
func (r *Repository) scanPages(rows *sql.Rows) ([]*page.Page, error) {
	var pages []*page.Page

	for rows.Next() {
		p := &page.Page{}
		err := rows.Scan(
			&p.ID, &p.Slug, &p.Title, &p.Content, &p.MetaTitle, &p.MetaDescription,
			&p.MetaKeywords, &p.OGImage, &p.Status, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		pages = append(pages, p)
	}

	return pages, nil
}
