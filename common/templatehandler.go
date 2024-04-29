package templatehandler

import (
	"github.com/gofiber/fiber/v2"
	"html/template"
)

func RenderTemplate(c *fiber.Ctx, htmlTemplate string, data any) error {
	tmpl, err := template.New("gohtml-template").Parse(htmlTemplate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error parsing template")
	}

	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
	err = tmpl.ExecuteTemplate(c, "gohtml-template", data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error executing template")
	}

	return nil
}
