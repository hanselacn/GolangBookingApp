package encrypt

import "golang.org/x/crypto/bcrypt"

type Bcrypt interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(hash string, password string) bool
}

func HashPassword(password string) (string, error) {
	passwordBytes := []byte(password)
	hashedPassword, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.MinCost)
	return string(hashedPassword), err
}

func CheckPasswordHash(hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
