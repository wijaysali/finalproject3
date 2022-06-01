package utils

import (
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func EncryptToBcrypt(plaintext *string) string {
	hashed, err := bcrypt.GenerateFromPassword([]byte(*plaintext), bcrypt.DefaultCost)
	if err != nil {
		ResponsePanic(fiber.StatusInternalServerError, err.Error())
	}
	return string(hashed)
}

func IsHashBcryptMatch(password *string, hashedPassword *string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(*hashedPassword), []byte(*password))
	return err == nil
}
