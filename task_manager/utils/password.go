package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes := []byte(password)

	hash, err := bcrypt.GenerateFromPassword(bytes, bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func CheckPassword(hash, password string) bool {
	bytes := []byte(password)

	err := bcrypt.CompareHashAndPassword([]byte(hash), bytes)

	if err != nil {
		return false
	}

	return true
}
