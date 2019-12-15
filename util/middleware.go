package util

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func NewAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Info("Starting auth")
		headerParts := strings.SplitN(c.GetHeader("Authorization"), " ", 2)
		if len(headerParts) != 2 || strings.ToLower(headerParts[0]) != "bearer" {
			log.Error("Malformed header")
			c.AbortWithError(401, errors.New("Invalid authorization")).SetType(gin.ErrorTypePublic)
			return
		}

		var jwtParts map[string]interface{}

		decErr := DecodeJWTOpen(headerParts[1], &jwtParts)
		if decErr != nil {
			log.Error("Bad token: ", decErr)
			c.AbortWithError(401, errors.New("Invalid authorization")).SetType(gin.ErrorTypePublic)
			return
		}

		c.Set("ValidAuth", true)
		c.Set("Identity", jwtParts["sub"].(string))
		c.Next()
	}
}
