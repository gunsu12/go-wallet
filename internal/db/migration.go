package db

import (
	"log"

	"github.com/gunsu12/go-wallet/internal/domain"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) {
	err := db.AutoMigrate(
		&domain.User{},
		// nanti tambahkan struct lain seperti Wallet, Transaction, Budget
	)

	if err != nil {
		log.Fatalf("Gagal melakukan migrasi: %v", err)
	}

	log.Println("Migrasi database berhasil ✅")
}
