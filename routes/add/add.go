package add

import (
	_ "embed"
	"encoding/base64"
	"fmt"
	"homepage/common"
	"html/template"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

//go:embed add.gohtml
var HtmlTemplate string

const MaxImageSize = 200 * 1024 // 200KB

type Data struct {
	Title        string
	SubmitLabel  string
	FormAction   string
	RemoveAction string
	Name         string
	URL          string
	ImageURL     template.URL
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
	name := c.FormValue("name")
	url := c.FormValue("url")

	customImage, err := ParseCustomImage(c)
	if err != nil {
		return err
	}

	_, err = common.AddShortcut(url, name, customImage)
	if err != nil {
		return err
	}

	return c.Redirect("/")
}

func ParseCustomImage(c *fiber.Ctx) (string, error) {
	file, err := c.FormFile("image")
	if err != nil || file == nil {
		return "", nil
	}

	if file.Size > MaxImageSize {
		return "", fmt.Errorf("image must be smaller than 200KB")
	}

	f, err := file.Open()
	if err != nil {
		return "", err
	}
	defer func() { _ = f.Close() }()

	data, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}

	mimeType := http.DetectContentType(data)
	encoded := base64.StdEncoding.EncodeToString(data)
	return fmt.Sprintf("data:%s;base64,%s", mimeType, encoded), nil
}
