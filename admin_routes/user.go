package admin_routes

import (
	"janus/model"

	"github.com/gin-gonic/gin"
)

func listUsers(c *gin.Context) {
	ctx, err := model.ListIdentitys()
	if err != nil {
		c.AbortWithStatusJSON(500, "Internal error")
		return
	}
	c.JSON(200, ctx)
}
