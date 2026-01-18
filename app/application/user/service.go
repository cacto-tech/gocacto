package user

import (
	"cacto-cms/app/domain/user"
	"cacto-cms/app/shared/errors"
	"time"
)

// Service handles business logic for users
type Service struct {
	repo user.Repository
}

// NewService creates a new user service
func NewService(repo user.Repository) *Service {
	return &Service{repo: repo}
}

// GetUserByID retrieves a user by ID
func (s *Service) GetUserByID(id int) (*user.User, error) {
	u, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrCodeNotFound, "User not found", 404)
	}
	return u, nil
}

// GetUserByEmail retrieves a user by email
func (s *Service) GetUserByEmail(email string) (*user.User, error) {
	u, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrCodeNotFound, "User not found", 404)
	}
	return u, nil
}

// GetAllUsers retrieves all users
func (s *Service) GetAllUsers() ([]*user.User, error) {
	return s.repo.FindAll()
}

// CreateUser creates a new user
func (s *Service) CreateUser(u *user.User) error {
	// Check if email already exists
	existing, err := s.repo.FindByEmail(u.Email)
	if err == nil && existing != nil {
		return errors.NewConflict("User with this email already exists")
	}

	// Set default role if not provided
	if u.Role == "" {
		u.Role = user.RoleViewer
	}

	// Set timestamps
	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now
	u.IsActive = true

	return s.repo.Create(u)
}

// UpdateUser updates an existing user
func (s *Service) UpdateUser(u *user.User) error {
	u.UpdatedAt = time.Now()
	return s.repo.Update(u)
}

// DeleteUser deletes a user by ID
func (s *Service) DeleteUser(id int) error {
	return s.repo.Delete(id)
}

// UpdateLastLogin updates user's last login time
func (s *Service) UpdateLastLogin(id int) error {
	return s.repo.UpdateLastLogin(id)
}
