package models

import (
	"time"
)

type SubSiswaJurusan struct {
	IDSubSiswaJurusan uint      `gorm:"primaryKey;autoIncrement" json:"id_sub_siswa_jurusan"`
	IDSiswa           int       `json:"id_siswa"`
	IDJurusan         int       `json:"id_jurusan"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func (SubSiswaJurusan) TableName() string {
	return "sub_siswa_jurusan"
}
