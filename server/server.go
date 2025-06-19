package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"todo-app/handlers"
)

func SetupRoutes() http.Handler {
	r := mux.NewRouter()

	// Auth routes
	r.HandleFunc("/register", handlers.RegisterHandler).Methods("POST")
	r.HandleFunc("/login", handlers.LoginHandler).Methods("POST")
	r.HandleFunc("/logout", handlers.LogoutHandler).Methods("POST")
	r.HandleFunc("/CreateTodo", handlers.CreateTodoHandler).Methods("POST")
	r.HandleFunc("/UpdateTodo", handlers.UpdateTodoHandler).Methods("PUT")
	r.HandleFunc("/GetTodo", handlers.GetTodoHandler).Methods("GET")          // expects ?id=<todo_id>
	r.HandleFunc("/GetAllTodos", handlers.GetAllTodosHandler).Methods("GET")  // fetch all for user
	r.HandleFunc("/DeleteTodo", handlers.DeleteTodoHandler).Methods("DELETE") // expects ?id=<todo_id>

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(map[string]bool{"ok": true})
		if err != nil {
			logrus.Errorf("Error encoding response: %v", err)
		}
	})

	return r
}
