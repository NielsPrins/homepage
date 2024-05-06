package home

import (
	_ "embed"
	"github.com/gofiber/fiber/v2"
	"homepage/common"
)

//go:embed home.gohtml
var htmlTemplate string

type Data struct {
	Title string
}

func Handler(c *fiber.Ctx) error {
	data := Data{
		Title: "Go setup",
	}

	return common.RenderTemplate(c, htmlTemplate, data)
}
