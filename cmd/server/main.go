package main

import (
	"log"
	"os"
	"some-pet/internal/database"
	"some-pet/internal/handlers"
	"some-pet/internal/repository"
	"some-pet/internal/routes"
	"some-pet/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {

	db := database.NewPostgres()

	database.RunMigrations()

	bookRepo := repository.NewBooks(db)
	bookService := service.NewBooks(bookRepo)
	bookHandler := handlers.NewBooks(bookService)

	r := gin.Default()

	routes.RegisterRouter(r, bookHandler)
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("server started on port", port)

	r.Run(":" + port)
}
