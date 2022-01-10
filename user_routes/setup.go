package user_routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ricecake/karma_chameleon/http_middleware"
	log "github.com/sirupsen/logrus"
	"janus/model"
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
}
