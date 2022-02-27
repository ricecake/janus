package admin_routes

import (
	"janus/model"

	"github.com/gin-gonic/gin"
)

func listContexts(c *gin.Context) {
	ctx, err := model.ListContexts()
	if err != nil {
		c.AbortWithStatusJSON(500, "Internal error")
		return
	}
	c.JSON(200, ctx)
}

func updateContext(c *gin.Context) {
	var ctx model.Context
	if err := c.ShouldBind(&ctx); err != nil {
		c.AbortWithError(400, err)
		return
	}

	if err := ctx.SaveChanges(); err != nil {
		c.AbortWithError(400, err)
		return
	}

	c.JSON(200, ctx)
}

func createContext(c *gin.Context) {
	var ctx model.Context
	if err := c.ShouldBind(&ctx); err != nil {
		c.AbortWithError(400, err)
		return
	}

	if err := model.CreateContext(&ctx); err != nil {
		c.AbortWithError(400, err)
		return
	}

	c.JSON(200, ctx)

}
