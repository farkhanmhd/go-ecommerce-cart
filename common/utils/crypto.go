package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(passwordStr string) string {
	hashedByteSlice, err := bcrypt.GenerateFromPassword([]byte(passwordStr), 12)
	PanicIfError(err)
	return string(hashedByteSlice)
}

func ComparePassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return err
	}
	return nil
}
