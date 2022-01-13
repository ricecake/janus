package user_routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"janus/model"
)

func listUserApps(c *gin.Context) {
	if !c.GetBool("ValidAuth") {
		c.AbortWithStatusJSON(401, "bad user")
		return
	}
	identity := c.GetString("Identity")
	data, err := model.IdentityAllowedClients(identity)
	if err != nil {
		fmt.Printf(err.Error())
		c.AbortWithStatusJSON(500, "Error")
		return
	}
	c.JSON(200, data)
}
