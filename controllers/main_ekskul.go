package controllers

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gofiber/fiber/v2"
	"sipandai/config"
	"sipandai/models"
)

func CreateDummyEkskulData(c *fiber.Ctx) error {
	db := config.DB

	tx := db.Begin()
	if tx.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to start transaction: " + tx.Error.Error()})
	}

	now := time.Now()
	hexColors := []string{
		"#10B981", "#3B82F6", "#EF4444", "#FACC15",
		"#4F46E5", "#8B5CF6", "#EC4899", "#9CA3AF",
		"#22C55E", "#14B8A6", "#F97316", "#06B6D4",
		"#A3E635", "#E11D48", "#6D28D9", "#F59E0B",
	}

	for i := 1; i <= 100; i++ {
		ekskul := models.MainEkskul{
			NamaEkskul: fmt.Sprintf("Ekskul Dummy %d", i),
			Warna:      hexColors[rand.Intn(len(hexColors))],
			IDSekolah:  1,
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		if err := tx.Create(&ekskul).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"error": "Failed to insert data: " + err.Error()})
		}
	}

	if err := tx.Commit().Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to commit transaction: " + err.Error()})
	}

	return c.JSON(fiber.Map{"message": "100 dummy ekskul berhasil dibuat!"})
}
