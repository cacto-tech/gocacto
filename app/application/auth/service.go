package auth

import (
	userservice "cacto-cms/app/application/user"
	"cacto-cms/app/domain/user"
	"cacto-cms/app/shared/auth"
	"cacto-cms/app/shared/errors"
	"time"
)

// Service handles authentication business logic
type Service struct {
	userService *userservice.Service
	jwtManager  *auth.JWTManager
	hasher      *auth.PasswordHasher
}

// NewService creates a new auth service
func NewService(userService *userservice.Service, jwtSecret string, tokenDuration time.Duration) *Service {
	return &Service{
		userService: userService,
		jwtManager:  auth.NewJWTManager(jwtSecret, tokenDuration),
		hasher:      auth.NewPasswordHasher(),
	}
}

// LoginRequest represents login request data
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// LoginResponse represents login response data
type LoginResponse struct {
	Token string      `json:"token"`
	User  *user.User  `json:"user"`
}

// Login authenticates a user and returns a JWT token
func (s *Service) Login(req *LoginRequest) (*LoginResponse, error) {
	// Get user by email
	u, err := s.userService.GetUserByEmail(req.Email)
	if err != nil {
		return nil, errors.NewUnauthorized("Invalid credentials")
	}

	// Check if user is active
	if !u.IsActive {
		return nil, errors.NewForbidden("User account is inactive")
	}

	// Verify password
	valid, err := s.hasher.VerifyPassword(req.Password, u.PasswordHash)
	if err != nil {
		return nil, errors.NewInternal("Failed to verify password", err)
	}

	if !valid {
		return nil, errors.NewUnauthorized("Invalid credentials")
	}

	// Generate JWT token
	token, err := s.jwtManager.GenerateToken(u.ID, u.Email, string(u.Role))
	if err != nil {
		return nil, errors.NewInternal("Failed to generate token", err)
	}

	// Update last login
	_ = s.userService.UpdateLastLogin(u.ID)

	return &LoginResponse{
		Token: token,
		User:  u,
	}, nil
}

// RegisterRequest represents registration request data
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Name     string `json:"name" validate:"required,min=2"`
	Role     string `json:"role,omitempty"`
}

// Register creates a new user account
func (s *Service) Register(req *RegisterRequest) (*user.User, error) {
	// Check if user already exists
	existing, err := s.userService.GetUserByEmail(req.Email)
	if err == nil && existing != nil {
		return nil, errors.NewConflict("User with this email already exists")
	}

	// Hash password
	passwordHash, err := s.hasher.HashPassword(req.Password)
	if err != nil {
		return nil, errors.NewInternal("Failed to hash password", err)
	}

	// Determine role
	role := user.RoleViewer
	if req.Role != "" {
		switch req.Role {
		case "admin":
			role = user.RoleAdmin
		case "editor":
			role = user.RoleEditor
		case "author":
			role = user.RoleAuthor
		}
	}

	// Create user
	newUser := &user.User{
		Email:        req.Email,
		PasswordHash: passwordHash,
		Name:         req.Name,
		Role:         role,
		IsActive:     true,
	}

	if err := s.userService.CreateUser(newUser); err != nil {
		return nil, errors.NewInternal("Failed to create user", err)
	}

	return newUser, nil
}

// ValidateToken validates a JWT token and returns claims
func (s *Service) ValidateToken(tokenString string) (*auth.Claims, error) {
	claims, err := s.jwtManager.ValidateToken(tokenString)
	if err != nil {
		return nil, errors.NewUnauthorized("Invalid or expired token")
	}
	return claims, nil
}
