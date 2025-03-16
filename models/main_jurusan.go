package models

import (
	"time"
)

type MainJurusan struct {
	IDJurusan   uint      `gorm:"primaryKey;autoIncrement" json:"id_jurusan"`
	NamaJurusan string    `json:"nama_jurusan"`
	IDSekolah   int       `json:"id_sekolah"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (MainJurusan) TableName() string {
	return "main_jurusan"
}
