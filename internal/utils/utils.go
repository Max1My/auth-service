package utils

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

func VerifyPassword(hashedPassword string, candidatePassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(candidatePassword))
	return err == nil
}

// HashPassword генерирует хэш для пароля
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	log.Printf("Generated hash for password: %s", string(bytes)) // Логируем хэш
	return string(bytes), nil
}
