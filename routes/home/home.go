package home

import (
	_ "embed"
	"github.com/gofiber/fiber/v2"
	"homepage/common"
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
