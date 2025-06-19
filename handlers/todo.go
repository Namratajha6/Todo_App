package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"todo-app/database"
	"todo-app/database/dbHelper"
	"todo-app/models"
)

func CreateTodoHandler(w http.ResponseWriter, r *http.Request) {
	//  Get session ID from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
		return
	}
	sessionID := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))

	// Validate session
	userID, err := dbHelper.ValidateSession(database.Todo, sessionID)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	// Decode the request body into a CreateTodo struct
	var req models.CreateTodo
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Attach the userID from session to the request
	req.UserID = userID

	//  Insert into DB
	if err := dbHelper.CreateTodo(database.Todo, &req); err != nil {
		http.Error(w, "Failed to create todo", http.StatusInternalServerError)
		return
	}

	// Respond with the created todo
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(req)
}

func UpdateTodoHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
		return
	}
	sessionID := strings.TrimPrefix(authHeader, "Bearer ")
	sessionID = strings.TrimSpace(sessionID)

	userID, err := dbHelper.ValidateSession(database.Todo, sessionID)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	var req models.UpdateTodo
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Println("Update Request Payload:", req)

	err = dbHelper.UpdateTodo(database.Todo, &req, userID)
	if err != nil {
		http.Error(w, "Failed to update todo", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Todo updated successfully"))
}

func GetAllTodosHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
		return
	}
	sessionID := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))

	userID, err := dbHelper.ValidateSession(database.Todo, sessionID)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	todos, err := dbHelper.GetAllTodos(database.Todo, userID)
	if err != nil {
		http.Error(w, "Failed to fetch todos: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(todos)
}

func GetTodoHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
		return
	}
	sessionID := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))

	userID, err := dbHelper.ValidateSession(database.Todo, sessionID)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	todoID := r.URL.Query().Get("id") // or from URL param if using a router like chi
	if todoID == "" {
		http.Error(w, "Missing todo ID", http.StatusBadRequest)
		return
	}

	todo, err := dbHelper.GetTodoByID(database.Todo, todoID, userID)
	if err != nil {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(todo)
}

func DeleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
		return
	}
	sessionID := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))

	userID, err := dbHelper.ValidateSession(database.Todo, sessionID)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	todoID := r.URL.Query().Get("id") // or from path param
	if todoID == "" {
		http.Error(w, "Missing todo ID", http.StatusBadRequest)
		return
	}

	err = dbHelper.ArchiveTodo(database.Todo, todoID, userID)
	if err != nil {
		http.Error(w, "Failed to archive todo", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Todo archived"))
}
