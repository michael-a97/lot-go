package handler

import (
	dto "lot/api/dto/auth"
	"lot/pkg/service/auth"

	"github.com/gofiber/fiber/v2"
)

func SignInHandler(authService service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var request dto.LoginRequest

		if err := c.BodyParser(&request); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		if err := request.Validate(); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				fiber.Map{
					"status": "error",
					"error":  err,
					"data":   nil,
				},
			)
		}

		result, err := authService.SignIn(request)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				fiber.Map{
					"status": "error",
					"data":   nil,
					"error":  err.Error(),
				},
			)
		}

		return c.JSON(
			fiber.Map{
				"status": "success",
				"data":   result,
				"error":  nil,
			},
		)
	}
}

func RefreshTokenHandler(authService service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var request dto.TokenRefreshRequest

		if err := c.BodyParser(&request); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		if err := request.Validate(); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				fiber.Map{
					"status": "error",
					"error":  err,
					"data":   nil,
				},
			)
		}

		result, err := authService.RefreshToken(request)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				fiber.Map{
					"status": "error",
					"data":   nil,
					"error":  err.Error(),
				},
			)
		}

		return c.JSON(
			fiber.Map{
				"status": "success",
				"data":   result,
				"error":  nil,
			},
		)
	}
}
