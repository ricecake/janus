package user_routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ricecake/janus/http_middleware"
	log "github.com/sirupsen/logrus"
)

/*
/profile/user
/profile/other stuff
/profile/activate
/profile/reset/password
*/

//Need some middleware for ensuring that the user is authenticated

func Configure(r *gin.RouterGroup) {
	log.Info("Configuring user routes...")

	r.GET("/activate", userActivate)

	apiGroup := r.Group("/api")
	apiGroup.Use(http_middleware.NewAuthMiddleware())

	apiGroup.POST("/activate", userActivateApi)
}
