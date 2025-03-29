package route

import (
	"github.com/gofiber/fiber/v2"
	"lot/api/handler"
	authService "lot/pkg/service/auth"
	userService "lot/pkg/service/user"
)

func SetupUserRoutes(
	router fiber.Router,
	userService userService.UserService,
	authService authService.AuthService,
) {
	router.Post("/signup", handler.SignUpHandler(userService, authService))
}
