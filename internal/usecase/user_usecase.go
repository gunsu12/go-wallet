package usecase

import (
	"errors"

	"github.com/google/uuid"
	"github.com/gunsu12/go-wallet/config"
	"github.com/gunsu12/go-wallet/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	repo domain.UserRepository
}

// NewUserUsecase adalah constructor untuk UserUsecase
func NewUserUsecase(r domain.UserRepository) *UserUsecase {
	return &UserUsecase{
		repo: r,
	}
}

// Register membuat user baru dengan validasi dan hash password
func (uc *UserUsecase) Register(user *domain.User) error {
	// Validasi apakah username sudah digunakan
	existingUserByUsername, err := uc.repo.FindByUsername(user.Username)
	if err != nil {
		return err
	}
	if existingUserByUsername != nil {
		return errors.New("username sudah digunakan")
	}

	// Validasi apakah email sudah digunakan
	existingUserByEmail, err := uc.repo.FindByEmail(user.Email)
	if err != nil {
		return err
	}
	if existingUserByEmail != nil {
		return errors.New("email sudah digunakan")
	}

	// Generate UUID untuk ID user
	user.ID = uuid.NewString() // Bisa ganti ke UUID v7 jika pakai lib yang support

	// Hash password sebelum disimpan
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("gagal mengenkripsi password")
	}
	user.Password = string(hashedPassword)

	// Simpan ke database
	err = uc.repo.Create(user)
	if err != nil {
		return err
	}

	return nil
}

// GetByID mengambil user berdasarkan ID
func (uc *UserUsecase) GetByID(id string) (*domain.User, error) {
	user, err := uc.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user tidak ditemukan")
	}
	return user, nil
}

// usecase login
func (uc *UserUsecase) Login(email, password string) (*domain.User, error) {
	// Cari user berdasarkan email
	user, err := uc.repo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("email tidak terdaftar")
	}

	// Cocokkan password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("password salah")
	}

	// Buat token
	token, err := config.GenerateJWT(user.ID, user.Email)
	if err != nil {
		return nil, errors.New("gagal generate token")
	}

	// Tambahkan token ke user
	user.Token = token

	return user, nil
}
