package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/gunsu12/go-wallet/internal/db"
	"github.com/gunsu12/go-wallet/internal/delivery"
	"github.com/gunsu12/go-wallet/internal/repository"
	"github.com/gunsu12/go-wallet/internal/usecase"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Gagal load file .env")
	}

	// Inisialisasi koneksi DB
	database := db.InitMySQL()

	// Jalankan migrasi
	db.RunMigrations(database)

	// Inisialisasi komponen
	userRepo := repository.NewUserRepository(database)
	userUsecase := usecase.NewUserUsecase(userRepo)

	// Router
	r := gin.Default()
	delivery.NewUserHandler(r, userUsecase)

	// Jalankan server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
