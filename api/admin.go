package api

import (
	"github.com/GoDev/Hotel-reservatrion/types"
	"github.com/gofiber/fiber/v2"
)

func AdminAuth(c *fiber.Ctx) error {
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return ErrUnAthorized()
	}

	if !user.IsAdmin {
		return ErrUnAthorized()
	}
	return c.Next()
}
