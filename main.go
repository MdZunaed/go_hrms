package main

import (
	"log"

	"github.com/MdZunaed/go_hrms/internal/db"
	employee "github.com/MdZunaed/go_hrms/internal/handler"
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

	app.Route("/employee", func(emp fiber.Router) {
		emp.Get("/", employee.GetEmployees())
		emp.Post("/", mw.VerifyToken(), employee.CreateEmployee())
		emp.Put("/:id", employee.UpdateEmployee())
		emp.Delete("/:id", employee.DeleteEmployee())
	})

	app.Get("/download", file.DownloadFile())
	app.Post("/upload", file.UploadMultiFile())

	app.Listen(":3000")
}
