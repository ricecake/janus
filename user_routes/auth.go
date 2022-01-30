package user_routes

import (
	"fmt"
	"janus/model"

	"github.com/gin-gonic/gin"
)

type ChangePassword struct {
	Password       string `form:"password"        json:"password"        binding:"omitempty,min=8"`
	VerifyPassword string `form:"verify_password" json:"verify_password" binding:"omitempty,min=8"`
}

func changePassword(c *gin.Context) {
	identity := c.GetString("Identity")
	user, err := model.FindIdentityById(identity)
	if err != nil {
		c.AbortWithStatusJSON(400, "bad user")
		return
	}

	var changePassword ChangePassword
	if err := c.ShouldBind(&changePassword); err != nil {
		c.AbortWithError(400, err)
		return
	}

	if changePassword.Password != changePassword.VerifyPassword {
		c.AbortWithError(400, fmt.Errorf("passwords do not match")).SetType(gin.ErrorTypePublic)
		return
	}

	passErr := user.SetPassword(changePassword.Password)
	if passErr != nil {
		c.AbortWithError(400, passErr)
		return
	}

	c.Status(204)
}

func listAuthenticators(c *gin.Context) {
	identity := c.GetString("Identity")
	user, err := model.FindIdentityById(identity)
	if err != nil {
		c.AbortWithStatusJSON(400, "bad user")
		return
	}

	c.JSON(200, model.AuthenticatorsForIdentity(user.Code))
}

func deleteAuthenticator(c *gin.Context) {
	identity := c.GetString("Identity")
	user, err := model.FindIdentityById(identity)
	if err != nil {
		c.AbortWithStatusJSON(400, "bad user")
		return
	}

	model.RemoveAuthenticator(user.Code, c.Query("name"))
	c.Status(204)
}
