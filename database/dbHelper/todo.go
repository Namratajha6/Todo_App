package dbHelper

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"todo-app/models"
)

// CreateTodo inserts a todo into the DB
func CreateTodo(db *sqlx.DB, todo *models.CreateTodo) error {
	query := `INSERT INTO todos (user_id, name, description) VALUES ($1, $2, $3)`
	_, err := db.Exec(query, todo.UserID, todo.Name, todo.Description)
	return err
}

func UpdateTodo(db *sqlx.DB, todo *models.UpdateTodo, userID string) error {
	query := `
		UPDATE todos
		SET name = $1,
		    description = $2,
		    is_completed = $3
		WHERE id = $4 AND user_id = $5
	`
	res, err := db.Exec(query, todo.Name, todo.Description, todo.IsCompleted, todo.ID, userID)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("not updated, check details")
	}
	return nil
}

func ValidateSession(db *sqlx.DB, sessionID string) (string, error) {
	query := `
		SELECT user_id FROM user_session
		WHERE id = $1 AND archived_at IS NULL AND expiry_at > NOW()
	`
	var userID string
	err := db.QueryRow(query, sessionID).Scan(&userID)
	if err != nil {
		return "", fmt.Errorf("invalid or expired session")
	}
	return userID, nil
}

func GetAllTodos(db *sqlx.DB, userID string) ([]models.Todo, error) {
	var todos []models.Todo
	query := `SELECT name, description, is_completed FROM todos WHERE user_id = $1 AND archived_at IS NULL`
	err := db.Select(&todos, query, userID)
	return todos, err
}

func GetTodoByID(db *sqlx.DB, todoID, userID string) (*models.Todo, error) {
	var todo models.Todo
	query := `SELECT name, description, is_completed FROM todos WHERE id = $1 AND user_id = $2 AND archived_at IS NULL`
	err := db.Get(&todo, query, todoID, userID)
	if err != nil {
		return nil, fmt.Errorf("todo not found")
	}
	return &todo, nil
}

func ArchiveTodo(db *sqlx.DB, todoID, userID string) error {
	query := `UPDATE todos SET archived_at = NOW() WHERE id = $1 AND user_id = $2 AND archived_at IS NULL`
	res, err := db.Exec(query, todoID, userID)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("todo not found or already archived")
	}
	return nil
}
