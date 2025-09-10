// /Users/admin/Documents/golang/controllers/main_siswa.go
package controllers

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"sipandai/config"
	"sipandai/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
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
			Foto:         toPtr(fmt.Sprintf("foto_%d.jpg", rand.Intn(6)+1)),
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
	status := c.Params("status")

	page, _ := strconv.Atoi(c.Query("page", "1")) // Default page 1
	limit := 100                                  // Data per halaman
	offset := (page - 1) * limit                  // Hitung offset
	searchQuery := c.Query("search", "")          // Ambil query pencarian
	sortBy := c.Query("sort", "nama_siswa")       // Kolom yang diurutkan (default: nama_siswa)
	order := c.Query("order", "asc")              // Urutan (asc atau desc)

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
		Warna      string `json:"warna"`
	}

	var siswaList []Siswa
	var siswaEkskulList []SiswaEkskul
	var total int64

	// Query untuk menghitung total data
	countQuery := db.Table("main_siswa").
		Joins("LEFT JOIN sub_siswa_kelas ON sub_siswa_kelas.id_siswa = main_siswa.id_siswa").
		Joins("LEFT JOIN main_kelas ON main_kelas.id_kelas = sub_siswa_kelas.id_kelas").
		Joins("LEFT JOIN sub_siswa_jurusan ON sub_siswa_jurusan.id_siswa = main_siswa.id_siswa").
		Joins("LEFT JOIN main_jurusan ON main_jurusan.id_jurusan = sub_siswa_jurusan.id_jurusan").
		Joins("LEFT JOIN sub_siswa_rombel ON sub_siswa_rombel.id_siswa = main_siswa.id_siswa").
		Joins("LEFT JOIN main_rombel ON main_rombel.id_rombel = sub_siswa_rombel.id_rombel").
		Where("main_siswa.status = ?", status)

	if searchQuery != "" {
		countQuery = countQuery.Where("main_siswa.nama_siswa LIKE ? OR main_siswa.nis LIKE ? OR main_siswa.nisn LIKE ?",
			"%"+searchQuery+"%", "%"+searchQuery+"%", "%"+searchQuery+"%")
	}

	countQuery.Count(&total)

	// Query untuk mengambil data siswa dengan pagination, search, dan sorting
	dataQuery := db.Table("main_siswa").
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
		Where("main_siswa.status = ?", status).
		Order(fmt.Sprintf("%s %s", sortBy, order)). // Tambahkan sorting
		Limit(limit).Offset(offset)

	if searchQuery != "" {
		dataQuery = dataQuery.Where("main_siswa.nama_siswa LIKE ? OR main_siswa.nis LIKE ? OR main_siswa.nisn LIKE ?",
			"%"+searchQuery+"%", "%"+searchQuery+"%", "%"+searchQuery+"%")
	}

	err := dataQuery.Find(&siswaList).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Query untuk mendapatkan data ekskul siswa beserta warnanya
	err = db.Table("sub_siswa_ekskul").
		Select("sub_siswa_ekskul.id_siswa, main_ekskul.nama_ekskul, main_ekskul.warna").
		Joins("LEFT JOIN main_ekskul ON main_ekskul.id_ekskul = sub_siswa_ekskul.id_ekskul").
		Find(&siswaEkskulList).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Membuat map untuk menyimpan ekskul beserta warnanya berdasarkan id_siswa
	ekskulMap := make(map[int][]map[string]string)
	for _, ekskul := range siswaEkskulList {
		ekskulMap[ekskul.IDSiswa] = append(ekskulMap[ekskul.IDSiswa], map[string]string{
			"nama":  ekskul.NamaEkskul,
			"warna": ekskul.Warna,
		})
	}

	// Menambahkan data ekskul beserta warnanya ke dalam data siswa
	type SiswaResponse struct {
		Siswa
		Ekskul []map[string]string `json:"ekskul"`
	}

	var response []SiswaResponse
	for _, siswa := range siswaList {
		response = append(response, SiswaResponse{
			Siswa:  siswa,
			Ekskul: ekskulMap[siswa.IDSiswa],
		})
	}

	return c.JSON(fiber.Map{
		"data_siswa": response,
		"total":      total, // Total data untuk pagination
	})
}

// Get all student data for export
func GetAllSiswaForExport(c *fiber.Ctx) error {
	db := config.DB
	status := c.Params("status")

	type Ekskul struct {
		Nama  string `json:"nama"`
		Warna string `json:"warna"`
	}

	type ExportSiswa struct {
		IDSiswa      int      `json:"id_siswa"`
		KodeSiswa    string   `json:"kode_siswa"`
		NISN         string   `json:"nisn"`
		NIS          string   `json:"nis"`
		NamaSiswa    string   `json:"nama_siswa"`
		JenisKelamin string   `json:"jenis_kelamin"`
		TahunMasuk   int      `json:"tahun_masuk"`
		Foto         string   `json:"foto"`
		Status       int      `json:"status"`
		NamaKelas    string   `json:"nama_kelas"`
		NamaJurusan  string   `json:"nama_jurusan"`
		NamaRombel   string   `json:"nama_rombel"`
		Ekskul       []Ekskul `json:"ekskul"`
	}

	var siswaList []struct {
		IDSiswa      int
		KodeSiswa    string
		NISN         string
		NIS          string
		NamaSiswa    string
		JenisKelamin string
		TahunMasuk   int
		Foto         string
		Status       int
		NamaKelas    string
		NamaJurusan  string
		NamaRombel   string
	}

	// Ambil data siswa dengan data kelas, jurusan, dan rombel
	if err := db.Table("main_siswa").
		Select(`main_siswa.id_siswa, main_siswa.kode_siswa, main_siswa.nisn, main_siswa.nis,
                main_siswa.nama_siswa, main_siswa.jenis_kelamin, main_siswa.tahun_masuk,
                main_siswa.foto, main_siswa.status,
                COALESCE(main_kelas.nama_kelas, '') as nama_kelas,
                COALESCE(main_jurusan.nama_jurusan, '') as nama_jurusan,
                COALESCE(main_rombel.nama_rombel, '') as nama_rombel`).
		Joins("LEFT JOIN sub_siswa_kelas ON sub_siswa_kelas.id_siswa = main_siswa.id_siswa").
		Joins("LEFT JOIN main_kelas ON main_kelas.id_kelas = sub_siswa_kelas.id_kelas").
		Joins("LEFT JOIN sub_siswa_jurusan ON sub_siswa_jurusan.id_siswa = main_siswa.id_siswa").
		Joins("LEFT JOIN main_jurusan ON main_jurusan.id_jurusan = sub_siswa_jurusan.id_jurusan").
		Joins("LEFT JOIN sub_siswa_rombel ON sub_siswa_rombel.id_siswa = main_siswa.id_siswa").
		Joins("LEFT JOIN main_rombel ON main_rombel.id_rombel = sub_siswa_rombel.id_rombel").
		Where("main_siswa.status = ?", status).
		Find(&siswaList).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Ambil data ekskul
	var ekskulList []struct {
		IDSiswa int
		Nama    string
		Warna   string
	}

	if err := db.Table("sub_siswa_ekskul").
		Select("sub_siswa_ekskul.id_siswa, main_ekskul.nama_ekskul as nama, main_ekskul.warna").
		Joins("LEFT JOIN main_ekskul ON main_ekskul.id_ekskul = sub_siswa_ekskul.id_ekskul").
		Find(&ekskulList).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Kelompokkan ekskul berdasarkan siswa
	ekskulMap := make(map[int][]Ekskul)
	for _, ekskul := range ekskulList {
		ekskulMap[ekskul.IDSiswa] = append(ekskulMap[ekskul.IDSiswa], Ekskul{
			Nama:  ekskul.Nama,
			Warna: ekskul.Warna,
		})
	}

	// Gabungkan data siswa dan ekskul
	var response []ExportSiswa
	for _, siswa := range siswaList {
		response = append(response, ExportSiswa{
			IDSiswa:      siswa.IDSiswa,
			KodeSiswa:    siswa.KodeSiswa,
			NISN:         siswa.NISN,
			NIS:          siswa.NIS,
			NamaSiswa:    siswa.NamaSiswa,
			JenisKelamin: siswa.JenisKelamin,
			TahunMasuk:   siswa.TahunMasuk,
			Foto:         siswa.Foto,
			Status:       siswa.Status,
			NamaKelas:    siswa.NamaKelas,
			NamaJurusan:  siswa.NamaJurusan,
			NamaRombel:   siswa.NamaRombel,
			Ekskul:       ekskulMap[siswa.IDSiswa],
		})
	}

	return c.JSON(fiber.Map{
		"data_siswa": response,
		"total":      len(response),
		"success":    true,
	})
}

// Get all student IDs
func GetAllStudentIDs(c *fiber.Ctx) error {
	db := config.DB
	status := c.Params("status")

	var ids []int
	query := db.Table("main_siswa").Where("status = ?", status).Pluck("id_siswa", &ids)

	if err := query.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(fiber.Map{"data": []int{}})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"data": ids,
	})
}

func UpdateStatusSiswa(c *fiber.Ctx) error {
	db := config.DB

	// Ambil status dari parameter dan konversi ke int
	status := c.Params("status")
	statusInt, err := strconv.Atoi(status)
	if err != nil || (statusInt != 0 && statusInt != 1) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Parameter status harus berupa 0 atau 1.",
		})
	}

	// Hitung status baru yang akan diset (toggle)
	newStatus := 1 - statusInt

	// Ambil data dari body
	var requestBody struct {
		Ids []int `json:"ids"` // Ambil ID dari body
	}

	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Gagal mem-parsing body request.",
		})
	}

	ids := requestBody.Ids // Ambil ID dari request body

	if len(ids) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Minimal 1 id_siswa harus diberikan.",
		})
	}

	// Update status siswa
	if err := db.Model(&models.Siswa{}).
		Where("id_siswa IN ?", ids).
		Update("status", newStatus).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message":     "Berhasil mengubah status siswa.",
		"id_siswa":    ids,
		"status_baru": newStatus,
	})
}

// BulkDeleteSiswa menghapus siswa secara permanen
func BulkDeleteSiswa(c *fiber.Ctx) error {
	db := config.DB

	// Ambil data dari body
	var requestBody struct {
		Ids []int `json:"ids"` // Array ID siswa yang akan dihapus
	}

	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Gagal mem-parsing body request",
		})
	}

	ids := requestBody.Ids

	if len(ids) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Minimal 1 id_siswa harus diberikan",
		})
	}

	// Mulai transaksi
	tx := db.Begin()
	if tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal memulai transaksi",
		})
	}

	// 1. Hapus data dari tabel relasi terlebih dahulu

	// Hapus dari sub_siswa_kelas
	if err := tx.Where("id_siswa IN ?", ids).Delete(&models.SubSiswaKelas{}).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal menghapus data kelas siswa: " + err.Error(),
		})
	}

	// Hapus dari sub_siswa_jurusan
	if err := tx.Where("id_siswa IN ?", ids).Delete(&models.SubSiswaJurusan{}).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal menghapus data jurusan siswa: " + err.Error(),
		})
	}

	// Hapus dari sub_siswa_rombel
	if err := tx.Where("id_siswa IN ?", ids).Delete(&models.SubSiswaRombel{}).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal menghapus data rombel siswa: " + err.Error(),
		})
	}

	// Hapus dari sub_siswa_ekskul
	if err := tx.Where("id_siswa IN ?", ids).Delete(&models.SubSiswaEkskul{}).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal menghapus data ekskul siswa: " + err.Error(),
		})
	}

	// 2. Hapus dari tabel utama
	if err := tx.Where("id_siswa IN ?", ids).Delete(&models.Siswa{}).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal menghapus data siswa: " + err.Error(),
		})
	}

	// Commit transaksi jika semua berhasil
	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal melakukan commit transaksi: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message":  "Berhasil menghapus siswa secara permanen",
		"id_siswa": ids,
		"total":    len(ids),
	})
}

// controllers/main_siswa.go
func GetReferenceData(c *fiber.Ctx) error {
	db := config.DB

	type ReferenceData struct {
		Ekskul  []models.MainEkskul  `json:"ekskul"`
		Jurusan []models.MainJurusan `json:"jurusan"`
		Kelas   []models.MainKelas   `json:"kelas"`
		Rombel  []models.MainRombel  `json:"rombel"`
	}

	var result ReferenceData

	// Get all ekskul
	if err := db.Find(&result.Ekskul).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch ekskul data"})
	}

	// Get all jurusan
	if err := db.Find(&result.Jurusan).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch jurusan data"})
	}

	// Get all kelas
	if err := db.Find(&result.Kelas).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch kelas data"})
	}

	// Get all rombel
	if err := db.Find(&result.Rombel).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch rombel data"})
	}

	return c.JSON(result)
}

func GetSiswaRelations(c *fiber.Ctx) error {
	idSiswa := c.Params("id")

	type Response struct {
		IDKelas     int    `json:"id_kelas"`
		NamaKelas   string `json:"nama_kelas"`
		IDJurusan   int    `json:"id_jurusan"`
		NamaJurusan string `json:"nama_jurusan"`
		IDRombel    int    `json:"id_rombel"`
		NamaRombel  string `json:"nama_rombel"`
	}

	var response Response

	// Get kelas
	var subKelas models.SubSiswaKelas
	if err := config.DB.Where("id_siswa = ?", idSiswa).First(&subKelas).Error; err == nil {
		var kelas models.MainKelas
		if err := config.DB.Where("id_kelas = ?", subKelas.IDKelas).First(&kelas).Error; err == nil {
			response.IDKelas = subKelas.IDKelas
			response.NamaKelas = kelas.NamaKelas
		}
	}

	// Get jurusan
	var subJurusan models.SubSiswaJurusan
	if err := config.DB.Where("id_siswa = ?", idSiswa).First(&subJurusan).Error; err == nil {
		var jurusan models.MainJurusan
		if err := config.DB.Where("id_jurusan = ?", subJurusan.IDJurusan).First(&jurusan).Error; err == nil {
			response.IDJurusan = subJurusan.IDJurusan
			response.NamaJurusan = jurusan.NamaJurusan
		}
	}

	// Get rombel
	var subRombel models.SubSiswaRombel
	if err := config.DB.Where("id_siswa = ?", idSiswa).First(&subRombel).Error; err == nil {
		var rombel models.MainRombel
		if err := config.DB.Where("id_rombel = ?", subRombel.IDRombel).First(&rombel).Error; err == nil {
			response.IDRombel = subRombel.IDRombel
			response.NamaRombel = rombel.NamaRombel
		}
	}

	return c.JSON(response)
}

type UpdateSiswaRequest struct {
	NamaSiswa         string `json:"nama_siswa"`
	NIS               string `json:"nis"`
	NISN              string `json:"nisn"`
	IDKelas           int    `json:"id_kelas"`
	IDJurusan         int    `json:"id_jurusan"`
	IDRombel          int    `json:"id_rombel"`
	Status            string `json:"status"`
	IDEkskul          []int  `json:"id_ekskul"`
	PreviousRelations struct {
		IDKelas   int `json:"id_kelas"`
		IDJurusan int `json:"id_jurusan"`
		IDRombel  int `json:"id_rombel"`
	} `json:"previous_relations"`
}

func UpdateSiswa(c *fiber.Ctx) error {
	idSiswa, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "ID siswa tidak valid",
		})
	}

	var req UpdateSiswaRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   "Gagal parsing data",
			"details": err.Error(),
		})
	}

	// Mulai transaksi database
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. Update data utama siswa
	var siswa models.Siswa
	if err := tx.First(&siswa, idSiswa).Error; err != nil {
		tx.Rollback()
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Siswa tidak ditemukan",
		})
	}

	siswa.NamaSiswa = req.NamaSiswa
	siswa.Nis = req.NIS
	siswa.Nisn = req.NISN
	siswa.Status = 0
	if req.Status == "1" {
		siswa.Status = 1
	}

	siswa.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")

	if err := tx.Save(&siswa).Error; err != nil {
		tx.Rollback()
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal update data siswa",
		})
	}

	// 2. Update relasi kelas jika berubah
	if req.IDKelas != req.PreviousRelations.IDKelas {
		// Hapus relasi lama
		if err := tx.Where("id_siswa = ?", idSiswa).Delete(&models.SubSiswaKelas{}).Error; err != nil {
			tx.Rollback()
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Gagal menghapus relasi kelas lama",
			})
		}

		// Tambah relasi baru
		newKelas := models.SubSiswaKelas{
			IDSiswa: idSiswa,
			IDKelas: req.IDKelas,
		}
		if err := tx.Create(&newKelas).Error; err != nil {
			tx.Rollback()
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Gagal menambahkan relasi kelas baru",
			})
		}
	}

	// 3. Update relasi jurusan jika berubah
	if req.IDJurusan != req.PreviousRelations.IDJurusan {
		// Hapus relasi lama
		if err := tx.Where("id_siswa = ?", idSiswa).Delete(&models.SubSiswaJurusan{}).Error; err != nil {
			tx.Rollback()
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Gagal menghapus relasi jurusan lama",
			})
		}

		// Tambah relasi baru
		newJurusan := models.SubSiswaJurusan{
			IDSiswa:   idSiswa,
			IDJurusan: req.IDJurusan,
		}
		if err := tx.Create(&newJurusan).Error; err != nil {
			tx.Rollback()
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Gagal menambahkan relasi jurusan baru",
			})
		}
	}

	// 4. Update relasi rombel jika berubah
	if req.IDRombel != req.PreviousRelations.IDRombel {
		// Hapus relasi lama
		if err := tx.Where("id_siswa = ?", idSiswa).Delete(&models.SubSiswaRombel{}).Error; err != nil {
			tx.Rollback()
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Gagal menghapus relasi rombel lama",
			})
		}

		// Tambah relasi baru
		newRombel := models.SubSiswaRombel{
			IDSiswa:  idSiswa,
			IDRombel: req.IDRombel,
		}
		if err := tx.Create(&newRombel).Error; err != nil {
			tx.Rollback()
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Gagal menambahkan relasi rombel baru",
			})
		}
	}

	// 5. Update ekstrakurikuler
	// Hapus semua relasi ekskul lama
	if err := tx.Where("id_siswa = ?", idSiswa).Delete(&models.SubSiswaEkskul{}).Error; err != nil {
		tx.Rollback()
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal menghapus relasi ekskul lama",
		})
	}

	// Tambahkan relasi ekskul baru
	for _, ekskulID := range req.IDEkskul {
		newEkskul := models.SubSiswaEkskul{
			IDSiswa:      idSiswa,
			IDEkskul:     ekskulID,
			TanggalMasuk: time.Now(),
		}
		if err := tx.Create(&newEkskul).Error; err != nil {
			tx.Rollback()
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": fmt.Sprintf("Gagal menambahkan ekskul ID %d", ekskulID),
			})
		}
	}

	// Commit transaksi jika semua berhasil
	if err := tx.Commit().Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal commit transaksi",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Data siswa berhasil diperbarui",
		"data": fiber.Map{
			"id_siswa":   idSiswa,
			"nama_siswa": req.NamaSiswa,
			"nis":        req.NIS,
			"nisn":       req.NISN,
			"status":     req.Status,
			"kelas":      req.IDKelas,
			"jurusan":    req.IDJurusan,
			"rombel":     req.IDRombel,
			"ekskul":     req.IDEkskul,
		},
	})
}
