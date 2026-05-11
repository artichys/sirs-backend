package models

type Dokter struct {
	IdDokter   int    `gorm:"primaryKey;column:id_dokter;autoIncrement" json:"id_dokter"`
	NamaDokter string `gorm:"column:nama_dokter" json:"nama_dokter"`
}

// TableName mengeset nama tabel yang benar di database
func (Dokter) TableName() string {
	return "dokter"
}