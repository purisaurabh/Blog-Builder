package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/purisaurabh/blog-backend/internal/helper"
)

func IsAuthenticated(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	if cookie == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	if _, err := helper.VerifyToken(cookie); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	return c.Next()
}
