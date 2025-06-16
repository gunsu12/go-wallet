package domain

import (
	"time"
)

type User struct {
	ID          string    `json:"id" gorm:"type:char(36);primaryKey"`
	Name        string    `json:"name"`
	Username    string    `json:"username" gorm:"unique"` // baru
	Email       string    `json:"email" gorm:"unique"`
	PhoneNumber string    `json:"phone_number"` // baru
	Password    string    `json:"password"`     // disembunyikan di JSON
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// tidak disimpan ke database
	Token string `json:"token,omitempty" gorm:"-"`
}

type UserRepository interface {
	Create(user *User) error
	FindByID(id string) (*User, error)
	FindByEmail(email string) (*User, error)
	FindByUsername(username string) (*User, error)
}

type UserUsecase interface {
	Register(user *User) error
	GetByID(id string) (*User, error)
}
