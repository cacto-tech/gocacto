package seeds

import (
	"cacto-cms/app/domain/user"
	"cacto-cms/app/shared/auth"
)

// SeedUsers returns user seed data
func SeedUsers() []UserSeed {
	hasher := auth.NewPasswordHasher()
	
	// Default admin password: "admin123" (change in production!)
	adminPasswordHash, _ := hasher.HashPassword("admin123")
	
	return []UserSeed{
		{
			User: user.User{
				Email:        "admin@cacto-cms.local",
				PasswordHash: adminPasswordHash,
				Name:         "Admin User",
				Role:         user.RoleAdmin,
				IsActive:     true,
			},
		},
		{
			User: user.User{
				Email:        "editor@cacto-cms.local",
				PasswordHash: adminPasswordHash, // Same password for demo
				Name:         "Editor User",
				Role:         user.RoleEditor,
				IsActive:     true,
			},
		},
	}
}

// UserSeed represents a user seed entry
type UserSeed struct {
	User user.User
}
