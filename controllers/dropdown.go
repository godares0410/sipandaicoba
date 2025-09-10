// controllers/dropdown.go
package controllers

import (
	"sipandai/config"
	"sipandai/models"
	"github.com/gofiber/fiber/v2"
)

func GetEkskulOptions(c *fiber.Ctx) error {
	var ekskul []models.MainEkskul
	if err := config.DB.Find(&ekskul).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(ekskul)
}

func GetKelasOptions(c *fiber.Ctx) error {
	var kelas []models.MainKelas
	if err := config.DB.Find(&kelas).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(kelas)
}

func GetJurusanOptions(c *fiber.Ctx) error {
	var jurusan []models.MainJurusan
	if err := config.DB.Find(&jurusan).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(jurusan)
}

func GetRombelOptions(c *fiber.Ctx) error {
	var rombel []models.MainRombel
	if err := config.DB.Find(&rombel).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(rombel)
}

func GetSiswaEkskul(c *fiber.Ctx) error {
	idSiswa := c.Params("id")
	var siswaEkskul []models.SubSiswaEkskul
	
	if err := config.DB.Where("id_siswa = ?", idSiswa).Find(&siswaEkskul).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	
	return c.JSON(siswaEkskul)
}