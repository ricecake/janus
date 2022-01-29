package user_routes

import (
	"janus/model"

	"github.com/gin-gonic/gin"
)

func listUserDetails(c *gin.Context) {
	identity := c.GetString("Identity")
	user, err := model.FindIdentityById(identity)
	if err != nil {
		c.AbortWithStatusJSON(400, "bad user")
		return
	}

	c.JSON(200, user)
}

func updateUserDetails(c *gin.Context) {
	identity := c.GetString("Identity")
	user, err := model.FindIdentityById(identity)
	if err != nil {
		c.AbortWithStatusJSON(400, "bad user")
		return
	}

	var updated model.Identity
	bindErr := c.ShouldBind(&updated)
	if bindErr != nil {
		c.AbortWithStatusJSON(400, "Bad input")
		return
	}

	user.PreferredName = updated.PreferredName
	user.GivenName = updated.GivenName
	user.FamilyName = updated.FamilyName

	saveErr := user.SaveChanges()
	if saveErr != nil {
		c.AbortWithStatusJSON(500, "Error saving changes")
	}

	c.JSON(200, user)
}
