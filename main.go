package main

import (
	"embed"
	"homepage/routes/add"
	"homepage/routes/edit"
	"homepage/routes/errorpage"
	"homepage/routes/home"
	"homepage/routes/order"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
)

//go:embed public
var res embed.FS

func main() {
	app := fiber.New(fiber.Config{
		ErrorHandler: errorpage.Render,
	})

	app.Use("/public", func(c *fiber.Ctx) error {
		file := strings.TrimLeft(c.Path(), "/")
		content, err := res.ReadFile(file)
		if err != nil {
			return c.SendStatus(fiber.StatusNotFound)
		}

		c.Set("Content-type", getContentType(file))
		return c.Send(content)
	})

	app.Get("/", home.Handler)
	app.Get("/add", add.Handler)
	app.Post("/add", add.PostHandler)
	app.Get("/edit/:id", edit.Handler)
	app.Post("/edit/:id", edit.PostHandler)
	app.Post("/remove/:id", edit.RemoveHandler)
	app.Post("/order", order.PostHandler)

	port := ":4700"
	if os.Getenv("APP_ENV") == "dev" {
		port = ":4600"
	}
	log.Fatal(app.Listen(port))
}

func getContentType(filename string) string {
	switch {
	case strings.HasSuffix(filename, ".css"):
		return "text/css"
	case strings.HasSuffix(filename, ".js"):
		return "application/javascript"
	default:
		return "text/plain"
	}
}
