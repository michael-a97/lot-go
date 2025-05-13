package middleware

import (
	"lot/api/dto"
	"lot/config"
	"lot/internal/service"
	"strings"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func ProtectedMiddleware() fiber.Handler {
	secret, _ := config.Config("secret")
	return jwtware.New(
		jwtware.Config{
			SigningKey: jwtware.SigningKey{
				Key: []byte(secret),
			},
			ErrorHandler: jwtError,
		},
	)
}
func AuthenticationMiddleware(authService service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := strings.Split(c.Get("Authorization"), " ")[1]
		user, err := authService.GetUserFromAccessToken(token)
		if err != nil {
			return c.JSON(
				dto.ApiResponse{
					Error:  err,
					Status: fiber.StatusUnauthorized,
				},
			)
		}

		c.Locals("user", user)

		return c.Next()
	}
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Missing or malformed JWT",
				"data":    nil,
			},
		)

	}
	return c.Status(fiber.StatusUnauthorized).JSON(
		fiber.Map{
			"status":  "error",
			"message": "Invalid or expired JWT",
			"data":    nil,
		},
	)
}
