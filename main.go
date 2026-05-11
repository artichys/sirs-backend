package main

import (
	"fmt" // Pastikan fmt di-import
	"os"
	"sirs-backend/config"
	"sirs-backend/routes"
)

func main() {
	fmt.Println("1. Memulai server...")

	fmt.Println("2. Mencoba koneksi ke database MySQL...")
	config.ConnectDB()
	fmt.Println("3. Database berhasil terhubung!")

	fmt.Println("4. Menyiapkan rute API...")
	r := routes.SetupRouter()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Jika di lokal (tidak ada env PORT), pakai 8080
	}

	fmt.Println("5. Server siap berjalan di port " + port)
	r.Run(":" + port)
}
