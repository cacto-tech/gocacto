package user

import "time"

// Role represents user roles
type Role string

const (
	RoleAdmin    Role = "admin"
	RoleEditor   Role = "editor"
	RoleAuthor   Role = "author"
	RoleViewer   Role = "viewer"
)

// User represents a user entity
type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Never expose password hash
	Name         string    `json:"name"`
	Role         Role      `json:"role"`
	IsActive     bool      `json:"is_active"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// HasPermission checks if user has a specific permission
func (u *User) HasPermission(permission string) bool {
	if !u.IsActive {
		return false
	}

	// Admin has all permissions
	if u.Role == RoleAdmin {
		return true
	}

	// Define role-based permissions
	permissions := map[Role][]string{
		RoleEditor: {"pages:read", "pages:write", "pages:delete", "components:read", "components:write"},
		RoleAuthor: {"pages:read", "pages:write", "components:read"},
		RoleViewer: {"pages:read", "components:read"},
	}

	rolePerms, exists := permissions[u.Role]
	if !exists {
		return false
	}

	for _, perm := range rolePerms {
		if perm == permission {
			return true
		}
	}

	return false
}

// CanEdit checks if user can edit content
func (u *User) CanEdit() bool {
	return u.HasPermission("pages:write")
}

// CanDelete checks if user can delete content
func (u *User) CanDelete() bool {
	return u.HasPermission("pages:delete")
}
