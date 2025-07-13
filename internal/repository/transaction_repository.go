package repository

import (
	"errors"

	"github.com/gunsu12/go-wallet/internal/domain"
	"gorm.io/gorm"
)

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) domain.TransactionRepository {
	return &transactionRepository{db}
}

func (r *transactionRepository) Create(trs *domain.Transaction) error {
	if err := r.db.Create(trs).Error; err != nil {
		return err
	}
	return nil
}

func (r *transactionRepository) FindByID(id string) (*domain.Transaction, error) {
	var transaction domain.Transaction

	result := r.db.First(&transaction, "id = ?", id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &transaction, nil
}

func (r *transactionRepository) FindByWallet(walet_id string) ([]domain.Transaction, error) {
	var transaction []domain.Transaction

	result := r.db.Where("wallet_id = ?", walet_id).Find(&transaction)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return transaction, nil
}

func (r *transactionRepository) FindByUser(user_id string) ([]domain.Transaction, error) {
	var transactions []domain.Transaction

	// Join dengan wallets dan filter berdasarkan wallets.user_id
	err := r.db.
		Joins("JOIN wallets ON wallets.id = transactions.wallet_id").
		Where("wallets.user_id = ?", user_id).
		Find(&transactions).Error

	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (r *transactionRepository) Delete(trs_id string) error {
	result := r.db.Where("id = ?", trs_id).Delete(&domain.Transaction{}).Error

	return result
}
