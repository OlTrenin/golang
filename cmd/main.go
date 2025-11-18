package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"number-service/internal/handler"
	"number-service/internal/repository"
	"number-service/internal/service"
)

func main() {
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "numbers_db")
	serverPort := getEnv("SERVER_PORT", "8080")

	var db *sql.DB
	var err error
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		db, err = repository.InitDB(dbHost, dbPort, dbUser, dbPassword, dbName)
		if err == nil {
			break
		}
		log.Printf("Failed to connect to database (attempt %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("Failed to initialize database after %d attempts: %v", maxRetries, err)
	}
	defer db.Close()

	log.Println("Database connection established successfully")

	repo := repository.NewPostgresRepository(db)
	svc := service.NewNumberService(repo)
	h := handler.NewNumberHandler(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("/numbers", h.AddNumber)
	mux.HandleFunc("/health", h.HealthCheck)

	addr := fmt.Sprintf(":%s", serverPort)
	log.Printf("Server starting on port %s", serverPort)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
