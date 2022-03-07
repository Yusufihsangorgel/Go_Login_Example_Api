package main

import (
	"backendtest/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

//hocam anydeskten dolayı ctrl + s yapamuyorum siz tüm ddosyalrı akydederbilr msiiniz tamamdır hocam
func main() {

	//database bağlandık
	database.Connect()

	//fiber app oluşturduk
	app := fiber.New(fiber.Config{
		AppName: "Test app",
	})

	//cors middleware ekledik
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	//routelarımızı ekledik
	Setup(app)

	//app'i başlatıyoruz
	app.Listen("localhost:3000")

}
