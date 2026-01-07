package seed

import (
	"log"

	"github.com/odink789/project-management/config"
	"github.com/odink789/project-management/models"
	"github.com/odink789/project-management/utils"
)




func SeedAdmin() {
	password,_ := utils.password("admin123")

	admin := models.User{
		Name: "Super admin",
		Email: "admin@example",
		Password: password,
		Role: "admin",
	}
	if err := config.DB.FirstOrCreate(&admin, models.User(Email: admin.Email)){
		log.Fatal("Failed to seed admin user:", err)
	} else {
		log.Println("Admin user seeded successfully")
	}

}