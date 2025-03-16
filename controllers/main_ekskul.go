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
	tailwindColors := []string{
		"bg-emerald-600", "bg-blue-500", "bg-red-400", "bg-yellow-300",
		"bg-indigo-700", "bg-purple-500", "bg-pink-600", "bg-gray-400",
		"bg-green-500", "bg-teal-600", "bg-orange-500", "bg-cyan-400",
		"bg-lime-500", "bg-rose-600", "bg-violet-700", "bg-amber-500",
	}

	for i := 1; i <= 100; i++ {
		ekskul := models.MainEkskul{
			NamaEkskul: fmt.Sprintf("Ekskul Dummy %d", i),
			Warna:      tailwindColors[rand.Intn(len(tailwindColors))],
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
