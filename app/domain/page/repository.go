package page

// Repository defines the interface for page data persistence
// This interface belongs to the domain layer and should not depend on infrastructure
type Repository interface {
	FindByID(id int) (*Page, error)
	FindBySlug(slug string) (*Page, error)
	FindAll() ([]*Page, error)
	FindPublished() ([]*Page, error)
	Create(page *Page) error
	Update(page *Page) error
	Delete(id int) error
	GetComponents(pageID int) ([]Component, error)
}
