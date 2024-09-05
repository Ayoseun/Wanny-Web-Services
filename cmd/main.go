// main.go
package main

import (
	"log"
	"net/http"
	"wanny-web-services/config"
	"wanny-web-services/internal/adapters/postgres"
	"wanny-web-services/internal/adapters/web"
	"wanny-web-services/internal/core/services"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	db, err := postgres.NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	userRepo := postgres.NewUserRepository(db)
	fileRepo := postgres.NewFileRepository(db)

	userService := services.NewUserService(userRepo)
	fileService := services.NewFileService(fileRepo, cfg.MaxStoragePerUser, cfg.EncryptionKey)

	handler := web.NewHandler(userService, fileService, cfg.JWTSecret)

	log.Printf("Server starting on %s", cfg.ServerAddress)
	log.Fatal(http.ListenAndServe(cfg.ServerAddress, handler.Router))
}

// config/config.go

// internal/core/domain/user.go

// internal/core/domain/file.go

// internal/core/ports/user_repository.go

// internal/core/ports/file_repository.go

// internal/core/services/user_service.go

// internal/core/services/file_service.go

// internal/adapters/postgres/db.go

// internal/adapters/postgres/user_repository.go

// internal/adapters/postgres/file_repository.go

// internal/adapters/http/handler.go
// internal/adapters/http/handler.go
