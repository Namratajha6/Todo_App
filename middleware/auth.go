package middleware

import (
	"context"
	"net/http"
	"strings"
	"todo-app/database"
	"todo-app/database/dbHelper"
)

type contextKey string

const UserIDKey = contextKey("userID")

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get session ID from Authorization header
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

		// Inject userID into context
		ctx := context.WithValue(r.Context(), UserIDKey, userID)

		// Call the next handler with updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserID(r *http.Request) (string, error) {
	userID, ok := r.Context().Value(UserIDKey).(string)
	if !ok || userID == "" {
		return "", http.ErrNoCookie // or return errors.New("unauthorized")
	}
	return userID, nil
}
