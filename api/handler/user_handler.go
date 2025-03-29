package handler

import (
	"lot/api/dto"
	app_errors "lot/internal/errors"
	"lot/internal/service"

	"github.com/gofiber/fiber/v2"
)

func SignUpHandler(userService service.UserService, authService service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var request dto.SignUpRequest
		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				dto.ApiResponse{
					Status: fiber.StatusBadRequest,
					Error:  err.Error(),
					Data:   nil,
				},
			)
		}

		if err := request.Validate(); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				dto.ApiResponse{
					Status: fiber.StatusBadRequest,
					Error:  err,
					Data:   nil,
				},
			)
		}

		isTokenValid, err := authService.VerifyPhoneNumberAuthenticationToken(request.PhoneNumberVerificationToken)

		if err != nil || !isTokenValid {
			return c.Status(fiber.StatusBadRequest).JSON(
				dto.ApiResponse{
					Status: fiber.StatusBadRequest,
					Data:   nil,
					Error:  app_errors.ErrInvalidPhoneNumberVerificationToken,
				},
			)
		}

		user, err := userService.SignUp(request)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				dto.ApiResponse{
					Status: fiber.StatusBadRequest,
					Error:  err.Error(),
					Data:   nil,
				},
			)
		}
		return c.JSON(user)
	}
}
