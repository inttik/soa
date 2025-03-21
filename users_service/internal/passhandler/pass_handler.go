package passhandler

import (
	"encoding/hex"
	"log"
	"users/oas"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password oas.PasswordString) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	hexPassword := hex.EncodeToString(hashedPassword)
	return hexPassword, nil
}

func ComparePassword(password string, hexPassword string) bool {
	hashedPassword, err := hex.DecodeString(hexPassword)
	if err != nil {
		log.Fatal(err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
