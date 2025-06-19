package dbHelper

import (
	"github.com/jmoiron/sqlx"
	"todo-app/models"
)

// CreateUser inserts a new user into the DB
func CreateUser(db *sqlx.DB, user models.User) error {
	query := `INSERT INTO users (name, email, password) VALUES ($1, $2, $3)`
	_, err := db.Exec(query, user.Name, user.Email, user.Password)
	return err
}

// GetUserByEmail fetches a user by email
func GetUserByEmail(db *sqlx.DB, email string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, name, email, password, created_at, archived_at FROM users WHERE email = $1`
	err := db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.ArchivedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}
