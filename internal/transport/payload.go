package transport

import "auth/internal/user"

type LoginResponse struct {
	User         user.User `json:"user"`
	AccessToken  string    `json:"authToken"`
	RefreshToken string    `json:"refreshToken"`
}
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func convertRegisterRequestToUser(r RegisterRequest) user.User {
	return user.User{
		Email:    r.Email,
		Password: r.Password,
	}
}
func convertLoginRequestToUser(r LoginRequest) user.User {
	return user.User{
		Email:    r.Email,
		Password: r.Password,
	}
}
