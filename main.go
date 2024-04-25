package main

import (
	"github.com/gofiber/fiber/v2"
	"homepage/routes/home"
	"log"
)

func main() {
	app := fiber.New()

	app.Static("/public", "./public")

	app.Get("/", home.Handler)

	log.Fatal(app.Listen(":3000"))
}
