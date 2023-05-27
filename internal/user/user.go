package user

import (
	"context"
)

type User struct {
	ID       string
	Email    string
	Password string
}

type UserStore interface {
	GetUser(context.Context, string) (User, error)
	PostUser(context.Context, User) (User, error)
	UpdateUser(context.Context, User) (User, error)
	DeleteUser(context.Context, string) error
	Ping(ctx context.Context) error
}
type Service struct {
	Store UserStore
}

func (s *Service) GetUser(ctx context.Context, id string) (User, error) {
	return s.Store.GetUser(ctx, id)
}

func (s *Service) PostUser(ctx context.Context, user User) (User, error) {
	return s.Store.PostUser(ctx, user)
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

func NewService(store UserStore) *Service {
	return &Service{
		Store: store,
	}
}
