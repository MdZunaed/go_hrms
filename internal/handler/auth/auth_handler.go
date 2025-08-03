package auth

import (
	"github.com/MdZunaed/go_hrms/internal/db"
	model "github.com/MdZunaed/go_hrms/internal/models"
	"github.com/MdZunaed/go_hrms/internal/utils"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func Signup() fiber.Handler {
	return func(c *fiber.Ctx) error {

		user := &model.User{
			Email:    c.FormValue("email"),
			Password: c.FormValue("password"),
		}
		if user.Email == "" || user.Password == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Email and Password is required",
			})
		}

		hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		user.Password = string(hashed)

		result, err := db.GetUserCollection().InsertOne(c.Context(), user)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		token, err := utils.GenerateToken(user)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"id":    result.InsertedID,
			"token": token,
		})
	}
}

func Login() fiber.Handler {
	return func(c *fiber.Ctx) error {

		authUser := &model.User{
			Email:    c.FormValue("email"),
			Password: c.FormValue("password"),
		}
		if authUser.Email == "" || authUser.Password == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Email and Password is required",
			})
		}

		var dbUser model.User
		query := fiber.Map{
			"email": authUser.Email,
		}
		err := db.GetUserCollection().FindOne(c.Context(), query).Decode(&dbUser)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}

		err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(authUser.Password))

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Invalid password",
			})
		}

		token, err := utils.GenerateToken(&dbUser)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"id":    dbUser.Id,
			"token": token,
		})
	}
}

func TestMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Token is valid"})
	}
}
