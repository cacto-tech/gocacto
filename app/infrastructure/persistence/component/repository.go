package component

import (
	"database/sql"
	"fmt"

	"cacto-cms/app/domain/component"
)

// Repository implements the component.Repository interface using SQLite
type Repository struct {
	db *sql.DB
}

// NewRepository creates a new component repository
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// FindByID retrieves a component by its ID
func (r *Repository) FindByID(id int) (*component.Component, error) {
	query := `
		SELECT id, type, name, title, subtitle, content,
		       image_url, link_url, link_text, data_json
		FROM components WHERE id = ?
	`

	c := &component.Component{}
	err := r.db.QueryRow(query, id).Scan(
		&c.ID, &c.Type, &c.Name, &c.Title, &c.Subtitle, &c.Content,
		&c.ImageURL, &c.LinkURL, &c.LinkText, &c.DataJSON,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("component not found")
	}
	if err != nil {
		return nil, err
	}

	return c, nil
}

// FindByName retrieves a component by its name
func (r *Repository) FindByName(name string) (*component.Component, error) {
	query := `
		SELECT id, type, name, title, subtitle, content,
		       image_url, link_url, link_text, data_json
		FROM components WHERE name = ?
	`

	c := &component.Component{}
	err := r.db.QueryRow(query, name).Scan(
		&c.ID, &c.Type, &c.Name, &c.Title, &c.Subtitle, &c.Content,
		&c.ImageURL, &c.LinkURL, &c.LinkText, &c.DataJSON,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("component not found")
	}
	if err != nil {
		return nil, err
	}

	return c, nil
}

// FindByType retrieves components by type
func (r *Repository) FindByType(componentType component.Type) ([]*component.Component, error) {
	query := `
		SELECT id, type, name, title, subtitle, content,
		       image_url, link_url, link_text, data_json
		FROM components WHERE type = ?
		ORDER BY id ASC
	`

	rows, err := r.db.Query(query, componentType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanComponents(rows)
}

// FindAll retrieves all components
func (r *Repository) FindAll() ([]*component.Component, error) {
	query := `
		SELECT id, type, name, title, subtitle, content,
		       image_url, link_url, link_text, data_json
		FROM components
		ORDER BY id ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanComponents(rows)
}

// Create creates a new component
func (r *Repository) Create(c *component.Component) error {
	query := `
		INSERT INTO components (type, name, title, subtitle, content,
		                       image_url, link_url, link_text, data_json)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(query,
		c.Type, c.Name, c.Title, c.Subtitle, c.Content,
		c.ImageURL, c.LinkURL, c.LinkText, c.DataJSON,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	c.ID = int(id)
	return nil
}

// Update updates an existing component
func (r *Repository) Update(c *component.Component) error {
	query := `
		UPDATE components 
		SET type = ?, name = ?, title = ?, subtitle = ?, content = ?,
		    image_url = ?, link_url = ?, link_text = ?, data_json = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(query,
		c.Type, c.Name, c.Title, c.Subtitle, c.Content,
		c.ImageURL, c.LinkURL, c.LinkText, c.DataJSON, c.ID,
	)

	return err
}

// Delete deletes a component by ID
func (r *Repository) Delete(id int) error {
	query := "DELETE FROM components WHERE id = ?"
	_, err := r.db.Exec(query, id)
	return err
}

// scanComponents is a helper method to scan multiple components from rows
func (r *Repository) scanComponents(rows *sql.Rows) ([]*component.Component, error) {
	var components []*component.Component

	for rows.Next() {
		c := &component.Component{}
		err := rows.Scan(
			&c.ID, &c.Type, &c.Name, &c.Title, &c.Subtitle, &c.Content,
			&c.ImageURL, &c.LinkURL, &c.LinkText, &c.DataJSON,
		)
		if err != nil {
			return nil, err
		}
		components = append(components, c)
	}

	return components, nil
}
