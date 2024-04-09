package main

import (
	"context"
	"log"
	"moria/cache"
	"moria/config"
	"moria/rest"
	"moria/source"
	"os"
	"os/signal"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	config.Validate()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	go startMoria()

	<-ctx.Done()
	log.Println("Shutdown now")
}

func startMoria() {
	source.Initialize()
	cache.Initialize()

	r := gin.Default()

	if config.IsSecurityCORSEnabled() {
		r.Use(cors())
	}

	rest.RegisterRoutes(&r.RouterGroup)

	err := r.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func cors() gin.HandlerFunc {
	origin := config.GetSecurityCORSOrigin()
	allowMethodsArr := []string{"GET", "OPTIONS"}
	allowHeadersArr := []string{"Content-Type", "Baggage", "Accept"}

	if config.IsAPIEnabled() {
		allowMethodsArr = append(allowMethodsArr, "POST")

		allowHeadersArr = append(allowHeadersArr, config.GetSecurityAPIAuthHeader())
	}

	if config.IsFrontendEnabled() {
		allowMethodsArr = append(allowMethodsArr, "DELETE")

		if !config.IsAPIEnabled() {
			allowMethodsArr = append(allowMethodsArr, "POST")
		}

		if !contains(allowHeadersArr, "Authorization") {
			allowHeadersArr = append(allowHeadersArr, "Authorization")
		}
	}

	allowMethods := strings.Join(allowMethodsArr, ", ")
	allowHeaders := strings.Join(allowHeadersArr, ", ")

	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Methods", allowMethods)
		c.Writer.Header().Set("Access-Control-Allow-Headers", allowHeaders)

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func contains(arr []string, val string) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}

	return false
}
