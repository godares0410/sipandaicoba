package controllers

import (
	"math/rand"
	"time"

	"github.com/gofiber/fiber/v2"
	"sipandai/config"
	"sipandai/models"
)

func CreateDummySubSiswaKelasData(c *fiber.Ctx) error {
	db := config.DB

	tx := db.Begin()
	if tx.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to start transaction: " + tx.Error.Error()})
	}

	now := time.Now()
	rand.Seed(time.Now().UnixNano())

	for i := 1; i <= 10000; i++ {
		subSiswaKelas := models.SubSiswaKelas{
			IDSiswa:   i,
			IDKelas:   rand.Intn(3) + 10, // ID kelas antara 10-12
			CreatedAt: now,
			UpdatedAt: now,
		}

		if err := tx.Create(&subSiswaKelas).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"error": "Failed to insert data: " + err.Error()})
		}
	}

	if err := tx.Commit().Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to commit transaction: " + err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Dummy data untuk sub siswa kelas berhasil dibuat!"})
}
