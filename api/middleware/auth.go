package middleware

import (
	"lot/config"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func Protected() fiber.Handler {
	secret, _ := config.Config("secret")
	return jwtware.New(
		jwtware.Config{
			SigningKey: jwtware.SigningKey{
				Key: []byte(secret),
			},
			ErrorHandler: jwtError,
		},
	)
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Missing or malformed JWT",
				"data":    nil,
			},
		)

	}
	return c.Status(fiber.StatusUnauthorized).JSON(
		fiber.Map{
			"status":  "error",
			"message": "Invalid or expired JWT",
			"data":    nil,
		},
	)
}
