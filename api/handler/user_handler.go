package handler

import (
	dto "lot/api/dto/user"
	"lot/pkg/service"

	"github.com/gofiber/fiber/v2"
)

func SignUpHandler(userService service.UserService, authService service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var request dto.SignUpRequest
		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				fiber.Map{
					"status": "error",
					"error":  err.Error(),
					"data":   nil,
				},
			)
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

		isTokenValid, err := authService.VerifyPhoneNumberAuthenticationToken(request.PhoneNumberVerificationToken)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				fiber.Map{
					"status": "error",
					"data":   nil,
					"error":  "invalid phone number verification token",
				},
			)
		}

		if !isTokenValid {
			return c.Status(fiber.StatusBadRequest).JSON(
				fiber.Map{
					"status": "error",
					"data":   nil,
					"error":  "verification failed",
				},
			)
		}

		user, err := userService.SignUp(request)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				fiber.Map{
					"status": "error",
					"error":  err.Error(),
					"data":   nil,
				},
			)
		}
		return c.JSON(user)
	}
}
