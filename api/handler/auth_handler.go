package handler

import (
	dto "lot/api/dto"
	app_errors "lot/internal/errors"
	"lot/internal/service"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func SignInHandler(authService service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var request dto.LoginRequest

		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				dto.ApiResponse{
					Status: 400,
					Error:  err,
					Data:   nil,
				},
			)
		}

		if err := request.Validate(); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				dto.ApiResponse{
					Status: 400,
					Error:  err,
					Data:   nil,
				},
			)
		}

		result, err := authService.SignIn(request)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				dto.ApiResponse{
					Status: 400,
					Data:   nil,
					Error:  err.Error(),
				},
			)
		}

		return c.JSON(
			dto.ApiResponse{
				Status: 200,
				Data:   result,
				Error:  nil,
			},
		)
	}
}

func RefreshTokenHandler(authService service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var request dto.TokenRefreshRequest

		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				dto.ApiResponse{
					Status: fiber.StatusBadRequest,
					Error:  err,
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

		result, err := authService.RefreshToken(request)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				dto.ApiResponse{
					Status: fiber.StatusBadRequest,
					Data:   nil,
					Error:  err.Error(),
				},
			)
		}

		return c.JSON(
			dto.ApiResponse{
				Status: fiber.StatusAccepted,
				Data:   result,
				Error:  nil,
			},
		)
	}
}

func ResetPasswordHandler(authService service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var passwordResetRequest dto.PasswordResetRequest
		if err := c.BodyParser(&passwordResetRequest); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				dto.ApiResponse{
					Status: fiber.StatusBadRequest,
					Error:  err.Error(),
					Data:   nil,
				},
			)
		}
		if err := passwordResetRequest.Validate(); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				dto.ApiResponse{
					Status: fiber.StatusBadRequest,
					Error:  err.Error(),
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
		if err = authService.ResetPassword(passwordResetRequest); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				dto.ApiResponse{
					Status: fiber.StatusInternalServerError,
					Data:   nil,
					Error:  err,
				},
			)
		}
		return nil
	}
}
