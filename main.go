package main

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"sipandai/config"
	"sipandai/routes"
)

func main() {
	// Load Configuration
	config.LoadConfig()

	// Koneksi ke database
	config.ConnectDatabase()

	// Setup Fiber App
	app := fiber.New()

	// Setup Routes
	routes.SetupRoutes(app)

	// Start Server on port 3000
	port := ":" + config.AppConfig.AppPort
	log.Fatal(app.Listen(port))
}
