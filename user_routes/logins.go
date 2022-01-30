package user_routes

import (
	"fmt"
	"janus/model"

	"github.com/gin-gonic/gin"
)

func listUserLogins(c *gin.Context) {
	if !c.GetBool("ValidAuth") {
		c.AbortWithStatusJSON(401, "bad user")
		return
	}
	identity := c.GetString("Identity")
	data, err := model.IdentityLogins(identity)
	if err != nil {
		fmt.Printf(err.Error())
		c.AbortWithStatusJSON(500, "Error")
		return
	}

	var session string
	var login string

	token, exists := c.Get("Token")
	if exists {
		coerced, ok := token.(model.AccessToken)
		if ok {
			session = coerced.Browser
			login = coerced.ContextCode
		}
	}

	c.JSON(200, gin.H{
		"Browser": session,
		"Access":  login,
		"Logins":  data,
	})
}
