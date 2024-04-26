package main

import (
	"os"

	"github.com/allegro/bigcache/v3"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	AUTH_TOKEN string
)

func initAuth() {
	AUTH_TOKEN = os.Getenv("AUTH_TOKEN")
	if AUTH_TOKEN == "" {
		panic("no auth token env")
	}
}

func authenticateBackend(c *gin.Context) bool {
	token := c.GetHeader("Authorization")

	if token != AUTH_TOKEN {
		c.JSON(401, gin.H{
			"error": "Unauthorized",
		})

		return false
	}

	return true
}

func authenticateCodeOrToken(c *gin.Context) bool {
	codeOrToken := c.GetHeader("Authorization")

	if codeOrToken == "" {
		c.JSON(401, gin.H{
			"error": "Unauthorized",
		})

		return false
	}

	if codeOrToken == AUTH_TOKEN {
		return true
	}

	_, err := AUTHCACHE.Get(codeOrToken)
	if err == bigcache.ErrEntryNotFound {
		c.JSON(401, gin.H{
			"error": "Unauthorized",
		})
		return false
	} else if err != nil {
		panic(err)
	}

	AUTHCACHE.Delete(codeOrToken)

	return true
}

func getNewInviteCode() string {
	generatedToken := uuid.New()

	AUTHCACHE.Set(generatedToken.String(), []byte(generatedToken.String()))

	return generatedToken.String()
}
