package helper

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func VerifyPassword(userPassword, passwordRequest string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(passwordRequest))
	return err == nil
}
