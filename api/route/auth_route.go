package route

import (
	"lot/api/handler"
	"lot/pkg/service"

	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(router fiber.Router, authService service.AuthService) {
	router.Post("/signin", handler.SignInHandler(authService))
	router.Post("/refresh", handler.RefreshTokenHandler(authService))
}
