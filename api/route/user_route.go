package route

import (
	"lot/api/handler"
	"lot/api/middleware"
	"lot/internal/service"

	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(
	router fiber.Router,
	userService service.UserService,
	authService service.AuthService,
) {
	router.Post("/signup", handler.SignUpHandler(userService, authService))
	router.Get("/me",
		middleware.ProtectedMiddleware(),
		middleware.AuthenticationMiddleware(authService),
		handler.FindUserHandler(userService),
	)
}
