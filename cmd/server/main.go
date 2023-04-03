package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Run() error {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.Run(":1664")

	return nil
}

func main() {
	if err := Run(); err != nil {
		fmt.Println(err)
	}
}
