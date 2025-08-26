package handlers

import (
	"go-auth-backend/auth/services"

	"github.com/gofiber/fiber/v2"
	"fmt"
	"strings"
)

type AuthHandler struct {
	Service *services.AuthService
}

func NewAuthHandler(s *services.AuthService) *AuthHandler {
	return &AuthHandler{s}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	body := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name	 string `json:"name,omitempty"`
	}{}
	if err := c.BodyParser(&body); err != nil {
		fmt.Println("Error parsing body:", err)
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	if err := h.Service.Register(body.Email, body.Password, body.Name); err != nil {
		fmt.Println("Error registering user:", err)
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "user created"})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	body := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	token, err := h.Service.Login(body.Email, body.Password)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"token": token})
}

func (h *AuthHandler) GetUser(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(401).JSON(fiber.Map{"error": "Missing Authorization header"})
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenStr == authHeader {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid Authorization header"})
	}

	user, err := h.Service.GetUserByToken(tokenStr)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid or expired token"})
	}

	return c.JSON(fiber.Map{
		"email": user.Email,
		"name":  user.Name,
	})
}