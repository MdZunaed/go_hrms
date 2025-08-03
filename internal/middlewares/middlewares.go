package mw

import (
	"fmt"
	"strings"

	"github.com/MdZunaed/go_hrms/internal/db"
	"github.com/MdZunaed/go_hrms/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func VerifyToken() fiber.Handler {
	var tokenString string
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
		tokenString = tokenParts[1]

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
			if t.Method.Alg() != jwt.GetSigningMethod("HS256").Alg() {
				return nil, fmt.Errorf("unexpected signin method: %v", t.Header["alg"])
			}
			return utils.JWTSecret, nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid or expired token"})
		}
		userId := token.Claims.(jwt.MapClaims)["userId"]
		objectId, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", userId))
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid userId format"})
		}

		var user bson.M
		err = db.GetUserCollection().FindOne(c.Context(), fiber.Map{"_id": objectId}).Decode(&user)

		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "User not found"})
		}
		c.Locals("userId", userId)
		return c.Next()
	}
}
