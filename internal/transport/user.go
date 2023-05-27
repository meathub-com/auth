package transport

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"net/http"
	"users/internal/user"
)

type UserService interface {
	GetUser(ctx context.Context, id string) (user.User, error)
	GetUserByEmail(ctx context.Context, email string) (user.User, error)
	PostUser(ctx context.Context, user user.User) (user.User, error)
	UpdateUser(ctx context.Context, user user.User) (user.User, error)
	DeleteUser(ctx context.Context, id string) error
	ReadyCheck(ctx context.Context) error
	GenerateToken(user user.User) (string, error)
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	user, err := h.Service.GetUser(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		log.WithError(err).Error("error getting user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		panic(err)
	}
}

func (h *Handler) PostUser(w http.ResponseWriter, r *http.Request) {
	var user user.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user, err := h.Service.PostUser(r.Context(), user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		panic(err)
	}
}
func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var user user.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user, err := h.Service.GetUserByEmail(r.Context(), user.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	token, err := h.Service.GenerateToken(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(LoginResponse{User: user, Token: token}); err != nil {
		panic(err)
	}
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user user.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user, err := h.Service.UpdateUser(r.Context(), user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		panic(err)
	}
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	err := h.Service.DeleteUser(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

type LoginResponse struct {
	User  user.User `json:"user"`
	Token string    `json:"token"`
}
