package edit

import (
	"homepage/common"
	"homepage/routes/add"

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
		Name:         shortcut.Name,
		URL:          shortcut.URL,
		AutofocusURL: true,
	})
}

func PostHandler(c *fiber.Ctx) error {
	postData, err := add.ParseShortcutForm(c)
	if err != nil {
		return err
	}

	err = common.EditShortcut(c.Params("id"), postData.URL, postData.Name)
	if err != nil {
		return err
	}

	return c.Redirect("/")
}
