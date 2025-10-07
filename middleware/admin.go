package middleware

import (
	"apps2pay/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" {
			return c.Status(401).JSON(fiber.Map{"error": "Missing Authorization header"})
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		claims, err := utils.ValidateJWTWithRole(tokenStr)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid or expired token"})
		}

		if role, ok := claims["role"].(string); ok && role == "admin" {
			return c.Next()
		}

		return c.Status(403).JSON(fiber.Map{"error": "Forbidden: admin only"})
	}
}
