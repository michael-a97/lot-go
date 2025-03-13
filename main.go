package main

import (
	"log"
	"lot/api/route"
	"lot/config"
	"lot/pkg/repository"
	"lot/pkg/service"
	"lot/platform/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	db := database.ConnectDb()

	app := fiber.New(fiber.Config{AppName: "lot"})
	app.Use(logger.New())
	api := app.Group("/api/v1")

	userRepository := repository.NewUserRepository(db)
	roleRepository := repository.NewRoleRepository(db)
	authRepository := repository.NewAuthRepository(db)
	userService := service.NewUserService(userRepository, roleRepository)
	authService := service.NewAuthService(authRepository, userRepository)
	route.SetupUserRoutes(api.Group("/user"), userService)
	route.SetupAuthRoutes(api.Group("/auth"), authService)

	log.Fatal(app.Listen(":" + config.Config("appPort")))
}
