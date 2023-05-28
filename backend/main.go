package main

import (
	"backend/databaze"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

var pocetSlov int = 10
var pocetPismenVeSlovu int = 4

func main() {
	databaze.DBConnect()

	app := fiber.New(fiber.Config{
		AppName: "pavouk",
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
		AllowHeaders: "Origin, Content-Type, Accept, Token",
	}))

	setupRouter(app)

	err := app.Listen(":8080")
	if err != nil {
		log.Fatal(err)
	} // http://localhost:8080
}
