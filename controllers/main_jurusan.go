package controllers

import (
	"fmt"
	"time"
	"github.com/gofiber/fiber/v2"
	"sipandai/config"
	"sipandai/models"
)

func CreateDummyMainJurusanData(c *fiber.Ctx) error {
	db := config.DB

	tx := db.Begin()
	if tx.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to start transaction: " + tx.Error.Error()})
	}

	now := time.Now()

	jurusanNames := []string{
		"Teknik Mesin",
		"Teknik Otomotif",
		"Teknik Elektronika",
		"Teknik Komputer dan Jaringan",
		"Rekayasa Perangkat Lunak",
		"Multimedia",
		"Animasi",
		"Desain Komunikasi Visual",
		"Teknik Bangunan",
		"Teknik Sipil",
		"Akuntansi",
		"Administrasi Perkantoran",
		"Pemasaran",
		"Perhotelan",
		"Tata Boga",
		"Tata Busana",
		"Kesehatan dan Farmasi",
		"Agribisnis Tanaman",
		"Agribisnis Peternakan",
		"Teknik Geologi",
	}

	for _, name := range jurusanNames {
		jurusan := models.MainJurusan{
			NamaJurusan: name,
			IDSekolah:   1,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		if err := tx.Create(&jurusan).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("Failed to insert data: %s", err.Error())})
		}
	}

	if err := tx.Commit().Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to commit transaction: " + err.Error()})
	}

	return c.JSON(fiber.Map{"message": "20 dummy data untuk main jurusan berhasil dibuat!"})
}
