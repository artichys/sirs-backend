package controllers

import (
	"net/http"
	"sirs-backend/config"
	"sirs-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 1. UPDATE STRUKTUR REQUEST
type BookingRequest struct {
	IdJadwal        int    `json:"id_jadwal"`
	NomorRekamMedis string `json:"nomor_rekam_medis"` // Diubah agar cocok dengan React
}

func CreateBooking(c *gin.Context) {
	var req BookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Input tidak valid", "error": err.Error()})
		return
	}

	var bookingBaru models.Booking
	var pasien models.Pasien
	var namaRuangan string = "POLI UMUM" // Fallback default
	var namaDokter string = "Dokter Jaga"

	// Memulai Database Transaction
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		var jadwal models.Jadwal

		// Pastikan jadwal ada dan sudah Confirmed
		if err := tx.Where("id_jadwal = ? AND status_jadwal = 'Confirmed'", req.IdJadwal).First(&jadwal).Error; err != nil {
			return err // Jadwal tidak valid
		}

		// Cari Pasien berdasarkan Nomor Rekam Medis dari React
		if err := tx.Where("nomor_rekam_medis = ?", req.NomorRekamMedis).First(&pasien).Error; err != nil {
			return err // Pasien tidak ditemukan, batalkan transaksi
		}

		// Hitung jumlah pasien yang sudah mendaftar (untuk nomor antrean)
		var jumlahAntrean int64
		tx.Model(&models.Booking{}).Where("id_jadwal = ?", req.IdJadwal).Count(&jumlahAntrean)

		// Buat record booking baru menggunakan IdPasien yang baru saja ditemukan
		bookingBaru = models.Booking{
			IdJadwal:      req.IdJadwal,
			IdPasien:      pasien.IdPasien, // Data valid dari database
			NomorAntrean:  int(jumlahAntrean) + 1,
			StatusBooking: "Menunggu",
		}

		// Simpan ke database
		if err := tx.Create(&bookingBaru).Error; err != nil {
			return err
		}

		// (Opsional) Ambil nama ruangan dan dokter untuk dicetak di struk
		type InfoJadwal struct {
			NamaRuangan string
			NamaDokter  string
		}
		var info InfoJadwal
		tx.Table("jadwal").
			Select("ruangan.nama_ruangan, dokter.nama_dokter").
			Joins("left join ruangan on ruangan.id_ruangan = jadwal.id_ruangan").
			Joins("left join dokter on dokter.id_dokter = jadwal.id_dokter").
			Where("jadwal.id_jadwal = ?", req.IdJadwal).
			Scan(&info)

		if info.NamaRuangan != "" {
			namaRuangan = info.NamaRuangan
		}
		if info.NamaDokter != "" {
			namaDokter = info.NamaDokter
		}

		return nil // Transaksi sukses di-commit
	})

	// Jika terjadi error di dalam blok Transaction (misal RM tidak ditemukan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal membuat antrean: " + err.Error()})
		return
	}

	// 2. KEMBALIKAN RESPONS LENGKAP UNTUK STRUK REACT
	c.JSON(http.StatusOK, gin.H{
		"message": "Booking berhasil, silahkan cek layar monitor antrean",
		"data": gin.H{
			"nomor_antrean": bookingBaru.NomorAntrean,
			"nama_pasien":   pasien.NamaPasien,
			"nama_ruangan":  namaRuangan,
			"nama_dokter":   namaDokter,
		},
	})
}

// ==============================================================
// FUNGSI UPDATE STATUS (Tetap sama seperti aslinya, tidak diubah)
// ==============================================================

type UpdateStatusInput struct {
	StatusBooking string `json:"status_booking" binding:"required"`
}

func UpdateStatusBooking(c *gin.Context) {
	idBooking := c.Param("id")
	var input UpdateStatusInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status booking wajib diisi (Menunggu, Diperiksa, atau Selesai)"})
		return
	}

	if input.StatusBooking != "Menunggu" && input.StatusBooking != "Diperiksa" && input.StatusBooking != "Selesai" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status tidak valid"})
		return
	}

	var booking models.Booking
	if err := config.DB.First(&booking, idBooking).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data antrean tidak ditemukan"})
		return
	}

	booking.StatusBooking = input.StatusBooking
	config.DB.Save(&booking)

	c.JSON(http.StatusOK, gin.H{
		"message": "Status pasien berhasil diperbarui menjadi " + input.StatusBooking,
		"data":    booking,
	})
}
