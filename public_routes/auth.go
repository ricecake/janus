package public_routes

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/ricecake/osin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	kcutil "github.com/ricecake/karma_chameleon/util"
	"janus/model"
)

var (
	server   *osin.Server
	tokenGen model.TokenGenerator
)

func init() {
	sconfig := osin.NewServerConfig()
	sconfig.RedirectUriSeparator = "|"
	sconfig.AllowClientSecretInParams = true
	sconfig.RequirePKCEForPublicClients = true
	sconfig.AllowedAuthorizeTypes = osin.AllowedAuthorizeType{osin.CODE, osin.TOKEN}
	sconfig.AllowedAccessTypes = osin.AllowedAccessType{
		osin.AUTHORIZATION_CODE,
		osin.CLIENT_CREDENTIALS,
		osin.REFRESH_TOKEN,
		osin.PASSWORD,
		osin.IMPLICIT,
	}
	server = osin.NewServer(sconfig, model.NewDbStorage())

	server.AccessTokenGen = &tokenGen
	server.AuthorizeTokenGen = &tokenGen
}

func loginPage(c *gin.Context) {
	response := server.NewResponse()
	defer response.Close()

	ar := server.HandleAuthorizeRequest(response, c.Request)
	if ar == nil {
		if response.IsError && response.InternalError != nil {
			log.Printf("internal error: %v", response.InternalError)
		}
		osin.OutputJSON(response, c.Writer, c.Request)
		return
	}

	client, clientErr := model.FindClientById(ar.Client.GetId())
	if clientErr != nil {
		c.AbortWithError(400, fmt.Errorf("Client Not Found")).SetType(gin.ErrorTypePublic)
		return
	}

	prompt := c.DefaultQuery("prompt", "any")
	if prompt != "login" {
		res := attemptIdentifyUser(c, model.IdentificationRequest{
			Strategy: model.SESSION_TOKEN,
			Context:  &client.Context,
		})
		if res.Success {
			permitted := res.Success

			//TODO: this and the other one needs to account for what happens if this is the initial login to activate the user
			//	maybe hitting the zip link should just activate the user?  At that point, they're confirmed...
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

			permitted = permitted && allowed

			if permitted {
				ar.Authorized = permitted

				userDetails, err := establishSession(c, client.Context, *res)
				if err != nil {
					c.Error(err).SetType(gin.ErrorTypePrivate)
					c.AbortWithError(500, fmt.Errorf("System Error")).SetType(gin.ErrorTypePublic)
				}
				userDetails.Nonce = c.Query("nonce")
				ar.UserData = userDetails
				server.FinishAuthorizeRequest(response, c.Request, ar)

				if response.IsError && response.InternalError != nil {
					log.Printf("internal error: %v", response.InternalError)
				}
				osin.OutputJSON(response, c.Writer, c.Request)
				return
			}
		}
	}

	if prompt == "none" {
		ar.Authorized = false
		server.FinishAuthorizeRequest(response, c.Request, ar)
		osin.OutputJSON(response, c.Writer, c.Request)
		return
	}

	c.HTML(200, "template.html", gin.H{
		"CspNonce": c.GetString("CspNonce"),
		"preload": gin.H{
			"Name":     client.DisplayName,
			"RawQuery": c.Request.URL.RawQuery,
		},
	})
}

func loginSubmit(c *gin.Context) {
	response := server.NewResponse()
	defer response.Close()

	ar := server.HandleAuthorizeRequest(response, c.Request)
	if ar == nil {
		if response.IsError && response.InternalError != nil {
			log.Printf("internal error: %v", response.InternalError)
		}
		osin.OutputJSON(response, c.Writer, c.Request)
		return
	}

	client, clientErr := model.FindClientById(ar.Client.GetId())
	if clientErr != nil {
		c.AbortWithError(400, fmt.Errorf("Client Not Found")).SetType(gin.ErrorTypePublic)
		return
	}

	res := attemptIdentifyUser(c, model.IdentificationRequest{
		Strategy: model.NONE,
		Context:  &client.Context,
	})

	permitted := res.Success

	if permitted {
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

		permitted = permitted && allowed
	}

	ar.Authorized = permitted

	if permitted {
		userDetails, err := establishSession(c, client.Context, *res)
		if err != nil {
			c.Error(err).SetType(gin.ErrorTypePrivate)
			c.AbortWithError(500, fmt.Errorf("System Error")).SetType(gin.ErrorTypePublic)
			return
		}
		userDetails.Nonce = c.Query("nonce")
		ar.UserData = userDetails
	}

	server.FinishAuthorizeRequest(response, c.Request, ar)

	if response.IsError && response.InternalError != nil {
		log.Printf("internal error: %v", response.InternalError)
	}
	osin.OutputJSON(response, c.Writer, c.Request)
}

func signupPage(c *gin.Context) {
	client, clientErr := model.FindClientById(c.Query("client_id"))
	if clientErr != nil {
		c.AbortWithError(400, fmt.Errorf("Client Not Found")).SetType(gin.ErrorTypePublic)
		return
	}

	c.HTML(200, "template.html", gin.H{
		"CspNonce": c.GetString("CspNonce"),
		"preload": gin.H{
			"Name":     client.DisplayName,
			"RawQuery": c.Request.URL.RawQuery,
		},
	})
}

type SignupParams struct {
	PreferredName string `form:"preferred_name" json:"preferred_name"`
	Email         string `form:"email"          json:"email"`
}

func signupSubmit(c *gin.Context) {
	client, clientErr := model.FindClientById(viper.GetString("identity.issuer_id"))
	if clientErr != nil {
		c.Error(clientErr).SetType(gin.ErrorTypePrivate)
		c.AbortWithError(500, fmt.Errorf("System Error")).SetType(gin.ErrorTypePublic)
		return
	}

	user := &model.Identity{}

	var signupParams SignupParams
	bindErr := c.ShouldBind(&signupParams)
	if bindErr != nil {
		c.Error(bindErr).SetType(gin.ErrorTypePrivate)
		c.AbortWithError(500, fmt.Errorf("System Error")).SetType(gin.ErrorTypePublic)
		return
	}

	user.Email = signupParams.Email
	user.PreferredName = signupParams.PreferredName
	user.Active = true

	// TODO: add the concept of active and verified.  Active users can do user things.  Verified users can log in outside of idp context.  Premature?

	if err := model.CreateIdentity(user); err != nil {
		c.Error(err).SetType(gin.ErrorTypePrivate)
		c.AbortWithError(500, fmt.Errorf("System Error")).SetType(gin.ErrorTypePublic)
		return
	}

	_, err := establishSession(c, client.Context, model.IdentificationResult{Identity: user})
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePrivate)
		c.AbortWithStatusJSON(500, "system error")
		return
	}

	c.Status(204)
}

type SignupPasswordArgs struct {
	Password       string `form:"password"        json:"password"        binding:"omitempty,min=8"`
	VerifyPassword string `form:"verify_password" json:"verify_password" binding:"omitempty,min=8"`
}

func signupPassword(c *gin.Context) {
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

	var signupPassword SignupPasswordArgs
	if err := c.ShouldBind(&signupPassword); err != nil {
		c.AbortWithError(400, err)
		return
	}

	if signupPassword.Password != signupPassword.VerifyPassword {
		c.AbortWithError(400, fmt.Errorf("passwords do not match")).SetType(gin.ErrorTypePublic)
		return
	}

	passErr := res.Identity.SetPassword(signupPassword.Password)
	if passErr != nil {
		c.AbortWithError(400, passErr)
		return
	}

	c.Status(201)
}

func logoutPage(c *gin.Context) {
	for _, cookie := range c.Request.Cookies() {
		if cookieVal := cookie.Value; cookieVal != "" {
			var encData model.IDToken
			if decodeErr := kcutil.DecodeJWTOpen(cookieVal, &encData); decodeErr == nil {
				revokeErr := model.RevokeSessionToken(encData.TokenId)
				if revokeErr != nil {
					log.Error(revokeErr)
				}
			}

			clearSessionCookie(c, cookie.Name, cookie.Domain)
		}
	}
	c.String(200, "Logged out")
}

type AuthParams struct {
	Username *string `form:"username" json:"username"`
	Email    *string `form:"email" json:"email"`
	Password *string `form:"password" json:"password"`
	Totp     *string `form:"totp" json:"totp"`
}

func attemptIdentifyUser(c *gin.Context, authData model.IdentificationRequest) *model.IdentificationResult {
	if authData.Strategy == model.NONE || authData.Strategy == model.PASSWORD {
		var authParams AuthParams
		if c.ShouldBind(&authParams) == nil {
			// For RFC compliance, sometimes this field
			// must be "Username", even though it's email in the structure
			if authParams.Username != nil {
				authData.Strategy = model.PASSWORD
				authData.Email = authParams.Username
			} else if authParams.Email != nil {
				authData.Strategy = model.PASSWORD
				authData.Email = authParams.Email
			}
			if authParams.Password != nil {
				authData.Strategy = model.PASSWORD
				authData.Password = authParams.Password
			}
			if authParams.Totp != nil {
				authData.Strategy = model.PASSWORD
				authData.Totp = authParams.Totp
			}
		}
	}

	if authData.Strategy == model.NONE || authData.Strategy == model.SESSION_TOKEN {
		if authData.Context != nil {
			cookieName := fmt.Sprintf("janus.auth.session.%s", *authData.Context)
			cookies := []string{}
			for _, cookie := range c.Request.Cookies() {
				if cookie.Name == cookieName {
					if cookieVal := cookie.Value; cookieVal != "" {
						authData.Strategy = model.SESSION_TOKEN
						cookies = append(cookies, cookieVal)
					}
				}
			}
			authData.SessionToken = &cookies
		}
	}

	authResult := model.IdentifyFromCredentials(authData)
	if !authResult.Success {
		if authData.Context != nil {
			cookieName := fmt.Sprintf("janus.auth.session.%s", *authData.Context)
			for _, cookie := range c.Request.Cookies() {
				if cookie.Name == cookieName {
					clearSessionCookie(c, cookieName, cookie.Domain)
				}
			}
		}
	}

	return authResult
}
