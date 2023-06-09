package transport

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	log "github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Handler struct {
	Router  *chi.Mux
	Service UserService
	Server  *http.Server
}

type Response struct {
	Message string `json:"message"`
}

func NewHandler(service UserService) *Handler {
	h := &Handler{
		Service: service,
		Router:  chi.NewRouter(),
	}

	// Configure CORS
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Change "*" to the appropriate origin URL(s)
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	})
	h.Router.Use(cors.Handler)

	h.mapRoutes()

	h.Server = &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      h.Router,
	}
	return h
}

func (h *Handler) mapRoutes() {
	h.Router.Get("/auth", h.AliveCheck)
	h.Router.Get("/auth/{id}", h.GetUser)
	h.Router.Post("/auth/register", h.RegisterUser)
	h.Router.Put("/auth/{id}", h.UpdateUser)
	h.Router.Delete("/auth/{id}", h.DeleteUser)
	h.Router.Post("/auth/login", h.LoginUser)
	h.Router.Get("/auth/refresh/{refreshToken}", h.RefreshToken)
	h.Router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))
}

func (h *Handler) AliveCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(Response{Message: "I am Alive!"}); err != nil {
		log.Errorf("Error getting profile: %v", err)
	}
}

func (h *Handler) ReadyCheck(w http.ResponseWriter, r *http.Request) {
	if err := h.Service.ReadyCheck(r.Context()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(Response{Message: "I am Ready!"}); err != nil {
		log.Errorf("Error getting profile: %v", err)
	}
}

// Serve - gracefully serves our newly set up handler function
func (h *Handler) Serve() error {
	go func() {
		if err := h.Server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	h.Server.Shutdown(ctx)

	log.Println("shutting down gracefully")
	return nil
}
