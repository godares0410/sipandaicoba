package models

import (
	"time"
)

type MainRombel struct {
	IDRombel   uint      `gorm:"primaryKey;autoIncrement" json:"id_rombel"`
	NamaRombel string    `json:"nama_rombel"`
	IDKelas    int       `json:"id_kelas"`
	IDJurusan  int       `json:"id_jurusan"`
	IDSekolah  int       `json:"id_sekolah"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (MainRombel) TableName() string {
	return "main_rombel"
}