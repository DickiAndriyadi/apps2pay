package middleware

import (
	"apps2pay/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" {
			return c.Status(401).JSON(fiber.Map{"error": "Missing Authorization header"})
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		_, err := utils.ValidateJWTWithRole(tokenStr)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid or expired token"})
		}

		return c.Next()
	}
}
