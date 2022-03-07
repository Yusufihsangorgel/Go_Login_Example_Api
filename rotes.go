package main

import (
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	app.Post("/register", Register)
	app.Post("/login", Login)
	app.Get("/user", User)
	app.Get("/logout", Logout)

}
