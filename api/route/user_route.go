package route

import (
	"github.com/gofiber/fiber/v2"
	"lot/api/handler"
	"lot/internal/service"
)

func SetupUserRoutes(
	router fiber.Router,
	userService service.UserService,
	authService service.AuthService,
) {
	router.Post("/signup", handler.SignUpHandler(userService, authService))
}
