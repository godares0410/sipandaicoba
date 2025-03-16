package models

import (
	"time"
)

type MainKelas struct {
	IDKelas   uint      `gorm:"primaryKey;autoIncrement" json:"id_kelas"`
	NamaKelas string    `json:"nama_kelas"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (MainKelas) TableName() string {
	return "main_kelas"
}
