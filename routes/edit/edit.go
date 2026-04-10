package edit

import (
	"homepage/common"
	"homepage/routes/add"
	"html/template"

	"github.com/gofiber/fiber/v2"
)

func Handler(c *fiber.Ctx) error {
	shortcut, err := common.GetShortcut(c.Params("id"))
	if err != nil {
		return err
	}

	return common.RenderTemplate(c, add.HtmlTemplate, add.Data{
		Title:        "Edit",
		SubmitLabel:  "Save",
		FormAction:   "/edit/" + shortcut.ID,
		RemoveAction: "/remove/" + shortcut.ID,
		Name:         shortcut.Name,
		URL:          shortcut.URL,
		Section:      shortcut.Section,
		ImageURL:     template.URL(shortcut.ImageURL),
		AutofocusURL: true,
	})
}

func PostHandler(c *fiber.Ctx) error {
	name := c.FormValue("name")
	url := c.FormValue("url")
	section := c.FormValue("section")

	customImage, err := add.ParseCustomImage(c)
	if err != nil {
		return err
	}

	err = common.EditShortcut(c.Params("id"), url, name, section, customImage)
	if err != nil {
		return err
	}

	return c.Redirect("/")
}

func RemoveHandler(c *fiber.Ctx) error {
	err := common.RemoveShortcut(c.Params("id"))
	if err != nil {
		return err
	}

	return c.Redirect("/")
}
