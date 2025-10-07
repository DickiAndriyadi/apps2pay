package handlers

import (
	"apps2pay/models"
	"apps2pay/utils"
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

var DB *pgx.Conn
var redisClient *redis.Client

func SetDB(conn *pgx.Conn, redis *redis.Client) {
	DB = conn
	redisClient = redis
}

func Login(c *fiber.Ctx) error {
	var req models.User
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	var storedHash, role string
	err := DB.QueryRow(context.Background(),
		"SELECT password_hash, role FROM users WHERE email = $1", req.Email).
		Scan(&storedHash, &role)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(req.Password)); err == nil {
		token, _ := utils.GenerateJWT(req.Email, role)
		return c.JSON(fiber.Map{"token": token, "role": role})
	}

	return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
}
