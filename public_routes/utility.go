package public_routes

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/ricecake/janus/model"
	"github.com/ricecake/janus/util"
)

func defaultPage(c *gin.Context)   {}
func checkUsername(c *gin.Context) {}

func checkAuth(c *gin.Context) {
	c.Header("Content-Type", "application/json; charset=utf-8")

	client, clientErr := model.FindClientById(c.Query("client_id"))
	if clientErr != nil {
		c.AbortWithStatusJSON(400, "Client Not Found")
		return
	}

	res := attemptIdentifyUser(c, model.IdentificationRequest{
		Strategy: model.NONE,
		Context:  &client.Context,
	})

	if res.Success {
		_, err := establishSession(c, client.Context, *res)
		if err != nil {
			c.Error(err).SetType(gin.ErrorTypePrivate)
			c.AbortWithStatusJSON(500, "system error")
		}
		c.Status(204)
		return
	}

	c.AbortWithStatusJSON(res.FailureCode, res.FailureReason)
}

// TODO this should just return the full redirect url in the location header
// That way it can stash much more away easier, and doesn't
// need to rely on the applicaiton to do the redirect stuff.
func checkAuthBackground(c *gin.Context) {
	c.Header("Content-Type", "application/json; charset=utf-8")

	client, clientErr := model.FindClientById(c.GetHeader("X-Client-Id"))
	if clientErr != nil {
		c.AbortWithStatusJSON(400, "Client Not Found")
		return
	}

	res := attemptIdentifyUser(c, model.IdentificationRequest{
		Strategy: model.SESSION_TOKEN,
		Context:  &client.Context,
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

		if allowed {
			c.Status(204)
			return
		}
	}

	redirect := c.GetHeader("X-Auth-Redirect")
	if redirect == "" {
		redirect = client.BaseUri
	}

	stashCode, err := model.Stash(map[string]string{"redirect": redirect})
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePrivate)
		c.AbortWithStatusJSON(500, "system error")
		return
	}

	idpClient, clientErr := model.FindClientById(viper.GetString("identity.issuer_id"))
	if clientErr != nil {
		c.Error(clientErr).SetType(gin.ErrorTypePrivate)
		c.AbortWithStatusJSON(500, "system error")
		return
	}

	redirectBase, err := url.Parse(idpClient.BaseUri)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePrivate)
		c.AbortWithStatusJSON(500, "system error")
		return
	}

	redirectBase.Path = "/check/auth/redirect"

	queryParams := url.Values{
		"scope":         []string{"openid"},
		"response_type": []string{"code"},
		"state":         []string{stashCode},
		"client_id":     []string{client.ClientId},
		"redirect_uri":  []string{redirectBase.String()},
	}
	redirectBase.Path = "/login"
	redirectBase.RawQuery = queryParams.Encode()

	c.Header("X-Auth-State", stashCode)
	c.Header("X-Auth-Scope", "openid")
	c.Header("X-Redirect-Location", redirectBase.String())

	c.AbortWithStatusJSON(401, res.FailureReason)
}

func checkAuthRedirect(c *gin.Context) {
	stateVar := c.Query("state")
	var data map[string]string
	if err := model.StashFetch(stateVar, &data); err != nil {
		c.Error(err).SetType(gin.ErrorTypePrivate)
		c.AbortWithStatusJSON(500, "system error")
		return
	}

	if errMsg := c.Query("error"); errMsg != "" {
		c.Error(fmt.Errorf(errMsg)).SetType(gin.ErrorTypePublic)
		c.String(401, c.Query("error_description"))
		c.Abort()
		return
	}

	var encData model.AuthCodeData
	if err := util.DecodeJWTClose(c.Query("code"), viper.GetString("security.passphrase"), &encData); err != nil {
		c.Error(err).SetType(gin.ErrorTypePrivate)
		c.AbortWithStatusJSON(500, "system error")
	}

	if encData.State != stateVar {
		c.AbortWithStatusJSON(500, "system error")
		return
	}
	//TODO: validate that the code isn't expired.
	// Need to make a generic function/interface for doing that, since it happens a lot

	client, clientErr := model.FindClientById(encData.ClientId)
	if clientErr != nil {
		c.AbortWithStatusJSON(400, "Client Not Found")
		return
	}

	res := attemptIdentifyUser(c, model.IdentificationRequest{
		Strategy: model.SESSION_TOKEN,
		Context:  &client.Context,
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

		if !allowed {
			c.String(401, "Access to this system is not allowed")
			c.Abort()
			return
		}

		idp, clientErr := model.FindClientById(viper.GetString("identity.issuer_id"))
		if clientErr != nil {
			log.Error("Error with own client?")
			c.Error(clientErr).SetType(gin.ErrorTypePrivate)
			c.AbortWithStatusJSON(500, "system error")
			return
		}

		user := res.Identity

		scopes := make(map[string]bool)
		for _, s := range strings.Fields(c.Query("scope")) {
			scopes[s] = true
		}

		token := user.IdentityToken(scopes)
		token.ClientID = idp.ClientId
		token.Context = client.Context
		encToken, err := util.EncodeJWTOpen(token)
		if err != nil {
			c.Error(err).SetType(gin.ErrorTypePrivate)
			c.AbortWithStatusJSON(500, "system error")
			return
		}

		redirectBase, err := url.Parse(data["redirect"])
		if err != nil {
			c.Error(err).SetType(gin.ErrorTypePrivate)
			c.AbortWithStatusJSON(500, "system error")
			return
		}

		idpBase, err := url.Parse(viper.GetString("identity.issuer"))
		if err != nil {
			c.Error(err).SetType(gin.ErrorTypePrivate)
			c.AbortWithStatusJSON(500, "system error")
			return
		}

		cookieName := fmt.Sprintf("janus.auth.session.%s", client.Context)
		http.SetCookie(c.Writer, &http.Cookie{
			Domain:   util.TrunkUrlFragment([]string{redirectBase.Host, idpBase.Host}),
			Name:     cookieName,
			Value:    encToken,
			Path:     "/",
			Secure:   !viper.GetBool("development.insecure"),
			HttpOnly: true,
		})

		c.Redirect(302, data["redirect"])
	}
}

func processZipCode(c *gin.Context) {
	code := c.Param("code")

	res := attemptIdentifyUser(c, model.IdentificationRequest{
		Strategy: model.ZIPCODE,
		ZipCode:  &code,
	})

	if res.Success {
		idp, clientErr := model.FindClientById(res.ZipCode.Client)
		if clientErr != nil {
			c.Error(clientErr).SetType(gin.ErrorTypePrivate)
			c.AbortWithStatusJSON(500, "system error")
			return
		}
		_, err := establishSession(c, idp.Context, *res)
		if err != nil {
			c.Error(clientErr).SetType(gin.ErrorTypePrivate)
			c.AbortWithStatusJSON(500, "system error")
			return
		}
		c.Redirect(302, res.ZipCode.RedirectUri)
		return
	}

	c.AbortWithStatusJSON(401, res.FailureReason)
}

func establishSession(c *gin.Context, context string, identData model.IdentificationResult) (*model.UserAuthDetails, error) {
	client, clientErr := model.FindClientById(viper.GetString("identity.issuer_id"))
	if clientErr != nil {
		log.Error("Error with own client?")
		return nil, clientErr
	}

	user := identData.Identity

	var sessionCode string
	if identData.Session == nil {
		token := user.IdentityToken(map[string]bool{})

		token.ClientID = client.ClientId
		token.Context = context

		encToken, err := util.EncodeJWTOpen(token)
		if err != nil {
			return nil, err
		}

		sessionToken := model.SessionToken{
			Code:      token.TokenId,
			Identity:  user.Code,
			UserAgent: c.Request.UserAgent(),
			IpAddress: c.ClientIP(),
			CreatedAt: time.Now(),
			ExpiresIn: int(token.Expiration - token.IssuedAt),
		}
		if err := model.CreateSessionToken(&sessionToken); err != nil {
			return nil, err
		}

		sessionCode = sessionToken.Code

		cookieName := fmt.Sprintf("janus.auth.session.%s", context)
		for _, cookie := range c.Request.Cookies() {
			if cookie.Name == cookieName {
				if cookieVal := cookie.Value; cookieVal != "" {
					var encData model.IDToken
					if err := util.DecodeJWTOpen(cookieVal, &encData); err == nil {
						if err := model.ReplaceSessionToken(encData.TokenId, sessionCode); err != nil {
							log.Error(err)
						}
					}
				}
				break
			}
		}

		http.SetCookie(c.Writer, &http.Cookie{
			Name:     cookieName,
			Value:    encToken,
			Path:     "/",
			Secure:   !viper.GetBool("development.insecure"),
			HttpOnly: true,
		})
	} else {
		// TODO: make sure that the session we have is valid, don't just trust the token
		sessionCode = *identData.Session
	}

	accessContext := model.AccessContext{
		Session:   sessionCode,
		Client:    client.ClientId,
		CreatedAt: time.Now(),
	}
	if err := model.EnsureAccessContext(&accessContext); err != nil {
		return nil, err
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "janus.user.email",
		Value:    user.Email,
		Path:     "/",
		Secure:   !viper.GetBool("development.insecure"),
		HttpOnly: false,
	})

	return &model.UserAuthDetails{
		Code:    user.Code,
		Browser: sessionCode,
		Context: accessContext.Code,
	}, nil
}
