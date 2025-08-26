package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"go-auth-backend/auth/config"
	"go-auth-backend/auth/db"
	"go-auth-backend/auth/handler"
	"go-auth-backend/auth/repository"
	"go-auth-backend/auth/services"
)

func main() {
	cfg := config.LoadConfig()

	db.ConnectMongo(cfg)
	fmt.Println("✅ MongoDB connected")

	userRepo := repository.NewUserRepository(db.Database)
	authService := services.NewAuthService(userRepo, cfg)
	authHandler := handlers.NewAuthHandler(authService)

	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New())

	api := app.Group("/api")
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Get("/me", authHandler.GetUser)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("natthawat")
	})

	port := cfg.Port
	if port == "" {
		port = "8080"
	}
	log.Println("Server running on port", port)
	log.Fatal(app.Listen(":" + port))
}
