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

	r.POST("/login/link", sendLoginLink)

	r.GET("/zip/:code", processZipCode)

	r.DELETE("/sessions", logout)

	r.GET("/signup", signupPage)
	r.POST("/signup", signupSubmit)

	r.POST("/signup/password", signupPassword)

	r.GET("/publickeys", publicKeys)
	r.GET("/userinfo", userInfo)
	r.POST("/token", accessToken)

	r.GET("/check/auth/background", checkAuthBackground)
	r.GET("/check/auth/redirect", checkAuthRedirect)

	r.GET("/check/username", checkUsername)
	r.POST("/check/authenticators", authDetails)
	r.POST("/check/auth", checkAuth)
	r.GET("/revocation", listRevocation)

	setupWebauthn()
	r.POST("/webauthn/register/start", registerStart)
	r.POST("/webauthn/register/finish", registerFinish)
	r.POST("/webauthn/login/start/:client_id/:email", loginStart)
	r.POST("/webauthn/login/finish/:client_id/:email", loginFinish)
}
