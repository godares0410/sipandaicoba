package controllers

import (
	"errors"
	"fmt"
	"math/rand"
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
		Where("main_siswa.status = 1")

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
        Where("main_siswa.status = 1").
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
        Where("main_siswa.status = 1").
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

	var ids []int
	query := db.Table("main_siswa").Where("status = 1").Pluck("id_siswa", &ids)

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

func UpdateStatusSiswaToInactive(c *fiber.Ctx) error {
	db := config.DB

	// Ambil semua query id (meskipun hanya 1 id tetap akan masuk slice)
	queryArgs := c.Request().URI().QueryArgs()
	var ids []int

	queryArgs.VisitAll(func(key, value []byte) {
		if string(key) == "id" {
			id, err := strconv.Atoi(string(value))
			if err == nil {
				ids = append(ids, id)
			}
		}
	})

	// Cek apakah ada id yang valid
	if len(ids) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Minimal 1 id_siswa harus diberikan.",
		})
	}

	// Update semua status siswa menjadi 0
	if err := db.Model(&models.Siswa{}).
		Where("id_siswa IN ?", ids).
		Update("status", 0).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message":  "Berhasil mengubah status siswa menjadi nonaktif.",
		"id_siswa": ids,
	})
}
