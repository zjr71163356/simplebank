package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassWord(password string) (string, error) {
	hashedPassWord, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("fail to hash password %w", err)
	}
	return string(hashedPassWord), err

}

func MatchPassWord(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
