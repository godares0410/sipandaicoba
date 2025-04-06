package routes

import (
	"github.com/gofiber/fiber/v2"
	"sipandai/controllers"
)

func SetupRoutes(app *fiber.App) {
	// Test route untuk memastikan server berjalan
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Server is running on port 3000 🚀")
	})

	app.Get("/siswa", controllers.GetSiswaData)
	app.Get("/idsiswa", controllers.GetAllStudentIDs)
	app.Post("/dummy", controllers.CreateDummyData)
	app.Post("/dummyjurusan", controllers.CreateDummyMainJurusanData)
	app.Post("/dummyekskul", controllers.CreateDummyEkskulData)
	app.Post("/dummyrombel", controllers.CreateDummyMainRombelData)
	app.Post("/dummykelas", controllers.CreateDummyMainRombelData)
	app.Post("/dummyse", controllers.CreateDummySubSiswaEkskulData)
	app.Post("/dummysk", controllers.CreateDummySubSiswaKelasData)
	app.Post("/dummysj", controllers.CreateDummySubSiswaJurusanData)
	app.Post("/dummysr", controllers.CreateDummySubSiswaRombelData)

	// Test route lain
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "pong",
		})
	})
}
