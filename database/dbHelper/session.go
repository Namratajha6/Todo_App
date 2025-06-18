package dbHelper

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"todo-app/models"
)

func CreateUserSession(db *sqlx.DB, session *models.UserSession) error {
	query := `INSERT INTO user_session (user_id) VALUES ($1) RETURNING id, created_at, archived_at`
	return db.QueryRow(query, session.UserID).Scan(&session.ID, &session.CreatedAt, &session.ArchivedAt)
}

func ArchiveSession(db *sqlx.DB, sessionID string) error {
	query := `UPDATE user_session SET archived_at = NOW() WHERE id = $1`
	_, err := db.Exec(query, sessionID)
	return err
}

func LogoutIfNotExpired(db *sqlx.DB, sessionID string) error {
	query := `
		UPDATE user_session
		SET archived_at = NOW()
		WHERE id = $1
		  AND archived_at IS NULL
		  AND expiry_at > NOW()
		RETURNING id;
	`

	var id string
	err := db.QueryRow(query, sessionID).Scan(&id)
	if err != nil {
		return fmt.Errorf("cannot logout: session is either expired or already archived")
	}
	return nil
}
