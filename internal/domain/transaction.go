package domain

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	ID                     string         `json:"id" gorm:"type:char(36);primaryKey"`
	WalletID               string         `json:"wallet_id"`
	TransactionAmount      int64          `jsom:"transaction_amount" binding:"required,numeric"`
	TransactionType        string         `json:"transaction_type"`
	TransactionDescription string         `json:"transaction_description"`
	CreatedAt              time.Time      `json:"created_at"`
	UpdatedAt              time.Time      `json:"updated_at"`
	DeletedAt              gorm.DeletedAt `json:"-" gorm:"index"`
}

type TransactionRepository interface {
	Create(trs *Transaction) error
	FindByID(id string) (*Transaction, error)
	FindByWallet(walet_id string) ([]Transaction, error)
	FindByUser(user_id string) ([]Transaction, error)
	Cancel(id string) error
}

type TransactionUsecase interface {
	Create(trs *Transaction) error
	FindByID(id string) (*Transaction, error)
	FindByWallet(walet_id string) ([]Transaction, error)
	FindByUser(user_id string) ([]Transaction, error)
	Cancel(id string) error
}
