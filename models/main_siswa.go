// /Users/admin/Documents/golang/models/main_siswa.go
package models

type Siswa struct {
    IDSiswa     uint      `gorm:"primaryKey;column:id_siswa;autoIncrement" json:"id_siswa"`
    KodeSiswa   string    `gorm:"unique" json:"kode_siswa"`
    Nisn        string    `json:"nisn"`
    Nis         string    `json:"nis"`
    NamaSiswa   string    `json:"nama_siswa"`
    JenisKelamin string   `json:"jenis_kelamin"`
    TahunMasuk  *int64    `json:"tahun_masuk"`
    Foto        *string   `json:"foto"`
    Status      int       `json:"status"`
    IDSekolah   int       `json:"id_sekolah"`
    CreatedAt   string    `json:"created_at"`
    UpdatedAt   string    `json:"updated_at"`
}


func (Siswa) TableName() string {
	return "main_siswa"
}
