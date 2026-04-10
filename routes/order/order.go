package order

import (
	_ "embed"
	"homepage/common"

	"github.com/gofiber/fiber/v2"
)

type ReorderRequest struct {
	Items []common.OrderShortcutDto `json:"items"`
}

func PostHandler(c *fiber.Ctx) error {
	var request ReorderRequest

	err := c.BodyParser(&request)
	if err != nil {
		return err
	}

	err = common.ReorderShortcuts(request.Items)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}
