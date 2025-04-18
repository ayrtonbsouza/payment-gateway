package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/ayrtonbsouza/payment-gateway/internal/repository"
	"github.com/ayrtonbsouza/payment-gateway/internal/service"
	"github.com/ayrtonbsouza/payment-gateway/internal/web/server"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		getEnv("DB_HOST", "db"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "postgres"),
		getEnv("DB_NAME", "payment_gateway"),
		getEnv("DB_SSL_MODE", "disable"),
	)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}

	defer db.Close()
	accountRepository := repository.NewAccountRepository(db)
	accountService := service.NewAccountService(accountRepository)
	port := getEnv("HTTP_PORT", "8080")
	srv := server.NewServer(accountService, port)
	srv.ConfigureRoutes()

	if err := srv.Start(); err != nil {
		log.Fatal("Error starting server: ", err)
	}
}

func getEnv(key string, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
