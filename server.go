package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Habit represents a user's habit
type Habit struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	CreatedAt   time.Time   `json:"created_at"`
	Logs        []time.Time `json:"logs,omitempty"` // Timestamps when habit was logged
}

// In-memory storage for habits (thread-safe)
var (
	habits   = make(map[string]*Habit)
	habitsMu sync.RWMutex // Mutex to protect concurrent access
)

// Middleware to mock JWT authentication
func jwtAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		// Mock: Accept any non-empty Bearer token
		if !strings.HasPrefix(auth, "Bearer ") || len(auth) <= 7 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// In real code, validate JWT here!
		next.ServeHTTP(w, r)
	})
}

// Health check endpoint
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// Hello endpoint (public)
func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello from Habit Tracker!"))
}

// List all habits (GET) or create a new habit (POST)
func habitsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// List all habits
		habitsMu.RLock()
		defer habitsMu.RUnlock()
		list := make([]*Habit, 0, len(habits))
		for _, h := range habits {
			list = append(list, h)
		}
		json.NewEncoder(w).Encode(list)
	case http.MethodPost:
		// Create a new habit
		var input struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		if input.Name == "" {
			http.Error(w, "Name is required", http.StatusBadRequest)
			return
		}
		id := fmt.Sprintf("%d", time.Now().UnixNano())
		habit := &Habit{
			ID:          id,
			Name:        input.Name,
			Description: input.Description,
			CreatedAt:   time.Now(),
		}
		habitsMu.Lock()
		habits[id] = habit
		habitsMu.Unlock()
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(habit)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Get, update, or delete a specific habit by ID
func habitByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/habits/")
	habitsMu.RLock()
	habit, exists := habits[id]
	habitsMu.RUnlock()
	if !exists {
		http.Error(w, "Habit not found", http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		json.NewEncoder(w).Encode(habit)
	case http.MethodPut:
		var input struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		habitsMu.Lock()
		habit.Name = input.Name
		habit.Description = input.Description
		habitsMu.Unlock()
		json.NewEncoder(w).Encode(habit)
	case http.MethodDelete:
		habitsMu.Lock()
		delete(habits, id)
		habitsMu.Unlock()
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Log a completion for a habit (POST /api/habits/{id}/log)
func habitLogHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	id := parts[3]
	habitsMu.Lock()
	habit, exists := habits[id]
	if !exists {
		habitsMu.Unlock()
		http.Error(w, "Habit not found", http.StatusNotFound)
		return
	}
	habit.Logs = append(habit.Logs, time.Now())
	habitsMu.Unlock()
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(habit)
}

// Get all logs for a habit (GET /api/habits/{id}/logs)
func habitLogsHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	id := parts[3]
	habitsMu.RLock()
	habit, exists := habits[id]
	habitsMu.RUnlock()
	if !exists {
		http.Error(w, "Habit not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(habit.Logs)
}

func main() {
	// Public endpoints
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/hello", helloHandler)

	// Protected endpoints (require JWT)
	http.Handle("/api/habits", jwtAuthMiddleware(http.HandlerFunc(habitsHandler)))
	http.Handle("/api/habits/", jwtAuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Route based on path
		if strings.HasSuffix(r.URL.Path, "/log") && r.Method == http.MethodPost {
			habitLogHandler(w, r)
			return
		}
		if strings.HasSuffix(r.URL.Path, "/logs") && r.Method == http.MethodGet {
			habitLogsHandler(w, r)
			return
		}
		habitByIDHandler(w, r)
	})))

	fmt.Println("Backend API server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

/*
How to run:
1. Save this file as server.go
2. Run: go run server.go
3. Test endpoints:
   - GET  http://localhost:8080/health
   - GET  http://localhost:8080/hello
   - For /api/habits* endpoints, add header: Authorization: Bearer testtoken
   - Example: curl -H "Authorization: Bearer testtoken" http://localhost:8080/api/habits
*/
