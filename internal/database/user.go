package database

import "users/internal/user"

type UserRow struct {
	ID       string
	Email    string
	Password string
}

func convertUserRowToUser(userRow UserRow) user.User {
	return user.User{
		ID:       userRow.ID,
		Email:    userRow.Email,
		Password: userRow.Password,
	}
}
