package main

import (
	log "github.com/sirupsen/logrus"
	"users/internal/database"
	"users/internal/transport"
	"users/internal/user"
)

func Run() error {
	var err error
	store, err := database.NewDatabase()
	if err != nil {
		log.WithError(err).Error("could not create database")
		return err
	}
	err = store.MigrateDB()
	if err != nil {
		log.WithError(err).Error("could not migrate database")
		return err
	}
	log.Info("database migrated")
	log.Info("creating new user service")
	userService := user.NewService(store)
	log.Info("creating new transport handler")
	handler := transport.NewHandler(userService)
	log.Info("starting server")
	if err := handler.Serve(); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := Run(); err != nil {
		log.WithError(err).Error("could not run server")
	}
}
