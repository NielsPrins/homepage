package home

import (
	"github.com/gofiber/fiber/v2"
)

type Data struct {
	Title string
}

func Handler(c *fiber.Ctx) error {

	data := Data{
		Title: "Go setup!!",
	}

	return c.Render("routes/home/home.gohtml", data)
}
