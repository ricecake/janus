package http_middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/ricecake/janus/util"
	// "context"
)

func NewEnvMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("Reqid", util.CompactUUID())
		c.Next()
	}
}
