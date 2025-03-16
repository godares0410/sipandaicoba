package models

import (
	"time"
)

type SubSiswaKelas struct {
	IDSubSiswaKelas uint      `gorm:"primaryKey;autoIncrement" json:"id_sub_siswa_kelas"`
	IDSiswa         int       `json:"id_siswa"`
	IDKelas         int       `json:"id_kelas"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (SubSiswaKelas) TableName() string {
	return "sub_siswa_kelas"
}
