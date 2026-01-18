package component

// Repository defines the interface for component data persistence
type Repository interface {
	FindByID(id int) (*Component, error)
	FindByName(name string) (*Component, error)
	FindByType(componentType Type) ([]*Component, error)
	FindAll() ([]*Component, error)
	Create(component *Component) error
	Update(component *Component) error
	Delete(id int) error
}
