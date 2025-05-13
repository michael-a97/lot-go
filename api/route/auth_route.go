package route

import (
	"lot/api/handler"
	"lot/api/middleware"
	"lot/internal/service"

	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(router fiber.Router, authService service.AuthService) {
	router.Post("/signin", handler.SignInHandler(authService))
	router.Post("/reset-password", handler.ResetPasswordHandler(authService))
	router.Post("/refresh", handler.RefreshTokenHandler(authService))
	router.Post(
		"/change-password",
		middleware.ProtectedMiddleware(),
		middleware.AuthenticationMiddleware(authService),
		handler.ChangePasswordHandler(authService),
	)
}
