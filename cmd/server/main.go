package main

import (
	"net/http"
	"os"

	db "github.com/salty-max/go-comments/pkg/database"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

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

	// cmtService := comment.NewService(store)

	// cmtService.CreateComment(
	// 	context.Background(),
	// 	comment.Comment{
	// 		Slug:   "manual-test",
	// 		Body:   "Hello world!",
	// 		Author: "Max",
	// 	},
	// )
	// fmt.Println(cmtService.GetComment(
	// 	context.Background(),
	// 	"43e99d25-2139-4dd4-b099-efd23c923c97",
	// ))

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	if err := r.Run(":" + os.Getenv("PORT")); err != nil {
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
