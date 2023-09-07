package models

import (
	"mini-account-service/db/pwhasher"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Uuid  string `gorm:"unique;not null;size:64"`
	Email string `gorm:"unique;not null;size:128"`

	LoginFuncVersion uint32 `gorm:"not null"`
	Salt             []byte `gorm:"not null;size:64"`
	PasswordHash     []byte `gorm:"not null;size:64"`
}

func NewUser(email string, password string) (*User, error) {
	hasher := pwhasher.NewPasswordHasher(password)
	hashed, err := hasher.Hash()
	if err != nil {
		return nil, err
	}

	return &User{
		Uuid:             uuid.NewString(),
		Email:            email,
		LoginFuncVersion: hashed.HashFuncVersion,
		Salt:             hashed.Salt,
		PasswordHash:     hashed.DK,
	}, nil
}
