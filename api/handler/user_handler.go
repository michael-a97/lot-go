package handler

import (
	"lot/api/dto"
	"lot/internal/errors"
	"lot/internal/service"
	"strings"

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
		token := strings.Split(c.Get("Authorization"), " ")[1]

		isTokenValid, err := authService.VerifyPhoneNumberAuthenticationToken(token)

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
		return c.JSON(
			dto.ApiResponse{
				Data: user, Status: 200, Message: "Success",
			},
		)
	}
}
