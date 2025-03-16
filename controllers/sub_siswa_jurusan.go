package controllers

import (
	"fmt"
	"math/rand"
	"sipandai/config"
	"sipandai/models"
	"time"

	"github.com/gofiber/fiber/v2"
)

func CreateDummySubSiswaJurusanData(c *fiber.Ctx) error {
	db := config.DB
	tx := db.Begin()
	if tx.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to start transaction: " + tx.Error.Error()})
	}

	now := time.Now()
	rand.Seed(now.UnixNano())

	for i := 1; i <= 10000; i++ {
		idJurusan := rand.Intn(20) + 1
		subSiswaJurusan := models.SubSiswaJurusan{
			IDSiswa:   i,
			IDJurusan: idJurusan,
			CreatedAt: now,
			UpdatedAt: now,
		}

		if err := tx.Create(&subSiswaJurusan).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("Failed to insert data for siswa %d: %s", i, err.Error())})
		}
	}

	if err := tx.Commit().Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to commit transaction: " + err.Error()})
	}

	return c.JSON(fiber.Map{"message": "10000 dummy data untuk sub siswa jurusan berhasil dibuat!"})
}
