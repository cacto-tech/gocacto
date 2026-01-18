package media

// Repository defines the interface for media data persistence
type Repository interface {
	FindByID(id int) (*Media, error)
	FindAll(limit, offset int) ([]*Media, error)
	FindByFilename(filename string) (*Media, error)
	Create(media *Media) error
	Update(media *Media) error
	Delete(id int) error
	Count() (int, error)
}
