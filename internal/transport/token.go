package transport

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v4"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// RefreshToken godoc
// @Summary Get refresh token for user
// @Description Get refresh token for user
// @Tags users
// @Accept  json
// @Produce  json
// @Param refreshToken path string true "Refresh token"
// @Success 200 {object} string
// @Router /auth/refresh/{refreshToken} [get]
func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken := chi.URLParam(r, "refreshToken")
	token, err := validateToken(refreshToken)
	if err != nil {
		log.WithError(err).Error("error validating refresh token")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid refresh token"))
		return
	}

	userID, ok := token["sub"].(string)
	if !ok {
		log.WithError(err).Error("error validating refresh token")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid refresh token"))
		return
	}

	usr, err := h.Service.GetUser(r.Context(), userID)
	if err != nil {
		log.WithError(err).Error("error getting user")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("user not found"))
		return
	}
	accessToken, err := h.Service.GenerateToken(usr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error generating refresh token"))
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"accessToken": accessToken}); err != nil {
		log.Errorf("Error getting profile: %v", err)
	}
}

// validateToken - validates an incoming JWT token
func validateToken(accessToken string) (map[string]interface{}, error) {
	var mySigningKey = []byte("missionimpossible")
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("could not validate auth token")
		}
		return mySigningKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("could not validate auth token")
}
