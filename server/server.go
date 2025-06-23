package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"todo-app/handlers"
	"todo-app/middleware"
)

func SetupRoutes() http.Handler {
	r := mux.NewRouter()

	// Auth routes
	r.HandleFunc("/register", handlers.RegisterHandler).Methods("POST")
	r.HandleFunc("/login", handlers.LoginHandler).Methods("POST")
	r.HandleFunc("/logout", handlers.LogoutHandler).Methods("POST")

	// Protected routes - wrap with JWT middleware
	protected := r.PathPrefix("/my").Subrouter()
	protected.Use(middleware.AuthMiddleware)

	protected.HandleFunc("/CreateTodo", handlers.CreateTodoHandler).Methods("POST")
	protected.HandleFunc("/UpdateTodo", handlers.UpdateTodoHandler).Methods("PUT")
	protected.HandleFunc("/GetTodo", handlers.GetTodoHandler).Methods("GET")          // expects ?id=<todo_id>
	protected.HandleFunc("/GetAllTodos", handlers.GetAllTodosHandler).Methods("GET")  // fetch all for user
	protected.HandleFunc("/DeleteTodo", handlers.DeleteTodoHandler).Methods("DELETE") // expects ?id=<todo_id>

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(map[string]bool{"ok": true})
		if err != nil {
			logrus.Errorf("Error encoding response: %v", err)
		}
	})

	return r
}
