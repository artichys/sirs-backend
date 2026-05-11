package main

import (
	"fmt" // Pastikan fmt di-import
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

	fmt.Println("5. Server siap berjalan di port 8080...")
	r.Run(":8080")
}
