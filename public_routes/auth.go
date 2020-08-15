package public_routes

import (
	"fmt"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/ricecake/osin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/ricecake/janus/model"
	"github.com/ricecake/karma_chameleon/util"
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
			ar.Authorized = res.Success
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

	if prompt == "none" {
		ar.Authorized = false
		server.FinishAuthorizeRequest(response, c.Request, ar)
		osin.OutputJSON(response, c.Writer, c.Request)
		return
	}

	body, renderErr := util.RenderHTMLTemplate("login", util.TemplateContext{
		"Name":     client.DisplayName,
		"Param":    c.Request.URL.Query(),
		"RawQuery": c.Request.URL.RawQuery,
		"CspNonce": c.GetString("CspNonce"),
	})

	if renderErr != nil {
		c.Error(renderErr).SetType(gin.ErrorTypePrivate)
		c.AbortWithError(500, fmt.Errorf("System Error")).SetType(gin.ErrorTypePublic)
		return
	}

	c.Data(200, "text/html", body)
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
	body, renderErr := util.RenderHTMLTemplate("signup", util.TemplateContext{
		"Name":     client.DisplayName,
		"Param":    c.Request.URL.Query(),
		"RawQuery": c.Request.URL.RawQuery,
		"CspNonce": c.GetString("CspNonce"),
	})

	if renderErr != nil {
		c.Error(renderErr).SetType(gin.ErrorTypePrivate)
		c.AbortWithError(500, fmt.Errorf("System Error")).SetType(gin.ErrorTypePublic)
		return
	}

	c.Data(200, "text/html", body)
}

type SignupParams struct {
	PreferredName string `form:"preferred_name" json:"preferred_name"`
	Email         string `form:"email"          json:"email"`
}

func signupSubmit(c *gin.Context) {
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

	if err := model.CreateIdentity(user); err != nil {
		c.Error(err).SetType(gin.ErrorTypePrivate)
		c.AbortWithError(500, fmt.Errorf("System Error")).SetType(gin.ErrorTypePublic)
		return
	}

	var code string
	if referer, err := url.Parse(c.Request.Header.Get("Referer")); err == nil {
		referer.Path = "/login"
		stashCode, stashErr := model.StashTTL(&map[string]string{
			"Redirect": referer.String(),
		}, 86400)
		if stashErr != nil {
			c.Error(stashErr).SetType(gin.ErrorTypePrivate)
			c.AbortWithError(500, fmt.Errorf("System Error")).SetType(gin.ErrorTypePublic)
			return
		}
		code = stashCode
	}

	zipCode := model.ZipCode{
		Identity:    user.Code,
		Client:      viper.GetString("identity.issuer_id"),
		TTL:         86400, // One day
		RedirectUri: "/profile/activate",
		Params: map[string]string{
			"code": code,
		},
	}
	if zipErr := zipCode.Save(); zipErr != nil {
		c.Error(zipErr).SetType(gin.ErrorTypePrivate)
		c.AbortWithError(500, fmt.Errorf("System Error")).SetType(gin.ErrorTypePublic)
		return
	}

	emailErr := util.SendMail(user.PreferredName, user.Email, "activation", util.TemplateContext{
		"Name":  user.PreferredName,
		"Email": user.Email,
		"Code":  zipCode.Code,
	})

	if emailErr != nil {
		c.Error(emailErr).SetType(gin.ErrorTypePrivate)
		c.AbortWithError(500, fmt.Errorf("System Error")).SetType(gin.ErrorTypePublic)
		return
	}

	c.Status(204)
}
func logoutPage(c *gin.Context)   {}
func logoutSubmit(c *gin.Context) {}

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
			for _, cookie := range c.Request.Cookies() {
				if cookie.Name == cookieName {
					if cookieVal := cookie.Value; cookieVal != "" {
						authData.Strategy = model.SESSION_TOKEN
						authData.SessionToken = &cookieVal
					}
					break
				}
			}
		}
	}

	return model.IdentifyFromCredentials(authData)
}
