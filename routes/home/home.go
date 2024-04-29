package home

import (
	_ "embed"
	"github.com/gofiber/fiber/v2"
	templatehandler "homepage/common"
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

	return templatehandler.RenderTemplate(c, htmlTemplate, data)
}
