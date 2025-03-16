package models

import (
	"time"
)

type SubSiswaEkskul struct {
	IDSubSiswaEkskul uint      `gorm:"primaryKey;autoIncrement" json:"id_sub_siswa_ekskul"`
	IDSiswa          int       `json:"id_siswa"`
	IDEkskul         int       `json:"id_ekskul"`
	TanggalMasuk     time.Time `gorm:"type:date" json:"tanggal_masuk"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (SubSiswaEkskul) TableName() string {
	return "sub_siswa_ekskul"
}
