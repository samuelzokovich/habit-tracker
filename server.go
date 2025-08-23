package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/microsoft/go-mssqldb" // Azure SQL driver
)

// AppName is the application name used in logs
const AppName = "Habit Tracker"

// DBConfig holds database credentials (dummy for now)
type DBConfig struct {
	Server   string
	Port     int
	User     string
	Password string
	Database string
}

// connectToDB initializes a SQL connection (dummy for now)
func connectToDB(cfg DBConfig) (*sql.DB, error) {
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s;port=%d",
		cfg.Server, cfg.User, cfg.Password, cfg.Database, cfg.Port)

	// Note: In real production code, use managed identity / secrets manager instead of raw credentials.
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		return nil, err
	}
	// Ping to validate connection (here it will fail unless real DB is provided)
	if err = db.Ping(); err != nil {
		log.Printf("[WARN] DB connection test failed: %v", err)
	}
	return db, nil
}

// loggingMiddleware logs request details
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("[INFO] Client: %s | Method: %s | URL: %s",
			r.RemoteAddr, r.Method, r.URL.Path)

		next.ServeHTTP(w, r)

		log.Printf("[INFO] Completed in %v", time.Since(start))
	})
}

// healthHandler checks server health
func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK - %s server is healthy\n", AppName)
}

// helloHandler simple hello endpoint
func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from %s server!\n", AppName)
}

func main() {
	// Load dummy DB config (in production: load from ENV or Vault)
	dbConfig := DBConfig{
		Server:   "dummy-sql-server.database.windows.net",
		Port:     1433,
		User:     "sqladmin",
		Password: "SuperSecretPassword123",
		Database: "HabitTrackerDB",
	}

	// Try DB connection (won’t work with dummy values, just logs warning)
	_, _ = connectToDB(dbConfig)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/hello", helloHandler)

	loggedMux := loggingMiddleware(mux)

	port := "8080"
	if val := os.Getenv("PORT"); val != "" {
		port = val
	}

	log.Printf("[START] %s server running on port %s", AppName, port)
	if err := http.ListenAndServe(":"+port, loggedMux); err != nil {
		log.Fatalf("[FATAL] Server failed: %v", err)
	}
}
