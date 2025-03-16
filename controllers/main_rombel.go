package controllers

import (
	"fmt"
	"time"
	"sipandai/config"
	"sipandai/models"

	"github.com/gofiber/fiber/v2"
)

func CreateDummyMainRombelData(c *fiber.Ctx) error {
	db := config.DB

	tx := db.Begin()
	if tx.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to start transaction: " + tx.Error.Error()})
	}

	now := time.Now()

	// Mock data for kelas and jurusan
	kelasNames := map[int]string{
		10: "X",
		11: "XI",
		12: "XII",
	}

	jurusanAbbr := map[int]string{
		1:  "TM",
		2:  "TO",
		3:  "TE",
		4:  "TKJ",
		5:  "RPL",
		6:  "MM",
		7:  "AN",
		8:  "DKV",
		9:  "TB",
		10: "TS",
		11: "AK",
		12: "AP",
		13: "PM",
		14: "PH",
		15: "TB",
		16: "TS",
		17: "KF",
		18: "AT",
		19: "AP",
		20: "TG",
	}

	for i := 1; i <= 20; i++ {
		idKelas := 10 + (i % 3) // 10-12
		idJurusan := i          // 1-20
		namaRombel := fmt.Sprintf("%s - %s", kelasNames[idKelas], jurusanAbbr[idJurusan])

		rombel := models.MainRombel{
			NamaRombel: namaRombel,
			IDKelas:    idKelas,
			IDJurusan:  idJurusan,
			IDSekolah:  1,
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		if err := tx.Create(&rombel).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("Failed to insert data: %s", err.Error())})
		}
	}

	if err := tx.Commit().Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to commit transaction: " + err.Error()})
	}

	return c.JSON(fiber.Map{"message": "20 dummy data untuk main rombel berhasil dibuat!"})
}
