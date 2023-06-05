package user

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrorUserExists = errors.New("error creating user")
)

// User godoc
// @Description User's login details
// @Param   username     string     "Username to use for login"
// @Param   password     string     "User's password"
type User struct {
	ID       string
	Email    string
	Password string
}

type UserStore interface {
	GetUserAndSaltByEmail(context.Context, string) (User, string, error)
	GetUser(context.Context, string) (User, error)
	GetUserByEmail(context.Context, string) (User, error)
	PostUser(context.Context, User) (User, error)
	UpdateUser(context.Context, User) (User, error)
	DeleteUser(context.Context, string) error
	Ping(ctx context.Context) error
}
type Service struct {
	Store UserStore
}

func NewService(store UserStore) *Service {
	return &Service{
		Store: store,
	}
}
func (s *Service) GetUser(ctx context.Context, id string) (User, error) {
	return s.Store.GetUser(ctx, id)
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (User, error) {
	return s.Store.GetUserByEmail(ctx, email)
}
func (s *Service) PostUser(ctx context.Context, usr User) (User, error) {
	u, err := s.Store.PostUser(ctx, usr)
	if err != nil && errors.Is(err, ErrorUserExists) {
		return User{}, ErrorUserExists
	}
	if err != nil {
		wrappedErr := fmt.Errorf("error creating user: %w", err)
		return User{}, wrappedErr
	}

	return u, nil
}

func (s *Service) UpdateUser(ctx context.Context, user User) (User, error) {
	return s.Store.UpdateUser(ctx, user)
}

func (s *Service) DeleteUser(ctx context.Context, id string) error {
	return s.Store.DeleteUser(ctx, id)
}

func (s *Service) ReadyCheck(ctx context.Context) error {
	return s.Store.Ping(ctx)
}
func (s *Service) Login(ctx context.Context, email string, password string) (User, error) {
	user, salt, err := s.Store.GetUserAndSaltByEmail(ctx, email)

	if err != nil {
		return User{}, err
	}
	if !validatePassword(password, salt, user.Password) {
		return User{}, fmt.Errorf("invalid password for user %s", email)
	}
	return user, nil
}
func validatePassword(cleanPassword, salt, encryptedPassword string) bool {
	saltedPassword := cleanPassword + salt

	err := bcrypt.CompareHashAndPassword([]byte(encryptedPassword), []byte(saltedPassword))
	if err != nil {
		return false
	}

	return true
}
