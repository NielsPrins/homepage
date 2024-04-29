package templatehandler

import (
	"embed"
	"github.com/gofiber/fiber/v2"
	"html/template"
)

func RenderTemplate(c *fiber.Ctx, htmlTemplate embed.FS, templateName string, data any) error {
	tmplData, err := htmlTemplate.ReadFile(templateName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error reading template file")
	}

	tmpl, err := template.New(templateName).Parse(string(tmplData))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error parsing template")
	}

	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
	err = tmpl.ExecuteTemplate(c, templateName, data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error executing template")
	}

	return nil
}
