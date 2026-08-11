package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var SecretKey = os.Getenv("SECRETE_KEY")

func init() {
	_ = godotenv.Load()
	SecretKey = os.Getenv("SECRETE_KEY")
}

func GenerateToken(userID, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(SecretKey))

	if err != nil {
		return "", err
	}

	return signedToken, nil
}
