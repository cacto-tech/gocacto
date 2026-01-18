package user

import (
	"database/sql"
	"fmt"
	"time"

	"cacto-cms/app/domain/user"
)

// Repository implements user.Repository interface
type Repository struct {
	db *sql.DB
}

// NewRepository creates a new user repository
func NewRepository(db *sql.DB) user.Repository {
	return &Repository{db: db}
}

// FindByID retrieves a user by ID
func (r *Repository) FindByID(id int) (*user.User, error) {
	query := `
		SELECT id, email, password_hash, name, role, is_active, 
		       last_login_at, created_at, updated_at
		FROM users WHERE id = ?
	`

	u := &user.User{}
	var lastLoginAt sql.NullTime

	err := r.db.QueryRow(query, id).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.Role,
		&u.IsActive, &lastLoginAt, &u.CreatedAt, &u.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, err
	}

	if lastLoginAt.Valid {
		u.LastLoginAt = &lastLoginAt.Time
	}

	return u, nil
}

// FindByEmail retrieves a user by email
func (r *Repository) FindByEmail(email string) (*user.User, error) {
	query := `
		SELECT id, email, password_hash, name, role, is_active,
		       last_login_at, created_at, updated_at
		FROM users WHERE email = ?
	`

	u := &user.User{}
	var lastLoginAt sql.NullTime

	err := r.db.QueryRow(query, email).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.Role,
		&u.IsActive, &lastLoginAt, &u.CreatedAt, &u.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, err
	}

	if lastLoginAt.Valid {
		u.LastLoginAt = &lastLoginAt.Time
	}

	return u, nil
}

// FindAll retrieves all users
func (r *Repository) FindAll() ([]*user.User, error) {
	query := `
		SELECT id, email, password_hash, name, role, is_active,
		       last_login_at, created_at, updated_at
		FROM users ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*user.User, 0)
	for rows.Next() {
		u := &user.User{}
		var lastLoginAt sql.NullTime

		err := rows.Scan(
			&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.Role,
			&u.IsActive, &lastLoginAt, &u.CreatedAt, &u.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if lastLoginAt.Valid {
			u.LastLoginAt = &lastLoginAt.Time
		}

		users = append(users, u)
	}

	return users, rows.Err()
}

// Create creates a new user
func (r *Repository) Create(u *user.User) error {
	query := `
		INSERT INTO users (email, password_hash, name, role, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(query,
		u.Email, u.PasswordHash, u.Name, u.Role, u.IsActive,
		u.CreatedAt, u.UpdatedAt,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	u.ID = int(id)
	return nil
}

// Update updates an existing user
func (r *Repository) Update(u *user.User) error {
	query := `
		UPDATE users 
		SET email = ?, password_hash = ?, name = ?, role = ?, 
		    is_active = ?, updated_at = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(query,
		u.Email, u.PasswordHash, u.Name, u.Role, u.IsActive,
		u.UpdatedAt, u.ID,
	)
	return err
}

// Delete deletes a user by ID
func (r *Repository) Delete(id int) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

// UpdateLastLogin updates user's last login time
func (r *Repository) UpdateLastLogin(id int) error {
	query := `UPDATE users SET last_login_at = ? WHERE id = ?`
	_, err := r.db.Exec(query, time.Now(), id)
	return err
}
