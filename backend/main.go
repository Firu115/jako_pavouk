package main

import (
	"backend/databaze"
	"log"

	"github.com/gofiber/fiber/v2"
)

var pocetSlov int = 10
var pocetPismenVeSlovu int = 4

func main() {
	databaze.DBConnect()

	app := fiber.New(fiber.Config{
		AppName:                 "pavouk",
		EnableTrustedProxyCheck: true,
		TrustedProxies:          []string{"localhost"},
	})
	setupRouter(app)

	err := app.Listen(":8080")
	if err != nil {
		log.Fatal(err)
	} // http://localhost:8080
}
