package rest

import (
	"moria/config"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/media/:category/:id", getImage)

	if config.IsAPIEnabled() {
		r.POST(config.GetAPIUploadEndpoint()+"/:category", postImage)
	}
}
