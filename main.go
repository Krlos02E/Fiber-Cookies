package main

import (
	controllers "SessionCookie/Controllers"
	database "SessionCookie/Database"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	LoadEnvVariables()
	db = database.NewConnection()
	database.MigrateData_User()
}

func main() {

	app := fiber.New()
	controllers.InitializeUserController(db, app)
	app.Listen(":" + os.Getenv("PORT"))
}

func LoadEnvVariables() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}
}
