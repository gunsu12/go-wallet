package repository

import (
	"errors"

	"github.com/gunsu12/go-wallet/internal/domain"
	"gorm.io/gorm"
)

type walletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) domain.WalletRepository {
	return &walletRepository{db}
}

func (r *walletRepository) Create(wallet *domain.Wallet) error {
	if err := r.db.Create(wallet).Error; err != nil {
		return err
	}
	return nil
}

func (r *walletRepository) FindByID(id string) (*domain.Wallet, error) {
	var wallet domain.Wallet

	result := r.db.First(&wallet, "id = ?", id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &wallet, nil
}

func (r *walletRepository) FindByUser(userId string) ([]domain.Wallet, error) {
	var wallets []domain.Wallet

	result := r.db.Where("user_id = ?", userId).Find(&wallets)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return wallets, nil
}

func (r *walletRepository) Update(wallet *domain.Wallet, id string) error {

	result := r.db.Model(&domain.Wallet{}).Where("id = ?", id).Updates(wallet)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *walletRepository) Delete(id string) error {
	result := r.db.Where("id = ?", id).Delete(&domain.Wallet{})
	return result.Error
}
