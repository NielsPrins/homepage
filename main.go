package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"log"
)

type Data struct {
	Title string
}

func main() {
	engine := html.New("./views", ".gohtml")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/public", "./public")

	app.Get("/", func(c *fiber.Ctx) error {
		data := Data{
			Title: "Go setup",
		}

		return c.Render("index", data)
	})

	log.Fatal(app.Listen(":3000"))
}
