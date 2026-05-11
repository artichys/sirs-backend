package controllers

import (
	"net/http"
	"sirs-backend/config"

	"github.com/gin-gonic/gin"
)

// Struct ditambahkan IdBooking
type MonitorResponse struct {
	IdBooking       int    `json:"id_booking"`
	NomorAntrean    int    `json:"nomor_antrean"`
	NomorRekamMedis string `json:"nomor_rekam_medis"`
	NamaPasien      string `json:"nama_pasien"`
	StatusBooking   string `json:"status_booking"`
}

func GetAntreanMonitor(c *gin.Context) {
	idJadwal := c.Param("id_jadwal")
	var antrean []MonitorResponse

	// Pastikan booking.id_booking ada di dalam fungsi Select() di bawah ini
	err := config.DB.Table("booking").
		Select("booking.id_booking, booking.nomor_antrean, pasien.nomor_rekam_medis, pasien.nama_pasien, booking.status_booking").
		Joins("JOIN pasien ON booking.id_pasien = pasien.id_pasien").
		Where("booking.id_jadwal = ? AND booking.status_booking != 'Selesai'", idJadwal).
		Order("booking.nomor_antrean ASC").
		Scan(&antrean).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memuat data layar monitor"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Data monitor berhasil dimuat",
		"id_jadwal": idJadwal,
		"data":      antrean,
	})
}
