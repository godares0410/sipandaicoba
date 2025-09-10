package controllers

import (
	"math/rand"
	"time"

	"sipandai/config"
	"sipandai/models"

	"github.com/gofiber/fiber/v2"
)

func CreateDummySubSiswaEkskulData(c *fiber.Ctx) error {
	db := config.DB

	tx := db.Begin()
	if tx.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to start transaction: " + tx.Error.Error()})
	}

	now := time.Now()
	rand.Seed(time.Now().UnixNano())

	for i := 1; i <= 10000; i++ {
		numEkskul := rand.Intn(3) + 1 // Setiap siswa memiliki 1-3 ekskul
		existingEkskul := map[int]bool{}

		for j := 0; j < numEkskul; j++ {
			var idEkskul int
			for {
				idEkskul = rand.Intn(100) + 1 // ID ekskul antara 1-100
				if !existingEkskul[idEkskul] {
					existingEkskul[idEkskul] = true
					break
				}
			}

			subSiswaEkskul := models.SubSiswaEkskul{
				IDSiswa:      i,
				IDEkskul:     idEkskul,
				TanggalMasuk: now.AddDate(0, -rand.Intn(12), -rand.Intn(30)),
				CreatedAt:    now,
				UpdatedAt:    now,
			}

			if err := tx.Create(&subSiswaEkskul).Error; err != nil {
				tx.Rollback()
				return c.Status(500).JSON(fiber.Map{"error": "Failed to insert data: " + err.Error()})
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to commit transaction: " + err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Dummy data untuk sub siswa ekskul berhasil dibuat!"})
}

// controllers/ekskul.go
func GetEkskulList(c *fiber.Ctx) error {
    var ekskulList []models.MainEkskul
    if err := config.DB.Find(&ekskulList).Error; err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    return c.JSON(fiber.Map{"data": ekskulList})
}

// controllers/siswa.go
func GetSiswaEkskulList(c *fiber.Ctx) error {
    idSiswa := c.Params("id")
    var ekskulList []models.SubSiswaEkskul
    
    if err := config.DB.Where("id_siswa = ?", idSiswa).Find(&ekskulList).Error; err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    
    return c.JSON(fiber.Map{"data": ekskulList})
}
