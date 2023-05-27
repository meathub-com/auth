package transport

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"log"
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
	h.mapRoutes()
	h.Server = &http.Server{
		Addr: "0.0.0.0:8080",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      h.Router,
	}
	// return our wonderful handler
	return h
}

func (h *Handler) mapRoutes() {
	h.Router.Get("/auth/{id}", h.GetUser)
	h.Router.Post("/auth", h.PostUser)
	h.Router.Put("/auth/{id}", h.UpdateUser)
	h.Router.Delete("/auth/{id}", h.DeleteUser)
}

func (h *Handler) AliveCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(Response{Message: "I am Alive!"}); err != nil {
		panic(err)
	}
}

func (h *Handler) ReadyCheck(w http.ResponseWriter, r *http.Request) {
	if err := h.Service.ReadyCheck(r.Context()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(Response{Message: "I am Ready!"}); err != nil {
		panic(err)
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
