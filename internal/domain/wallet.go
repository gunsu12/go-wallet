package domain

import (
	"time"

	"gorm.io/gorm"
)

type Wallet struct {
	ID          string         `json:"id" gorm:"type:char(36);primaryKey"`
	Name        string         `json:"name" binding:"required"`
	Description string         `json:"description"`
	UserID      string         `json:"-" gorm:"column:user_id"`
	Amount      int64          `json:"amount" binding:"required,numeric"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type WalletRepository interface {
	Create(wallet *Wallet) error
	FindByID(id string) (*Wallet, error)
	FindByUser(user_id string) ([]Wallet, error)
	Update(wallet *Wallet, id string) error
	Delete(id string) error
}

type WalletUsecase interface {
	Create(wallet *Wallet) error
	Update(Wallet *Wallet, id string) error
	FindByID(id string) (*Wallet, error)
	FindByUser(id string) ([]Wallet, error)
	Delete(id string) error
}
