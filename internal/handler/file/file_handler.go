package file

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/MdZunaed/go_hrms/internal/db"
	"github.com/gofiber/fiber/v2"
)

func DownloadFile() fiber.Handler {
	return func(c *fiber.Ctx) error {
		basePath, erro := os.Getwd()
		if erro != nil {
			return c.JSON(fiber.Map{"error": erro.Error()})
		}
		filePath := filepath.Join(basePath, "internal", "db", "test_file.txt")
		err := c.Download(filePath)
		if err != nil {
			return c.JSON(fiber.Map{"error": err.Error()})
		}
		return nil
	}
}

func UploadFile() fiber.Handler {
	return func(c *fiber.Ctx) error {
		file, err := c.FormFile("doc")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		tempPath, _ := db.GetUploadDir()
		err = c.SaveFile(file, filepath.Join(tempPath, file.Filename))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"Message": "Success"})
	}
}

func UploadMultiFile() fiber.Handler {
	return func(c *fiber.Ctx) error {
		form, err := c.MultipartForm()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		files := form.File["docs"]
		for _, file := range files {
			fmt.Println(file.Filename, file.Size, file.Header["Content-Type"][0])
			tempPath, _ := db.GetUploadDir()
			if err := c.SaveFile(file, filepath.Join(tempPath, file.Filename)); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
			}
		}
		return c.JSON(fiber.Map{"Message": "Success"})
	}
}
