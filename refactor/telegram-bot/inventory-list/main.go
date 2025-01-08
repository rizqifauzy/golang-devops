package main

import (
	"log"
	"telegram-inventory/clients"
	"telegram-inventory/config"
	"telegram-inventory/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
	}))
	bot := clients.Init()
	handlers.Init(bot)
	log.Fatal(app.Listen(":" + config.Config("PORT")))
}
