package controllers

import (
	"fmt"
	"math/rand"
	"sipandai/config"
	"sipandai/models"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GetSiswaRombel(c *fiber.Ctx) error {
    idSiswa := c.Params("id")
    var siswaRombel models.SubSiswaRombel
    
    if err := config.DB.Where("id_siswa = ?", idSiswa).First(&siswaRombel).Error; err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "Data tidak ditemukan"})
    }
    
    return c.JSON(siswaRombel)
}

func CreateDummySubSiswaRombelData(c *fiber.Ctx) error {
	db := config.DB
	tx := db.Begin()
	if tx.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to start transaction: " + tx.Error.Error()})
	}

	now := time.Now()
	rand.Seed(now.UnixNano())

	for i := 1; i <= 10000; i++ {
		idRombel := rand.Intn(20) + 1
		subSiswaRombel := models.SubSiswaRombel{
			IDSiswa:   i,
			IDRombel:  idRombel,
			CreatedAt: now,
			UpdatedAt: now,
		}

		if err := tx.Create(&subSiswaRombel).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("Failed to insert data for siswa %d: %s", i, err.Error())})
		}
	}

	if err := tx.Commit().Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to commit transaction: " + err.Error()})
	}

	return c.JSON(fiber.Map{"message": "10000 dummy data untuk sub siswa rombel berhasil dibuat!"})
}
