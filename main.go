package main

import (
	"log"

	"github.com/MdZunaed/go_hrms/internal/db"
	"github.com/MdZunaed/go_hrms/internal/handler/auth"
	employee "github.com/MdZunaed/go_hrms/internal/handler/employee"
	"github.com/MdZunaed/go_hrms/internal/handler/file"
	mw "github.com/MdZunaed/go_hrms/internal/middlewares"
	"github.com/gofiber/fiber/v2"
)

func main() {
	if err := db.ConnectDb(); err != nil {
		log.Fatal("connect DB error", err)
	}

	app := fiber.New()

	app.Get("any/:any", func(c *fiber.Ctx) error {
		return c.SendString(c.Params("any"))
	})

	// Authentication

	app.Route("/auth", func(route fiber.Router) {
		route.Post("/signup", auth.Signup())
		route.Post("/login", auth.Login())
	})

	// Employee

	app.Route("/employee", func(emp fiber.Router) {
		emp.Get("/", employee.GetEmployees())
		emp.Post("/", mw.VerifyToken(), employee.CreateEmployee())
		emp.Put("/:id", mw.VerifyToken(), employee.UpdateEmployee())
		emp.Delete("/:id", mw.VerifyToken(), employee.DeleteEmployee())
	})

	// Protected

	admin := app.Group("/admin", mw.VerifyToken())
	admin.Get("verify", auth.TestMiddleware())

	app.Get("/download", file.DownloadFile())
	app.Post("/upload", file.UploadMultiFile())

	app.Listen(":3000")
}
