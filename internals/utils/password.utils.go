package utils

import (
	"crypto/sha512"
	"fmt"
)

func GeneratePasswordHash(password, salt string) string {
	hash := sha512.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func DoPasswordMatch(hashedPassword, currentPassword, salt string) bool {
	var currentPasswordHash = GeneratePasswordHash(currentPassword, salt)
	return currentPasswordHash == hashedPassword
}
