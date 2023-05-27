package database

import (
	"crypto/rand"
	"encoding/base64"
	"golang.org/x/crypto/bcrypt"
	"users/internal/user"
)

type UserRow struct {
	ID       string
	Email    string
	Password string
	Salt     string
}

func convertUserRowToUser(userRow UserRow) user.User {
	return user.User{
		ID:       userRow.ID,
		Email:    userRow.Email,
		Password: userRow.Password,
	}
}
func generateRandomSalt(length int) (string, error) {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(salt), nil
}
func encryptPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}
