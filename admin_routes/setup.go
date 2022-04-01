package admin_routes

/*
/basically-crud-support-for-everything
*/
import (
	"janus/model"

	"github.com/gin-gonic/gin"
	"github.com/ricecake/karma_chameleon/http_middleware"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Configure(r *gin.RouterGroup) {
	log.Info("Configuring admin routes...")

	apiGroup := r.Group("/api")
	apiGroup.Use(http_middleware.NewAuthMiddleware(model.NewLocalVerifierCache()))
	apiGroup.Use(func(c *gin.Context) {
		// This should be able to use data straight from the token, once that gets cleaned up
		// to use a more generic token interface that can handle role and action fields.
		// Also, should really be setting up a system to check the api path and method against specific actions,
		// so that things can be slightly more granular.
		// Not quite there yet, but a solid TODO.
		identity := c.GetString("Identity")
		user, err := model.FindIdentityById(identity)
		if err != nil {
			c.AbortWithStatusJSON(400, "bad user")
			return
		}

		idpClient, clientErr := model.FindClientById(viper.GetString("identity.issuer_id"))
		if clientErr != nil {
			c.Error(clientErr).SetType(gin.ErrorTypePrivate)
			c.AbortWithStatusJSON(500, "system error")
			return
		}

		allowed, err := model.HasContextRole(idpClient.Context, user.Code, "Admin")
		if err != nil {
			c.Error(err).SetType(gin.ErrorTypePrivate)
			c.AbortWithStatusJSON(500, "system error")
			return

		}

		if !allowed {
			c.AbortWithStatusJSON(403, "Need administrator access")
			return
		}

	})

	apiGroup.GET("/context", listContexts)
	apiGroup.PATCH("/context", updateContext)
	apiGroup.POST("/context", createContext)

	apiGroup.GET("/client", listClients)
	apiGroup.PATCH("/client", updateClient)
	apiGroup.POST("/client", createClient)

	apiGroup.GET("/user", listUsers)

	apiGroup.GET("/role", listRoles)

	apiGroup.GET("/action", listActions)

	apiGroup.GET("/group", listGroups)
}
