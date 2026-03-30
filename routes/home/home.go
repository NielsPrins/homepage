package home

import (
	_ "embed"
	"homepage/common"

	"github.com/gofiber/fiber/v2"
)

//go:embed home.gohtml
var htmlTemplate string

type Data struct {
	Shortcuts common.Shortcuts
}

func Handler(c *fiber.Ctx) error {
	shortcuts, _ := common.GetAllShortcuts()

	data := Data{
		Shortcuts: shortcuts,
	}

	return common.RenderTemplate(c, htmlTemplate, data)
}
