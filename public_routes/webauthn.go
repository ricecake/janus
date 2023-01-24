package public_routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/ricecake/karma_chameleon/util"

	"janus/model"
)

var webAuthn *webauthn.WebAuthn

func setupWebauthn() {
	var err error
	webAuthn, err = webauthn.New(&webauthn.Config{
		RPDisplayName: viper.GetString("basic.name"),
		RPID:          viper.GetString("basic.domain"),
		RPOrigin:      viper.GetString("basic.site"),
	})

	if err != nil {
		log.Fatal("failed to create WebAuthn from config:", err)
	}
}

func registerStart(c *gin.Context) {
	idp, clientErr := model.FindClientById(viper.GetString("identity.issuer_id"))
	if clientErr != nil {
		log.Error("Error with own client?")
		c.Error(clientErr).SetType(gin.ErrorTypePrivate)
		c.AbortWithStatusJSON(500, "system error")
		return
	}
	res := attemptIdentifyUser(c, model.IdentificationRequest{
		Strategy: model.SESSION_TOKEN,
		Context:  &idp.Context,
	})
	if !res.Success {
		c.AbortWithError(res.FailureCode, fmt.Errorf(res.FailureReason))
		return
	}

	registerOptions := func(credCreationOpts *protocol.PublicKeyCredentialCreationOptions) {
		credCreationOpts.CredentialExcludeList = res.Identity.CredentialExcludeList()
	}

	options, sessionData, err := webAuthn.BeginRegistration(
		res.Identity,
		registerOptions,
	)

	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePrivate)
		c.AbortWithError(400, fmt.Errorf("Input Error")).SetType(gin.ErrorTypePublic)
		return
	}

	sessionCookie, err := util.EncodeJWTClose(sessionData, viper.GetString("security.passphrase"))
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "janus.webauthn.register.state",
		Value:    sessionCookie,
		Path:     "/",
		Secure:   !viper.GetBool("development.insecure"),
		HttpOnly: true,
		MaxAge:   600,
	})

	c.JSON(200, options)
}

func registerFinish(c *gin.Context) {
	idp, clientErr := model.FindClientById(viper.GetString("identity.issuer_id"))
	if clientErr != nil {
		log.Error("Error with own client?")
		c.Error(clientErr).SetType(gin.ErrorTypePrivate)
		c.AbortWithStatusJSON(500, "system error")
		return
	}
	res := attemptIdentifyUser(c, model.IdentificationRequest{
		Strategy: model.SESSION_TOKEN,
		Context:  &idp.Context,
	})
	if !res.Success {
		c.AbortWithError(res.FailureCode, fmt.Errorf(res.FailureReason))
		return
	}

	var sessionData webauthn.SessionData
	cookieName := fmt.Sprintf("janus.webauthn.register.state")
	for _, cookie := range c.Request.Cookies() {
		if cookie.Name == cookieName {
			if cookieVal := cookie.Value; cookieVal != "" {
				if err := util.DecodeJWTClose(cookieVal, viper.GetString("security.passphrase"), &sessionData); err != nil {
					c.AbortWithError(400, err)
					return
				}
			}
			break
		}
	}
	clearWebauthnCookie(c, "register")

	credential, err := webAuthn.FinishRegistration(res.Identity, sessionData, c.Request)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	name := c.Query("name")

	if err := res.Identity.AddWebauthnCredential(name, credential); err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.Status(204)
}

func loginStart(c *gin.Context) {
	email := c.Param("email")
	client_id := c.Param("client_id")

	ident, err := model.FindIdentityByEmail(email)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	client, clientErr := model.FindClientById(client_id)
	if clientErr != nil {
		c.AbortWithError(400, fmt.Errorf("Client Not Found")).SetType(gin.ErrorTypePublic)
		return
	}

	allowed, err := model.AclCheck(model.AclCheckRequest{
		Identity: ident.Code,
		Context:  client.Context,
		Action:   client.ClientId,
	})

	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePrivate)
		c.AbortWithError(500, fmt.Errorf("System Error")).SetType(gin.ErrorTypePublic)
		return
	}

	if !allowed {
		c.AbortWithStatus(401)
		return
	}

	options, sessionData, err := webAuthn.BeginLogin(ident)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	sessionCookie, err := util.EncodeJWTClose(sessionData, viper.GetString("security.passphrase"))
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "janus.webauthn.login.state",
		Value:    sessionCookie,
		Path:     "/",
		Secure:   !viper.GetBool("development.insecure"),
		HttpOnly: true,
		MaxAge:   600,
	})

	c.JSON(200, options)
}

func loginFinish(c *gin.Context) {
	email := c.Param("email")
	client_id := c.Param("client_id")

	ident, err := model.FindIdentityByEmail(email)
	if err != nil {
		c.AbortWithError(400, fmt.Errorf("Bad input"))
		return
	}

	var sessionData webauthn.SessionData
	cookieName := fmt.Sprintf("janus.webauthn.login.state")
	for _, cookie := range c.Request.Cookies() {
		if cookie.Name == cookieName {
			if cookieVal := cookie.Value; cookieVal != "" {
				if err := util.DecodeJWTClose(cookieVal, viper.GetString("security.passphrase"), &sessionData); err != nil {
					c.AbortWithError(400, err)
					return
				}
			}
			break
		}
	}
	clearWebauthnCookie(c, "login")

	credential, err := webAuthn.FinishLogin(ident, sessionData, c.Request)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	client, clientErr := model.FindClientById(client_id)
	if clientErr != nil {
		c.AbortWithError(400, fmt.Errorf("Client Not Found")).SetType(gin.ErrorTypePublic)
		return
	}

	res := attemptIdentifyUser(c, model.IdentificationRequest{
		Strategy:   model.WEBAUTHN,
		Credential: credential,
		Context:    &client.Context,
	})

	if res.Success {
		allowed, err := model.AclCheck(model.AclCheckRequest{
			Identity: res.Identity.Code,
			Context:  client.Context,
			Action:   client.ClientId,
		})
		if err != nil {
			c.Error(err).SetType(gin.ErrorTypePrivate)
			c.AbortWithError(500, fmt.Errorf("System Error")).SetType(gin.ErrorTypePublic)
			return
		}

		// TODO: increment webauthn credential auth count
		// check 'credential.Authenticator.CloneWarning'

		if allowed {
			_, err = establishSession(c, client.Context, *res)
			if err != nil {
				c.Error(err).SetType(gin.ErrorTypePrivate)
				c.AbortWithStatusJSON(500, "system error")
			}
			c.Status(204)
			return
		}
	}

	c.Status(401)
}

func clearWebauthnCookie(c *gin.Context, cookieName string) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     fmt.Sprintf("janus.webauthn.%s.state", cookieName),
		Value:    "",
		Path:     "/",
		Secure:   !viper.GetBool("development.insecure"),
		HttpOnly: true,
		MaxAge:   -1,
	})
}
