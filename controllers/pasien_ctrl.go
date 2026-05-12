package controllers

import (
	"fmt"
	"net/http"
	"sirs-backend/config"
	"sirs-backend/models"
	"time"

	"github.com/gin-gonic/gin"
)

type PasienInput struct {
	NamaPasien string `json:"nama_pasien"`
	NoHp       string `json:"no_hp"`
}

// CreatePasien: Mendaftarkan Pasien Baru dan generate RM
func CreatePasien(c *gin.Context) {
	var input PasienInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Buat object pasien
	pasienBaru := models.Pasien{
		NamaPasien: input.NamaPasien,
		TipePasien: "Pasien Baru", // Otomatis di-set Pasien Baru
		NoHp:       input.NoHp,
	}

	// Generate Nomor Rekam Medis sementara (Format: RM-TahunBulanTanggal-Random)
	// Catatan: GORM akan auto-increment ID setelah data di-save, jadi kita save dulu
	// dengan RM sementara, lalu update RM-nya.
	pasienBaru.NomorRekamMedis = fmt.Sprintf("TEMP-%d", time.Now().Unix())

	if err := config.DB.Create(&pasienBaru).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendaftarkan pasien"})
		return
	}

	// Update dengan RM yang rapi (RM - Tahun - ID)
	rmResmi := fmt.Sprintf("RM-%d-%04d", time.Now().Year(), pasienBaru.IdPasien)
	config.DB.Model(&pasienBaru).Update("nomor_rekam_medis", rmResmi)
	pasienBaru.NomorRekamMedis = rmResmi // Update response object

	c.JSON(http.StatusOK, gin.H{
		"message": "Pasien berhasil didaftarkan",
		"data":    pasienBaru,
	})

}
func GetAllPasien(c *gin.Context) {
	var pasiens []models.Pasien

	// config.DB.Find(&pasiens) akan mengambil seluruh isi tabel pasien
	if err := config.DB.Find(&pasiens).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pasien"})
		return
	}

	c.JSON(http.StatusOK, pasiens)
}
