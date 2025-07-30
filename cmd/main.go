package main

import (
	"log"

	"github.com/MdZunaed/go_hrms/internal/db"
	employee "github.com/MdZunaed/go_hrms/internal/handler"
	"github.com/gofiber/fiber/v2"
)

func main() {
	if err := db.ConnectDb(); err != nil {
		log.Fatal("connect DB error", err)
	}

	app := fiber.New()

	app.Get("/employee", employee.GetEmployees())
	app.Post("/employee", employee.CreateEmployee())
	app.Put("/employee/:id", employee.UpdateEmployee())
	app.Delete("/employee/:id", employee.DeleteEmployee())

	app.Listen(":3000")
}
