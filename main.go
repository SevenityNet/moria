package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	go startup()

	<-ctx.Done()

	log.Println("Shutting down ...")
}

func startup() {
	initCache()
	initIO()
	initAuth()

	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	r.GET("/", getHello)
	r.POST("/media/:folder/upload", uploadFile)
	r.DELETE("/media/:folder/:file", deleteFile)
	r.GET("/media/:folder/:file", serveFile)
	r.GET("/auth/invite", invite)

	r.Run(":1980")
}

func getHello(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "I am alive! 👋",
	})
}
