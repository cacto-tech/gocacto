package seeds

import (
	"database/sql"
	"fmt"
	"log"
)

// Seeder handles database seeding
type Seeder struct {
	db *sql.DB
}

// NewSeeder creates a new seeder
func NewSeeder(db *sql.DB) *Seeder {
	return &Seeder{db: db}
}

// SeedAll seeds all data
func (s *Seeder) SeedAll() error {
	log.Println("🌱 Starting database seeding...")

	if err := s.SeedComponents(); err != nil {
		return fmt.Errorf("failed to seed components: %w", err)
	}

	if err := s.SeedPages(); err != nil {
		return fmt.Errorf("failed to seed pages: %w", err)
	}

	if err := s.SeedUsers(); err != nil {
		return fmt.Errorf("failed to seed users: %w", err)
	}

	log.Println("✅ Database seeding completed")
	return nil
}

// SeedComponents seeds component data
func (s *Seeder) SeedComponents() error {
	log.Println("  📦 Seeding components...")

	components := Components()
	
	for _, comp := range components {
		// Check if component already exists
		var exists bool
		err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM components WHERE name = ?)", comp.Name).Scan(&exists)
		if err != nil {
			return err
		}

		if exists {
			log.Printf("    ⏭️  Component '%s' already exists, skipping", comp.Name)
			continue
		}

		query := `
			INSERT INTO components (type, name, title, subtitle, content,
								image_url, link_url, link_text, data_json)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		`

		_, err = s.db.Exec(query,
			comp.Type, comp.Name, comp.Title, comp.Subtitle, comp.Content,
			comp.ImageURL, comp.LinkURL, comp.LinkText, comp.DataJSON,
		)

		if err != nil {
			return fmt.Errorf("failed to insert component %s: %w", comp.Name, err)
		}

		log.Printf("    ✓ Seeded component: %s", comp.Name)
	}

	return nil
}

// SeedPages seeds page data
func (s *Seeder) SeedPages() error {
	log.Println("  📄 Seeding pages...")

	pages := Pages()

	for _, pageSeed := range pages {
		// Check if page already exists
		var exists bool
		var query string
		if pageSeed.Page.Slug == "" {
			// Home page check (slug is empty)
			query = "SELECT EXISTS(SELECT 1 FROM pages WHERE slug = '' OR slug IS NULL)"
		} else {
			query = "SELECT EXISTS(SELECT 1 FROM pages WHERE slug = ?)"
		}
		
		var err error
		if pageSeed.Page.Slug == "" {
			err = s.db.QueryRow(query).Scan(&exists)
		} else {
			err = s.db.QueryRow(query, pageSeed.Page.Slug).Scan(&exists)
		}
		
		if err != nil {
			return err
		}

		pageName := pageSeed.Page.Slug
		if pageName == "" {
			pageName = "home"
		}
		
		if exists {
			log.Printf("    ⏭️  Page '%s' already exists, skipping", pageName)
			continue
		}

		// Insert page
		pageQuery := `
			INSERT INTO pages (slug, title, content, meta_title, meta_description,
							meta_keywords, og_image, status, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, datetime('now'), datetime('now'))
		`

		result, err := s.db.Exec(pageQuery,
			pageSeed.Page.Slug, pageSeed.Page.Title, pageSeed.Page.Content,
			pageSeed.Page.MetaTitle, pageSeed.Page.MetaDescription,
			pageSeed.Page.MetaKeywords, pageSeed.Page.OGImage, pageSeed.Page.Status,
		)

		if err != nil {
			pageName := pageSeed.Page.Slug
			if pageName == "" {
				pageName = "home"
			}
			return fmt.Errorf("failed to insert page %s: %w", pageName, err)
		}

		pageID, err := result.LastInsertId()
		if err != nil {
			return err
		}

		// Associate components with page
		for position, componentName := range pageSeed.ComponentNames {
			// Get component ID
			var componentID int
			err = s.db.QueryRow("SELECT id FROM components WHERE name = ?", componentName).Scan(&componentID)
			if err != nil {
				log.Printf("    ⚠️  Component '%s' not found, skipping association", componentName)
				continue
			}

			// Insert page_component association
			associationQuery := `
				INSERT INTO page_components (page_id, component_id, position)
				VALUES (?, ?, ?)
			`

			_, err = s.db.Exec(associationQuery, pageID, componentID, position)
			if err != nil {
				return fmt.Errorf("failed to associate component %s with page %s: %w", componentName, pageSeed.Page.Slug, err)
			}
		}

		log.Printf("    ✓ Seeded page: %s", pageName)
	}

	return nil
}

// SeedUsers seeds user data
func (s *Seeder) SeedUsers() error {
	log.Println("  👤 Seeding users...")

	users := SeedUsers()

	for _, userSeed := range users {
		// Check if user already exists
		var exists bool
		err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", userSeed.User.Email).Scan(&exists)
		if err != nil {
			return err
		}

		if exists {
			log.Printf("    ⏭️  User '%s' already exists, skipping", userSeed.User.Email)
			continue
		}

		// Insert user
		query := `
			INSERT INTO users (email, password_hash, name, role, is_active, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, datetime('now'), datetime('now'))
		`

		_, err = s.db.Exec(query,
			userSeed.User.Email, userSeed.User.PasswordHash, userSeed.User.Name,
			userSeed.User.Role, userSeed.User.IsActive,
		)

		if err != nil {
			return fmt.Errorf("failed to insert user %s: %w", userSeed.User.Email, err)
		}

		log.Printf("    ✓ Seeded user: %s (%s)", userSeed.User.Email, userSeed.User.Role)
	}

	return nil
}
