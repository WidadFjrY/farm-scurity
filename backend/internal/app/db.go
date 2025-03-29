package app

import (
	"farm-scurity/domain/model"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"
)

func NewDB() *gorm.DB {
	dsn := "file:farm-scurity.db?_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)"
	db, err := gorm.Open(sqlite.Dialector{DriverName: "sqlite", DSN: dsn}, &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal menghubungkan ke database: %v", err)
	}

	if err := db.AutoMigrate(&model.Picture{}, &model.History{}, &model.Device{}); err != nil {
		log.Fatalf("Gagal melakukan migrasi database: %v", err)
	}

	log.Println("Database berhasil diinisialisasi!")

	return db
}
