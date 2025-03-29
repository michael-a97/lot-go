package main

import (
	"context"
	"log"
	"lot/api/route"
	"lot/config"
	authRepository "lot/pkg/repository/auth"
	roleRepository "lot/pkg/repository/role"
	userRepository "lot/pkg/repository/user"
	authService "lot/pkg/service/auth"
	smsTokenVerifierService "lot/pkg/service/sms_token_verifier"
	userService "lot/pkg/service/user"

	"lot/platform/database"
	firebaseApp "lot/platform/firebase"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	config.LoadEnv()
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

	userRepository := userRepository.NewUserRepository(db)
	roleRepository := roleRepository.NewRoleRepository(db)
	authRepository := authRepository.NewAuthRepository(db)

	userService := userService.NewUserService(userRepository, roleRepository)
	authService := authService.NewAuthService(
		authRepository, userRepository,
		smsTokenVerifierService.NewFirebaseSmsTokenVerifier(firebaseAuthClient),
	)

	route.SetupUserRoutes(api.Group("/user"), userService, authService)
	route.SetupAuthRoutes(api.Group("/auth"), authService)

	appPort, err := config.Config("appPort")
	if err != nil {
		log.Fatal("please specify an `appPort`")
	}
	log.Fatal(app.Listen(":" + appPort))
}
