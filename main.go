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

	go videoEncodingConsumer()

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	r.Use(cors())

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

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		corsAllowedOrigin := os.Getenv("CORS_ALLOWED_ORIGINS")

		if corsAllowedOrigin == "*" || corsAllowedOrigin == "" {
			requestDomain := c.Request.Header.Get("Origin")
			c.Writer.Header().Set("Access-Control-Allow-Origin", requestDomain)
		} else {
			c.Writer.Header().Set("Access-Control-Allow-Origin", corsAllowedOrigin)
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, PATCH, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Baggage, Accept, Sentry-Trace")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Authorization, Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	}
}
