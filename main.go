package main

import (
	"context"
	"log"
	"lot/api/route"
	"lot/config"
	"lot/pkg/repository"
	"lot/pkg/service"
	"lot/platform/database"
	firebaseApp "lot/platform/firebase"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	db := database.ConnectDb()

	firebaseApp, err := firebaseApp.ConnectFirebaseApp()
	if err != nil {
		log.Fatal("coun't connect firebase " + err.Error())
	}

	firebaseAuthClient, err := firebaseApp.Auth(context.Background())
	if err != nil {
		log.Fatal("coun't connect firebase " + err.Error())
	}

	app := fiber.New(fiber.Config{AppName: "lot"})
	app.Use(logger.New())
	api := app.Group("/api/v1")

	userRepository := repository.NewUserRepository(db)
	roleRepository := repository.NewRoleRepository(db)
	authRepository := repository.NewAuthRepository(db)

	userService := service.NewUserService(userRepository, roleRepository)
	authService := service.NewAuthService(
		authRepository, userRepository,
		service.NewFirebaseSmsTokenVerifier(firebaseAuthClient),
	)

	route.SetupUserRoutes(api.Group("/user"), userService, authService)
	route.SetupAuthRoutes(api.Group("/auth"), authService)

	log.Fatal(app.Listen(":" + config.Config("appPort")))
}
