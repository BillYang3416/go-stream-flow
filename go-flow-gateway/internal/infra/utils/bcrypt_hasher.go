package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type BcryptHasher struct{}

func NewBcryptHasher() *BcryptHasher {
	return &BcryptHasher{}
}

func (bh *BcryptHasher) GenerateHash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("BcryptHasher - GenerateHash - bcrypt.GenerateFromPassword: %w", err)
	}

	return string(hashedPassword), nil
}

func (bh *BcryptHasher) CompareHash(password, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return fmt.Errorf("BcryptHasher - CompareHashAndPassword: %w", err)
	}

	return nil
}
