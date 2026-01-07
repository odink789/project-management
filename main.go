package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/odink789/project-management/config"
	"github.com/odink789/project-management/controllers"
	"github.com/odink789/project-management/database/seed"
	"github.com/odink789/project-management/repositories"
	"github.com/odink789/project-management/routes"
	"github.com/odink789/project-management/services"
)

func main() {
	config.LoadEnv()
	config.ConnectDB()

	seed.SeedAdmin()
	//inisialisasi fiber

	app := fiber.New()

	userRepo := repositories.NewUserRepository()
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	routes.Setup(app, userController)

	port := config.AppConfig.AppPort
	log.Println("Server Is running On port :", port)
	log.Fatal(app.Listen(":" + port))

}
