package main

import (
	"log"
	"sipandai/config"
	"sipandai/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// Load Configuration
	config.LoadConfig()

	// Koneksi ke database
	config.ConnectDatabase()

	// Setup Fiber App
	app := fiber.New()

	// Setup CORS Middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",          // Izinkan permintaan dari frontend Next.js
		AllowHeaders: "Origin, Content-Type, Accept",   // Header yang diizinkan
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH", // Metode HTTP yang diizinkan
	}))

	// Setup Routes
	routes.SetupRoutes(app)

	// Start Server on port 3000
	port := ":" + config.AppConfig.AppPort
	log.Fatal(app.Listen(port))
}
