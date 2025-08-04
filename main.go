package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/Ana-Gabs/actividadr-back/routes"
)

func main() {
	app := fiber.New()

	// Cargar rutas
	routes.SetupRoutes(app)

	// Levantar servidor
	app.Listen(":3000")
}