package controllers

import (
	"net/http"
	"sirs-backend/config"
	"sirs-backend/models" // Asumsi Anda sudah membuat struct modelnya

	"github.com/gin-gonic/gin"
)

func ApproveJadwal(c *gin.Context) {
	jadwalID := c.Param("id")
	var jadwal models.Jadwal

	// 1. Cari jadwal yang diajukan
	if err := config.DB.First(&jadwal, jadwalID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Jadwal tidak ditemukan"})
		return
	}

	// 2. Cek Bentrok (Business Logic)
	// Pastikan tidak ada jadwal 'Confirmed' di ruangan yang sama, kecuali jadwal ini sendiri
	var konflik int64
	config.DB.Model(&models.Jadwal{}).Where(
		"id_ruangan = ? AND tanggal = ? AND status_jadwal = 'Confirmed' AND id_jadwal != ? AND ((jam_mulai < ? AND jam_selesai > ?) OR (jam_mulai < ? AND jam_selesai > ?))",
		jadwal.IdRuangan, jadwal.Tanggal, jadwal.IdJadwal, jadwal.JamSelesai, jadwal.JamMulai, jadwal.JamSelesai, jadwal.JamMulai,
	).Count(&konflik)

	// Tambahan: Cek apakah jadwal sebenarnya sudah di-approve sebelumnya
	if jadwal.StatusJadwal == "Confirmed" {
		c.JSON(http.StatusOK, gin.H{"message": "Jadwal ini sudah dikonfirmasi sebelumnya."})
		return
	}

	if konflik > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Jadwal ditolak! Terdapat konflik ruangan pada jam tersebut."})
		return
	}

	// 3. Update status jika aman
	jadwal.StatusJadwal = "Confirmed"
	config.DB.Save(&jadwal)

	c.JSON(http.StatusOK, gin.H{"message": "Jadwal berhasil divalidasi dan dikonfirmasi"})
}

type JadwalInput struct {
	IdDokter   int    `json:"id_dokter"`
	IdRuangan  int    `json:"id_ruangan"`
	Tanggal    string `json:"tanggal"`
	JamMulai   string `json:"jam_mulai"`
	JamSelesai string `json:"jam_selesai"`
}

// CreateJadwal: Dokter mengajukan jadwal baru
func CreateJadwal(c *gin.Context) {
	var input JadwalInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jadwalBaru := models.Jadwal{
		IdDokter:     input.IdDokter,
		IdRuangan:    input.IdRuangan,
		Tanggal:      input.Tanggal,
		JamMulai:     input.JamMulai,
		JamSelesai:   input.JamSelesai,
		StatusJadwal: "Draft", // Otomatis Draft saat pertama kali diajukan
	}

	if err := config.DB.Create(&jadwalBaru).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengajukan jadwal"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Jadwal berhasil diajukan dan menunggu validasi Admin",
		"data":    jadwalBaru,
	})
}

func GetSemuaJadwal(c *gin.Context) {
	var daftarJadwal []models.Jadwal
	// Ambil semua jadwal, urutkan dari yang terbaru
	if err := config.DB.Order("id_jadwal DESC").Find(&daftarJadwal).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data jadwal"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": daftarJadwal})
}

// Fungsi untuk mengambil riwayat jadwal dokter berdasarkan ID
func GetJadwalByDokter(c *gin.Context) {
	idDokter := c.Param("id")
	var jadwals []models.Jadwal

	// Kita gunakan Preload("Ruangan") agar nama ruangan ikut terbawa ke React,
	// dan Order untuk mengurutkan dari yang terbaru (tanggal & jam).
	err := config.DB.Preload("Ruangan").
		Where("id_dokter = ?", idDokter).
		Order("id_jadwal DESC"). // Urutkan dari pengajuan paling baru
		Find(&jadwals).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal mengambil riwayat jadwal"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil",
		"data":    jadwals,
	})
}

// Mengambil SEMUA jadwal beserta data dokter dan ruangan
func GetAllJadwal(c *gin.Context) {
	var jadwals []models.Jadwal
	err := config.DB.Preload("Dokter").Preload("Ruangan").Order("id_jadwal DESC").Find(&jadwals).Error

	if err != nil {
		c.JSON(500, gin.H{"message": "Gagal mengambil data jadwal"})
		return
	}
	c.JSON(200, gin.H{"data": jadwals})
}

// Update status jadwal (Setuju / Tolak)
func UpdateStatusJadwal(c *gin.Context) {
	idJadwal := c.Param("id")
	var input struct {
		StatusJadwal string `json:"status_jadwal"`
		Catatan      string `json:"catatan"` // Jika ditolak
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"message": "Input tidak valid"})
		return
	}

	var jadwal models.Jadwal
	if err := config.DB.First(&jadwal, idJadwal).Error; err != nil {
		c.JSON(404, gin.H{"message": "Jadwal tidak ditemukan"})
		return
	}

	jadwal.StatusJadwal = input.StatusJadwal
	if input.Catatan != "" {
		jadwal.Catatan = input.Catatan
	}

	config.DB.Save(&jadwal)
	c.JSON(200, gin.H{"message": "Status berhasil diupdate"})
}
