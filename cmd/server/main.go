package main

import (
	"github.com/salty-max/go-comments/pkg/comment"
	db "github.com/salty-max/go-comments/pkg/database"
	transportHttp "github.com/salty-max/go-comments/pkg/transport/http"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

// Run - is responsible for the instantiation
// and startup of the app
func Run() error {
	log.SetFormatter(&log.JSONFormatter{})
	log.Info("Setting up server...")

	var err error

	if err := godotenv.Load(); err != nil {
		log.Error("failed to load env")
		return err
	}

	store, err := db.NewDatabase()
	if err != nil {
		log.Error("failed to connect to the database")
		return err
	}

	if err := store.MigrateDB(); err != nil {
		log.Error("failed to migrate database")
		return err
	}

	log.Info("successfully connected and pinged database")

	cmtService := comment.NewService(store)

	httpHandler := transportHttp.NewHandler(cmtService)
	if err := httpHandler.Serve(); err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func main() {
	if err := Run(); err != nil {
		log.Error(err)
	}
}
