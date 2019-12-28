package user_routes

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/ricecake/janus/model"
	"github.com/ricecake/janus/util"
	"github.com/spf13/viper"
)

func userActivate(c *gin.Context) {
	clientParams, marshErr := json.Marshal(map[string]string{
		"client_id": viper.GetString("identity.issuer_id"),
	})

	if marshErr != nil {
		c.Error(marshErr).SetType(gin.ErrorTypePrivate)
		c.AbortWithError(500, fmt.Errorf("System Error")).SetType(gin.ErrorTypePublic)
		return
	}

	body, renderErr := util.RenderHTMLTemplate("activate", util.TemplateContext{
		"ClientParams": string(clientParams),
	})

	if renderErr != nil {
		c.Error(renderErr).SetType(gin.ErrorTypePrivate)
		c.AbortWithError(500, fmt.Errorf("System Error")).SetType(gin.ErrorTypePublic)
		return
	}

	c.Data(200, "text/html", body)
}

type IdentityActivationArgs struct {
	PreferredName  string `form:"preferred_name"  json:"preferred_name"  binding:"omitempty"`
	Password       string `form:"password"        json:"password"        binding:"omitempty,min=8"`
	VerifyPassword string `form:"verify_password" json:"verify_password" binding:"omitempty,min=8"`
}

func userActivateApi(c *gin.Context) {
	identity, identErr := model.FindIdentityById(c.GetString("Identity"))
	if identErr != nil {
		c.Error(identErr).SetType(gin.ErrorTypePrivate)
		c.AbortWithError(401, fmt.Errorf("bad user")).SetType(gin.ErrorTypePublic)
		return
	}

	var activateArgs IdentityActivationArgs
	if err := c.ShouldBind(&activateArgs); err != nil {
		c.AbortWithError(400, err)
		return
	}

	if activateArgs.Password != activateArgs.VerifyPassword {
		c.AbortWithError(400, fmt.Errorf("passwords do not match")).SetType(gin.ErrorTypePublic)
		return
	}

	passErr := identity.SetPassword(activateArgs.Password)
	if passErr != nil {
		c.AbortWithError(400, passErr)
		return
	}

	identity.PreferredName = activateArgs.PreferredName
	identity.Active = true

	saveErr := identity.SaveChanges()
	if saveErr != nil {
		c.AbortWithError(400, saveErr)
		return
	}

	c.JSON(200, identity)
}
