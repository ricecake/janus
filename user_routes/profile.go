package user_routes

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/ricecake/janus/util"
)

func userActivate(c *gin.Context) {
	body, renderErr := util.RenderTemplate("activate", util.TemplateContext{})

	if renderErr != nil {
		c.Error(renderErr).SetType(gin.ErrorTypePrivate)
		c.AbortWithError(500, fmt.Errorf("System Error")).SetType(gin.ErrorTypePublic)
		return
	}

	c.Data(200, "text/html", body)
}
