package add

import (
	_ "embed"
	"github.com/gofiber/fiber/v2"
	"homepage/common"
)

//go:embed add.gohtml
var htmlTemplate string

type Data struct {
}

func Handler(c *fiber.Ctx) error {
	return common.RenderTemplate(c, htmlTemplate, fiber.Map{})
}

func PostHandler(c *fiber.Ctx) error {
	postData := struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}{}

	if err := c.BodyParser(&postData); err != nil {
		return err
	}

	_, err := common.AddShortcut(postData.Url, postData.Name)
	if err != nil {
		return err
	}

	return c.Redirect("/")
}
