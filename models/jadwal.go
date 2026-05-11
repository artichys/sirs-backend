package models

type Jadwal struct {
	IdJadwal     int    `gorm:"primaryKey;column:id_jadwal;autoIncrement" json:"id_jadwal"`
	IdDokter     int    `gorm:"column:id_dokter" json:"id_dokter"`
	IdRuangan    int    `gorm:"column:id_ruangan" json:"id_ruangan"`
	Tanggal      string `gorm:"column:tanggal" json:"tanggal"`
	JamMulai     string `gorm:"column:jam_mulai" json:"jam_mulai"`
	JamSelesai   string `gorm:"column:jam_selesai" json:"jam_selesai"`
	StatusJadwal string `gorm:"column:status_jadwal" json:"status_jadwal"`
	Catatan      string `gorm:"column:catatan" json:"catatan"`

	// Relasi yang dibutuhkan oleh Preload()
	Ruangan Ruangan `gorm:"foreignKey:IdRuangan;references:IdRuangan" json:"ruangan"`
	Dokter  Dokter  `gorm:"foreignKey:IdDokter;references:IdDokter" json:"dokter"` // <-- TAMBAHKAN BARIS INI
}

func (Jadwal) TableName() string {
	return "jadwal"
}
