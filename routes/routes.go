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

	app.Get("/siswa/:status", controllers.GetSiswaData)
	app.Get("/idsiswa/:status", controllers.GetAllStudentIDs)
	app.Get("/siswa/all/:status", controllers.GetAllSiswaForExport)
	app.Put("/siswa/status/:status", controllers.UpdateStatusSiswa)
	app.Post("/dummy", controllers.CreateDummyData)
	app.Post("/dummyjurusan", controllers.CreateDummyMainJurusanData)
	app.Post("/dummyekskul", controllers.CreateDummyEkskulData)
	app.Post("/dummyrombel", controllers.CreateDummyMainRombelData)
	app.Post("/dummykelas", controllers.CreateDummyMainRombelData)
	app.Post("/dummyse", controllers.CreateDummySubSiswaEkskulData)
	app.Post("/dummysk", controllers.CreateDummySubSiswaKelasData)
	app.Post("/dummysj", controllers.CreateDummySubSiswaJurusanData)
	app.Post("/dummysr", controllers.CreateDummySubSiswaRombelData)
	app.Delete("/siswa", controllers.BulkDeleteSiswa)
	// routes.go
	app.Get("/ekskul", controllers.GetEkskulOptions)
	app.Get("/kelas", controllers.GetKelasOptions)
	app.Get("/jurusan", controllers.GetJurusanOptions)
	app.Get("/rombel", controllers.GetRombelOptions)
	app.Get("/siswa/:id/ekskul", controllers.GetSiswaEkskul)
	app.Get("/siswa/{id}/kelas", controllers.GetSiswaKelas)
	app.Get("/siswa/{id}/jurusan", controllers.GetSiswaJurusan)
	app.Get("/siswa/{id}/rombel", controllers.GetSiswaRombel)
	app.Get("/siswa/:id/relations", controllers.GetSiswaRelations)
	app.Put("/siswa/:id", controllers.UpdateSiswa)

	// Test route lain
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "pong",
		})
	})
}
