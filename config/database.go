package config

import (
	"log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	// Sesuaikan dengan username dan password mysql Anda (default xampp: root, tanpa password)
	dsn := "root:@tcp(127.0.0.1:3306)/sirs_penjadwalan?charset=utf8mb4&parseTime=True&loc=Local"
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Gagal koneksi ke database:", err)
	}

	DB = database
	log.Println("Database terkoneksi dengan sukses!")
}