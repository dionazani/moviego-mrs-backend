package database

import (
	"log"
	"time"

	"github.com/dionazani/moviego-mrs-backend/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(cfg *config.Config) *gorm.DB {
	// Membuka koneksi ke Postgres menggunakan GORM
	// Kita set logger ke Info agar bisa melihat SQL Query yang dieksekusi selama masa development
	db, err := gorm.Open(postgres.Open(cfg.DBDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Mengatur Connection Pool
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get generic database interface: %v", err)
	}

	// SetMaxIdleConns menetapkan jumlah maksimum koneksi dalam pool koneksi menganggur.
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns menetapkan jumlah maksimum koneksi terbuka ke database.
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime menetapkan jumlah waktu maksimum koneksi dapat digunakan kembali.
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Database connection successfully configured.")
	return db
}
