package usecase

import (
	"errors"

	"github.com/google/uuid"
	"github.com/gunsu12/go-wallet/internal/domain"
)

type TransactionUsecase struct {
	repo       domain.TransactionRepository
	walletRepo domain.WalletRepository
}

func NewTransactionUsecase(repo domain.TransactionRepository) *TransactionUsecase {
	return &TransactionUsecase{repo: repo}
}

func (uc *TransactionUsecase) Create(trs *domain.Transaction) error {
	trs.ID = uuid.New().String()

	// 1. Ambil wallet terlebih dahulu
	wallet, err := uc.walletRepo.FindByID(trs.WalletID)
	if err != nil {
		return err
	}
	if wallet == nil {
		return errors.New("wallet tidak ditemukan")
	}

	// 2. Hitung saldo baru berdasarkan tipe transaksi
	switch trs.TransactionType {
	case "debit":
		if wallet.Amount < trs.TransactionAmount {
			return errors.New("saldo tidak mencukupi untuk debit")
		}
		wallet.Amount -= trs.TransactionAmount
	case "credit":
		wallet.Amount += trs.TransactionAmount
	default:
		return errors.New("jenis transaksi tidak valid: harus debit atau credit")
	}

	// 3. Simpan transaksi
	if err := uc.repo.Create(trs); err != nil {
		return err
	}

	// 4. Simpan kembali saldo wallet yang baru
	if err := uc.walletRepo.Update(wallet, wallet.ID); err != nil {
		// Jika update wallet gagal, rollback transaksi yang sudah dibuat
		_ = uc.repo.Delete(trs.ID)
		return err
	}

	return nil
}

func (uc *TransactionUsecase) FindByID(id string) (*domain.Transaction, error) {
	return uc.repo.FindByID(id)
}

func (uc *TransactionUsecase) FindByUser(id string) ([]domain.Transaction, error) {
	return uc.repo.FindByUser(id)
}

func (uc *TransactionUsecase) FindByWallet(id string) ([]domain.Transaction, error) {
	return uc.repo.FindByWallet(id)
}
