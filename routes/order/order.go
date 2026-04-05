package order

import (
	_ "embed"
	"homepage/common"

	"github.com/gofiber/fiber/v2"
)

type ReorderRequest struct {
	IDs []string `json:"ids"`
}

func PostHandler(c *fiber.Ctx) error {
	var request ReorderRequest

	err := c.BodyParser(&request)
	if err != nil {
		return err
	}

	err = common.ReorderShortcuts(request.IDs)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}
