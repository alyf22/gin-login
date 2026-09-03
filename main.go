package main

import (
	"context"
	"log"
	"net/http"
	"os"

	auth "gin-login/Auth"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	db, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("connect to database: %v", err)
	}
	defer db.Close()
	if err := db.Ping(context.Background()); err != nil {
		log.Fatalf("ping database: %v", err)
	}

	repository := auth.NewAuthRepository(db)
	service := auth.NewAuthService(repository, jwtSecret)
	handler := auth.NewAuthHandler(service)

	router := gin.Default()
	router.POST("/login", handler.Login)

	router.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, Gin!",
		})
	})

	router.Run(":8080")
}
