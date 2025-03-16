package controllers

import (
	"fmt"
	"math/rand"
	"sipandai/config"
	"sipandai/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func generateRandomNumber(length int) string {
	min := int64(1)
	for i := 1; i < length; i++ {
		min *= 10
	}
	max := min*10 - 1
	return strconv.FormatInt(rand.Int63n(max-min)+min, 10)
}

func CreateDummyData(c *fiber.Ctx) error {
	db := config.DB

	tx := db.Begin()
	if tx.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to start transaction: " + tx.Error.Error()})
	}

	toPtr := func(s string) *string {
		return &s
	}

	now := time.Now().Format("2006-01-02 15:04:05")

	for i := 1; i <= 10000; i++ {
		tahun := int64(rand.Intn(5) + 2018)

		siswa := models.Siswa{
			KodeSiswa:    fmt.Sprintf("SIS%04d", i),
			Nisn:         generateRandomNumber(10),
			Nis:          generateRandomNumber(10),
			NamaSiswa:    fmt.Sprintf("Siswa Dummy %d", i),
			JenisKelamin: []string{"L", "P"}[rand.Intn(2)],
			TahunMasuk:   &tahun,
			Foto:         toPtr(fmt.Sprintf("foto_%d.jpg", i)),
			Status:       1,
			IDSekolah:    1,
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		
		if err := tx.Create(&siswa).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"error": "Failed to insert data: " + err.Error()})
		}
	}

	if err := tx.Commit().Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to commit transaction: " + err.Error()})
	}

	return c.JSON(fiber.Map{"message": "10000 dummy siswa berhasil dibuat!"})
}

func GetSiswaData(c *fiber.Ctx) error {
	db := config.DB

	type Siswa struct {
		IDSiswa      int       `json:"id_siswa"`
		KodeSiswa    string    `json:"kode_siswa"`
		NISN         string    `json:"nisn"`
		NIS          string    `json:"nis"`
		NamaSiswa    string    `json:"nama_siswa"`
		JenisKelamin string    `json:"jenis_kelamin"`
		TahunMasuk   int       `json:"tahun_masuk"`
		Foto         string    `json:"foto"`
		Status       string    `json:"status"`
		IDSekolah    int       `json:"id_sekolah"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		NamaKelas    string    `json:"nama_kelas"`
		NamaJurusan  string    `json:"nama_jurusan"`
		NamaRombel   string    `json:"nama_rombel"`
	}

	type SiswaEkskul struct {
		IDSiswa    int    `json:"id_siswa"`
		NamaEkskul string `json:"nama_ekskul"`
	}

	var siswaList []Siswa
	var siswaEkskulList []SiswaEkskul

	// Query untuk mendapatkan data siswa
	err := db.Table("main_siswa").
		Select(`main_siswa.id_siswa, main_siswa.kode_siswa, main_siswa.nisn, main_siswa.nis, 
	        main_siswa.nama_siswa, main_siswa.jenis_kelamin, main_siswa.tahun_masuk, 
	        main_siswa.foto, main_siswa.status, main_siswa.id_sekolah, 
	        main_siswa.created_at, main_siswa.updated_at, main_kelas.nama_kelas, main_jurusan.nama_jurusan, main_rombel.nama_rombel`).
		Joins("LEFT JOIN sub_siswa_kelas ON sub_siswa_kelas.id_siswa = main_siswa.id_siswa").
		Joins("LEFT JOIN main_kelas ON main_kelas.id_kelas = sub_siswa_kelas.id_kelas").
		Joins("LEFT JOIN sub_siswa_jurusan ON sub_siswa_jurusan.id_siswa = main_siswa.id_siswa").
		Joins("LEFT JOIN main_jurusan ON main_jurusan.id_jurusan = sub_siswa_jurusan.id_jurusan").
		Joins("LEFT JOIN sub_siswa_rombel ON sub_siswa_rombel.id_siswa = main_siswa.id_siswa").
		Joins("LEFT JOIN main_rombel ON main_rombel.id_rombel = sub_siswa_rombel.id_rombel").
		Find(&siswaList).Error // Gunakan Find() untuk memetakan hasil query ke slice of struct

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Query untuk mendapatkan data ekskul siswa
	err = db.Table("sub_siswa_ekskul").
		Select("sub_siswa_ekskul.id_siswa, main_ekskul.nama_ekskul").
		Joins("LEFT JOIN main_ekskul ON main_ekskul.id_ekskul = sub_siswa_ekskul.id_ekskul").
		Find(&siswaEkskulList).Error // Gunakan Find() untuk memetakan hasil query ke slice of struct

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Membuat map untuk menyimpan ekskul berdasarkan id_siswa
	ekskulMap := make(map[int][]string)
	for _, ekskul := range siswaEkskulList {
		ekskulMap[ekskul.IDSiswa] = append(ekskulMap[ekskul.IDSiswa], ekskul.NamaEkskul)
	}

	// Menambahkan data ekskul ke dalam data siswa
	type SiswaResponse struct {
		Siswa
		NamaEkskul []string `json:"nama_ekskul"`
	}

	var response []SiswaResponse
	for _, siswa := range siswaList {
		response = append(response, SiswaResponse{
			Siswa:      siswa,
			NamaEkskul: ekskulMap[siswa.IDSiswa],
		})
	}

	return c.JSON(fiber.Map{"data_siswa": response})
}