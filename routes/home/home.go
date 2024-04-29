package home

import (
	"embed"
	"github.com/gofiber/fiber/v2"
	"html/template"
)

//go:embed home.gohtml
var htmlTemplate embed.FS

type Data struct {
	Title string
}

func Handler(c *fiber.Ctx) error {
	data := Data{
		Title: "Go setup!!",
	}

	tmplData, err := htmlTemplate.ReadFile("home.gohtml")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error reading template file")
	}

	tmpl, err := template.New("home").Parse(string(tmplData))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error parsing template")
	}

	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
	err = tmpl.ExecuteTemplate(c, "home", data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error executing template")
	}

	return nil
}
