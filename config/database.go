package config

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() { // Sesuaikan dengan nama fungsi Anda (misal: ConnectDatabase)
	// 1. Ganti <PASSWORD> dengan password asli TiDB Anda
	// 2. Perhatikan kata "sirs_penjadwalan", pastikan ini sama dengan nama database yang Anda buat di TiDB
	dsn := "2ctUt9uw377dKfr.root:cYaBOt8fnRXPCL7G@tcp(gateway01.ap-southeast-1.prod.aws.tidbcloud.com:4000)/sirs_penjadwalan?charset=utf8mb4&parseTime=True&loc=Local&tls=true"

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal terhubung ke TiDB Cloud:", err)
	}

	DB = database
	log.Println("Berhasil terhubung ke TiDB Cloud!")
}
