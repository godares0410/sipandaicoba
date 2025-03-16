package models

import (
	"time"
)

type SubSiswaRombel struct {
	IDSubSiswaRombel uint      `gorm:"primaryKey;autoIncrement" json:"id_sub_siswa_rombel"`
	IDSiswa          int       `json:"id_siswa"`
	IDRombel         int       `json:"id_rombel"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (SubSiswaRombel) TableName() string {
	return "sub_siswa_rombel"
}