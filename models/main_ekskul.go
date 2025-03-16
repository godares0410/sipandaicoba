package models

import (
	"time"
)

type MainEkskul struct {
	IDEkskul   uint      `gorm:"primaryKey;autoIncrement" json:"id_ekskul"`
	NamaEkskul string    `json:"nama_ekskul"`
	Warna      string    `json:"warna"`
	IDSekolah  int       `json:"id_sekolah"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (MainEkskul) TableName() string {
	return "main_ekskul"
}