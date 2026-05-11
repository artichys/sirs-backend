package models

// Pasien merepresentasikan tabel 'pasien' di database
type Pasien struct {
	IdPasien        int    `gorm:"primaryKey;column:id_pasien;autoIncrement" json:"id_pasien"`
	NomorRekamMedis string `gorm:"column:nomor_rekam_medis;unique" json:"nomor_rekam_medis"`
	NamaPasien      string `gorm:"column:nama_pasien" json:"nama_pasien"`
	TipePasien      string `gorm:"column:tipe_pasien;type:enum('Pasien Baru', 'Pasien Lama')" json:"tipe_pasien"`
	NoHp            string `gorm:"column:no_hp" json:"no_hp"`
}

func (Pasien) TableName() string {
	return "pasien"
}