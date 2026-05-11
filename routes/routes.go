package routes

import (
    "sirs-backend/controllers"
    "time"

    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
    r := gin.Default()

    // Pengaturan CORS yang mengizinkan metode PUT
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:5173"}, // Sesuaikan dengan port React
        AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))

    api := r.Group("/api")
    {
        // ---- FASE 1: Jadwal Dokter ----
        api.POST("/jadwal", controllers.CreateJadwal)
        api.GET("/jadwal", controllers.GetAllJadwal) // Menggunakan fungsi baru yang ada Preload-nya
        api.GET("/jadwal/dokter/:id", controllers.GetJadwalByDokter)
        api.PUT("/jadwal/:id/status", controllers.UpdateStatusJadwal) // Fungsi dinamis untuk Approve/Reject

        // ---- FASE 2: Pasien & Pendaftaran ----
        api.POST("/pasien", controllers.CreatePasien)
        api.POST("/booking", controllers.CreateBooking)

        // ---- FASE 3: Layar Monitor ----
        api.GET("/monitor/antrean/:id_jadwal", controllers.GetAntreanMonitor)

        // ---- FASE 4: Aksi Ruang Poli ----
        api.PUT("/booking/:id/status", controllers.UpdateStatusBooking)

        // ---- FASE 5: Laporan Manajemen ----
        api.GET("/laporan/utilisasi", controllers.GetLaporanUtilisasi)
    }

    return r
}