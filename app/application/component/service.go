package component

import (
	"cacto-cms/app/domain/component"
)

// Service handles business logic for components
type Service struct {
	repo component.Repository
}

// NewService creates a new component service
func NewService(repo component.Repository) *Service {
	return &Service{repo: repo}
}

// GetComponentByID retrieves a component by ID
func (s *Service) GetComponentByID(id int) (*component.Component, error) {
	c, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return c.MergeWithDefaults(), nil
}

// GetComponentByName retrieves a component by name
func (s *Service) GetComponentByName(name string) (*component.Component, error) {
	c, err := s.repo.FindByName(name)
	if err != nil {
		return nil, err
	}
	return c.MergeWithDefaults(), nil
}

// GetComponentsByType retrieves components by type
func (s *Service) GetComponentsByType(componentType component.Type) ([]*component.Component, error) {
	components, err := s.repo.FindByType(componentType)
	if err != nil {
		return nil, err
	}
	
	// Apply defaults to all components
	for _, c := range components {
		c.MergeWithDefaults()
	}
	
	return components, nil
}

// GetAllComponents retrieves all components
func (s *Service) GetAllComponents() ([]*component.Component, error) {
	return s.repo.FindAll()
}

// CreateComponent creates a new component
func (s *Service) CreateComponent(c *component.Component) error {
	return s.repo.Create(c)
}

// UpdateComponent updates an existing component
func (s *Service) UpdateComponent(c *component.Component) error {
	return s.repo.Update(c)
}

// DeleteComponent deletes a component by ID
func (s *Service) DeleteComponent(id int) error {
	return s.repo.Delete(id)
}
