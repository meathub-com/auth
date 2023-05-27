package user

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

func (s *Service) GenerateToken(user User) (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tokenString, err := token.SignedString([]byte("your-secret-key"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
func (s *Service) GenerateRefreshToken(user User) (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()

	tokenString, err := token.SignedString([]byte("your-secret-key"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}