package admin_routes

import (
	"janus/model"

	"github.com/gin-gonic/gin"
)

func listRoles(c *gin.Context) {
	ctx, err := model.ListRoles()
	if err != nil {
		c.AbortWithStatusJSON(500, "Internal error")
		return
	}
	c.JSON(200, ctx)
}
