package transport

import (
	"auth/internal/database"
	"auth/internal/user"
	"context"
	"encoding/json"
	"errors"
	chi "github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type UserService interface {
	GetUser(ctx context.Context, id string) (user.User, error)
	GetUserByEmail(ctx context.Context, email string) (user.User, error)
	Login(ctx context.Context, email string, password string) (user.User, error)
	PostUser(ctx context.Context, user user.User) (user.User, error)
	UpdateUser(ctx context.Context, user user.User) (user.User, error)
	DeleteUser(ctx context.Context, id string) error
	ReadyCheck(ctx context.Context) error
	GenerateToken(user user.User) (string, error)
	GenerateRefreshToken(user user.User) (string, error)
}

// GetUser @Summary Get a user
// @Description get string by ID
// @ID get-string-by-int
// @Accept  json
// @Produce  json
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	user, err := h.Service.GetUser(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		log.WithError(err).Error("error getting user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Errorf("Error getting profile: %v", err)
	}
}

// RegisterUser godoc
// @Summary Register a new user
// @Description Register a new user by username and password
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body RegisterRequest true "User info"
// @Success 200 {object} RegisterResponse
// @Router /auth/register [post]
func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var rr RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&rr); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	usr, err := h.Service.PostUser(r.Context(), convertRegisterRequestToUser(rr))
	if err != nil && errors.Is(err, user.ErrorUserExists) {
		log.WithError(err).Error("error creating user")
		errMsg := "User already exists"
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errMsg))
		return
	}
	if err != nil && errors.Is(err, database.ErrorUserExists) {
		log.WithError(err).Error("error creating user")
		errMsg := "User already exists"
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errMsg))
		return
	}

	accessToken, err := h.Service.GenerateToken(usr)
	if err != nil {
		log.WithError(err).Error("error creating user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	refreshToken, err := h.Service.GenerateRefreshToken(usr)
	if err != nil {
		log.WithError(err).Error("error creating user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(RegisterResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         usr,
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

}

// LoginUser godoc
// @Summary Log in a user
// @Description Log in a user by username and password
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body user.User true "Login user"
// @Success 200 {object} LoginResponse
// @Router /auth/login [post]
func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var lr LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&lr); err != nil {
		log.WithError(err).Error("error decoding user")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	usr := convertLoginRequestToUser(lr)
	usr, err := h.Service.Login(r.Context(), usr.Email, usr.Password)
	if err != nil {
		log.WithError(err).Error("invalid login or password")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid login or password"))
		return
	}
	accessToken, err := h.Service.GenerateToken(usr)
	if err != nil {
		log.WithError(err).Error("error generating token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	refreshToken, err := h.Service.GenerateToken(usr)
	if err != nil {
		log.WithError(err).Error("error generating refresh token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         usr,
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
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
		log.Errorf("Error getting profile: %v", err)
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
