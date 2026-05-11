package controllers

import (
	"net/http"
	"sirs-backend/config"
	"github.com/gin-gonic/gin"
)

// Struct untuk menampung hasil query raw (rekap)
type RekapRuangan struct {
	IdRuangan    int    `json:"id_ruangan"`
	NamaRuangan  string `json:"nama_ruangan"`
	TotalJadwal  int    `json:"total_jadwal"`
	TotalPasien  int    `json:"total_pasien"`
}

// GetLaporanUtilisasi: API untuk direksi melihat produktivitas rumah sakit
func GetLaporanUtilisasi(c *gin.Context) {
	// Ambil parameter bulan dan tahun dari URL (opsional)
	// Contoh: /api/laporan/utilisasi?bulan=05&tahun=2026
	bulan := c.DefaultQuery("bulan", "05") // Default contoh: Mei
	tahun := c.DefaultQuery("tahun", "2026") // Default contoh: 2026

	var rekapData []RekapRuangan

	// Query SQL Agregasi: Menghitung total jadwal dan total pasien per ruangan
	// Query ini mengelompokkan (GROUP BY) berdasarkan ruangan
	query := `
		SELECT 
			r.id_ruangan, 
			r.nama_ruangan, 
			COUNT(DISTINCT j.id_jadwal) as total_jadwal,
			COUNT(b.id_booking) as total_pasien
		FROM ruangan r
		LEFT JOIN jadwal j ON r.id_ruangan = j.id_ruangan 
			AND MONTH(j.tanggal) = ? 
			AND YEAR(j.tanggal) = ?
			AND j.status_jadwal = 'Confirmed'
		LEFT JOIN booking b ON j.id_jadwal = b.id_jadwal
		GROUP BY r.id_ruangan, r.nama_ruangan
		ORDER BY total_pasien DESC
	`

	// Eksekusi Raw SQL menggunakan GORM
	err := config.DB.Raw(query, bulan, tahun).Scan(&rekapData).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghasilkan laporan utilisasi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Laporan utilisasi ruangan berhasil dimuat",
		"periode": tahun + "-" + bulan,
		"data":    rekapData,
	})
}