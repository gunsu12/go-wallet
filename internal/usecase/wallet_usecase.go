package usecase

import (
	"errors"

	"github.com/google/uuid"
	"github.com/gunsu12/go-wallet/internal/domain"
)

type WalletUsecase struct {
	repo domain.WalletRepository
}

// NewWalletUsecase membuat instance baru WalletUsecase
func NewWalletUsecase(repo domain.WalletRepository) *WalletUsecase {
	return &WalletUsecase{repo: repo}
}

// Create menambahkan wallet baru
func (uc *WalletUsecase) Create(wallet *domain.Wallet) error {
	// Validasi dasar (opsional)
	if wallet.Name == "" || wallet.UserID == "" {
		return errors.New("nama wallet atau user id tidak boleh kosong") // bisa buat custom error
	}

	wallet.ID = uuid.New().String()

	return uc.repo.Create(wallet)
}

// Update memperbarui wallet berdasarkan ID
func (uc *WalletUsecase) Update(wallet *domain.Wallet, id string) error {
	return uc.repo.Update(wallet, id)
}

// FindByID mencari wallet berdasarkan ID
func (uc *WalletUsecase) FindByID(id string) (*domain.Wallet, error) {
	return uc.repo.FindByID(id)
}

// FindByUser mencari wallet berdasarkan UserID
func (uc *WalletUsecase) FindByUser(userID string) ([]domain.Wallet, error) {
	return uc.repo.FindByUser(userID)
}

// fungsi delete wallet
func (uc *WalletUsecase) Delete(id string) error {
	return uc.repo.Delete(id)
}
