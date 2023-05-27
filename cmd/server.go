// @title Swagger Example API
// @version 1.0
// @description This is a sample server Auth server.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /

package main

import (
	"auth/internal/database"
	"auth/internal/transport"
	"auth/internal/user"
	log "github.com/sirupsen/logrus"
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
