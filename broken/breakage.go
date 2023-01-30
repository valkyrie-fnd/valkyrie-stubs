package broken

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func breakage(c *fiber.Ctx, e string) error {
	switch e {
	case "timeout 5s":
		// This is a timeout error
		time.Sleep(5 * time.Second)
		return c.SendStatus(fiber.StatusRequestTimeout)
	case "connection close":
		// This is a connection close error
		return c.Context().Conn().Close()
	default:
		return nil
	}
}
