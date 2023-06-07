package user

import (
	jwt "github.com/golang-jwt/jwt/v4"
	"time"
)

func (s *Service) GenerateToken(user User) (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	claims["iat"] = time.Now().Unix()
	claims["sub"] = user.ID

	tokenString, err := token.SignedString([]byte("missionimpossible"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
func (s *Service) GenerateRefreshToken(user User) (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()
	claims["iat"] = time.Now().Unix()
	claims["sub"] = user.ID

	tokenString, err := token.SignedString([]byte("missionimpossible"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
