package errorpage

import (
	_ "embed"
	"errors"
	"homepage/common"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

//go:embed errorpage.gohtml
var htmlTemplate string

type Data struct {
	StatusCode int
	StatusText string
	Title      string
	Detail     string
}

func Render(c *fiber.Ctx, err error) error {
	statusCode := fiber.StatusInternalServerError

	var fiberError *fiber.Error
	if errors.As(err, &fiberError) {
		statusCode = fiberError.Code
	}

	data := Data{
		StatusCode: statusCode,
		StatusText: fiber.ErrInternalServerError.Message,
		Title:      "Something went wrong",
	}

	switch statusCode {
	case fiber.StatusNotFound:
		data.StatusText = fiber.ErrNotFound.Message
		data.Title = "Page not found"
	case fiber.StatusBadRequest:
		data.StatusText = fiber.ErrBadRequest.Message
		data.Title = "Bad request"
	default:
		data.StatusText = http.StatusText(statusCode)
		data.Detail = err.Error()
	}

	return common.RenderTemplateWithStatus(c, statusCode, htmlTemplate, data)
}
