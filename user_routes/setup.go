package user_routes

import (
	"janus/model"

	"github.com/gin-gonic/gin"
	"github.com/ricecake/karma_chameleon/http_middleware"
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

	apiGroup := r.Group("/api")
	apiGroup.Use(http_middleware.NewAuthMiddleware(model.NewLocalVerifierCache()))

	apiGroup.GET("/applist", listUserApps)

	apiGroup.GET("/detail", listUserDetails)
	apiGroup.POST("/detail", updateUserDetails)

	apiGroup.GET("/login", listUserLogins)
	apiGroup.DELETE("/login", clearAccess)
	apiGroup.DELETE("/login/session", clearSession)

	apiGroup.GET("/authenticator", listAuthenticators)
	apiGroup.DELETE("/authenticator", deleteAuthenticator)

	apiGroup.POST("/password", changePassword)
}
