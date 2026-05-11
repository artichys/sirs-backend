package models

import "time"

// Booking merepresentasikan tabel 'booking' di database
type Booking struct {
	IdBooking     int       `gorm:"primaryKey;column:id_booking;autoIncrement" json:"id_booking"`
	IdJadwal      int       `gorm:"column:id_jadwal" json:"id_jadwal"`
	IdPasien      int       `gorm:"column:id_pasien" json:"id_pasien"`
	NomorAntrean  int       `gorm:"column:nomor_antrean" json:"nomor_antrean"`
	StatusBooking string    `gorm:"column:status_booking;type:enum('Menunggu', 'Diperiksa', 'Selesai');default:'Menunggu'" json:"status_booking"`
	WaktuBooking  time.Time `gorm:"column:waktu_booking;autoCreateTime" json:"waktu_booking"`
}

// Custom nama tabel agar sesuai dengan MySQL
func (Booking) TableName() string {
	return "booking"
}