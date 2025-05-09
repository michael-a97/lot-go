package route

import (
	"lot/api/handler"
	"lot/internal/service"

	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(router fiber.Router, authService service.AuthService) {
	router.Post("/signin", handler.SignInHandler(authService))
	router.Post("/reset-password", handler.ResetPasswordHandler(authService))
	router.Post("/refresh", handler.RefreshTokenHandler(authService))
}
