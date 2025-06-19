package models

import (
	"time"
)

type Todo struct {
	ID          string     `json:"id" db:"id"`
	UserID      string     `json:"user_id" db:"user_id"`
	Name        string     `json:"name" db:"name"`
	Description string     `json:"description" db:"description"`
	IsCompleted bool       `json:"is_completed" db:"is_completed"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	ArchivedAt  *time.Time `json:"archived_at" db:"archived_at"`
}

type CreateTodo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	UserID      string `json:"-"`
}

type UpdateTodo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsCompleted bool   `json:"is_completed"`
}
