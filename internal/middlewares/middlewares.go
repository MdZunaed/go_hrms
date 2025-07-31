package mw

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func VerifyToken() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHader := c.Get("Authorization")
		if authHader == "" {
			log.Warnf("empty authorization header")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
		}
		tokenParts := strings.Split(authHader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			log.Warnf("invalid token parts")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid token parts"})
		}
		return c.Next()
	}
}
