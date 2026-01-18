package user

// Repository defines the interface for user data persistence
type Repository interface {
	FindByID(id int) (*User, error)
	FindByEmail(email string) (*User, error)
	FindAll() ([]*User, error)
	Create(user *User) error
	Update(user *User) error
	Delete(id int) error
	UpdateLastLogin(id int) error
}
