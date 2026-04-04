package add

import (
	_ "embed"
	"homepage/common"

	"github.com/gofiber/fiber/v2"
)

//go:embed add.gohtml
var HtmlTemplate string

type Data struct {
	Title        string
	SubmitLabel  string
	FormAction   string
	RemoveAction string
	Name         string
	URL          string
	AutofocusURL bool
}

func Handler(c *fiber.Ctx) error {
	return common.RenderTemplate(c, HtmlTemplate, Data{
		Title:        "Add",
		SubmitLabel:  "Add",
		FormAction:   "/add",
		AutofocusURL: false,
	})
}

func PostHandler(c *fiber.Ctx) error {
	postData, err := ParseShortcutForm(c)
	if err != nil {
		return err
	}

	_, err = common.AddShortcut(postData.URL, postData.Name)
	if err != nil {
		return err
	}

	return c.Redirect("/")
}

type ShortcutFormData struct {
	Name string `json:"name" form:"name"`
	URL  string `json:"url" form:"url"`
}

func ParseShortcutForm(c *fiber.Ctx) (ShortcutFormData, error) {
	postData := ShortcutFormData{}

	err := c.BodyParser(&postData)
	if err != nil {
		return ShortcutFormData{}, err
	}

	return postData, nil
}
