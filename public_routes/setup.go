package public_routes

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

/*
/login
/logout
/token
/userinfo
/wellknown-endpoint
/auth/list
/auth/check
*/

func Configure(r *gin.RouterGroup) {
	log.Info("Configuring Public routes...")

	r.GET("/", defaultPage)
	r.GET("/.well-known/openid-configuration", discovery)

	r.GET("/login", loginPage)
	r.POST("/login", loginSubmit)

	r.GET("/logout", logoutPage)
	r.POST("/logout", logoutSubmit)

	r.GET("/signup", signupPage)
	r.POST("/signup", signupSubmit)

	r.GET("/publickeys", publicKeys)
	r.GET("/userinfo", userInfo)
	r.POST("/token", accessToken)

	r.GET("/check/auth/background", checkAuthBackground)
	r.GET("/check/auth/redirect", checkAuthRedirect)

	r.GET("/check/username", checkUsername)
	r.POST("/check/auth", checkAuth)
}
