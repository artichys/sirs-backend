package models

type Ruangan struct {
	IdRuangan   int    `gorm:"primaryKey;column:id_ruangan;autoIncrement" json:"id_ruangan"`
	NamaRuangan string `gorm:"column:nama_ruangan" json:"nama_ruangan"`
}

// TableName mengeset nama tabel yang benar di database
func (Ruangan) TableName() string {
	return "ruangan"
}
