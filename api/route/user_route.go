package route

import (
	"github.com/gofiber/fiber/v2"
	"lot/api/handler"
	"lot/pkg/service"
)

func SetupUserRoutes(router fiber.Router, userService service.UserService) {
	router.Post("/signup", handler.SignUpHandler(userService))
}
