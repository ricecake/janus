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

	kcutil "github.com/ricecake/karma_chameleon/util"
	"janus/model"
	"janus/util"
)

func defaultPage(c *gin.Context)   {}
func checkUsername(c *gin.Context) {}

type AuthDetailParams struct {
	Email string `form:"email" json:"email"`
}

func authDetails(c *gin.Context) {
	var authParams AuthDetailParams
	if err := c.ShouldBind(&authParams); err != nil {
		c.Error(err).SetType(gin.ErrorTypePrivate)
		c.AbortWithStatusJSON(400, "Bad request")
		return
	}

	ident, err := model.FindIdentityByEmail(authParams.Email)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePrivate)
		c.AbortWithStatusJSON(400, "Bad request")
		return
	}

	methods, err := ident.ValidAuthMethods()
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePrivate)
		c.AbortWithStatusJSON(400, "Bad request")
		return
	}

	c.JSON(200, methods)
}

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
			c.Header("X-Identity-Code", res.Identity.Code)
			c.Header("X-Identity-Email", res.Identity.Email)
			c.Header("X-Identity-Name", res.Identity.PreferredName)
			return
		}
	}

	log.Errorf("Failed to auth: %s", res.FailureReason)

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
	if err := kcutil.DecodeJWTClose(c.Query("code"), viper.GetString("security.passphrase"), &encData); err != nil {
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

	trunkDomain := kcutil.TrunkUrlFragment([]string{redirectBase.Host, idpBase.Host})
	cookieName := fmt.Sprintf("janus.auth.session.%s", client.Context)

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
		encToken, err := kcutil.EncodeJWTOpen(token)
		if err != nil {
			c.Error(err).SetType(gin.ErrorTypePrivate)
			c.AbortWithStatusJSON(500, "system error")
			return
		}

		http.SetCookie(c.Writer, &http.Cookie{
			Domain:   trunkDomain,
			Name:     cookieName,
			Value:    encToken,
			Path:     "/",
			Secure:   !viper.GetBool("development.insecure"),
			HttpOnly: true,
			MaxAge:   int(time.Until(time.Unix(token.Expiration, 0)).Seconds()),
		})

		c.Redirect(302, data["redirect"])
	} else {
		clearSessionCookie(c, cookieName, trunkDomain)
	}
}

type LoginLinkRequest struct {
	Email  string `form:"email" json:"email"`
	Client string `form:"client" json:"client"`
	Url    string `form:"url" json:"url"`
	State  string `form:"state" json:"state"`
}

func sendLoginLink(c *gin.Context) {
	var linkRequest LoginLinkRequest
	if err := c.ShouldBind(&linkRequest); err != nil {
		c.Error(err).SetType(gin.ErrorTypePrivate)
		c.AbortWithStatusJSON(400, "Bad request")
		return
	}

	ident, err := model.FindIdentityByEmail(linkRequest.Email)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePrivate)
		c.AbortWithStatusJSON(400, "Bad email")
		return
	}

	client, clientErr := model.FindClientById(linkRequest.Client)
	if clientErr != nil {
		c.Error(clientErr).SetType(gin.ErrorTypePrivate)
		c.AbortWithStatusJSON(400, "bad client")
		return
	}

	redirect, err := url.Parse(linkRequest.Url)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePrivate)
		c.AbortWithStatusJSON(400, "bad url")
		return
	}

	idpBase, err := url.Parse(viper.GetString("identity.issuer"))
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePrivate)
		c.AbortWithStatusJSON(500, "system error")
		return
	}

	params := map[string]string{
		"state": linkRequest.State,
	}

	// Redirects can happen to the client base uri, or into the idp, to facilitate background auth flow.
	if redirect.Host == idpBase.Host {
		var data map[string]string
		if err := model.StashFetch(linkRequest.State, &data); err != nil {
			c.Error(err).SetType(gin.ErrorTypePrivate)
			c.AbortWithStatusJSON(400, "Bad state")
			return
		}
		redirect, err = url.Parse(data["redirect"])
		if err != nil {
			c.Error(err).SetType(gin.ErrorTypePrivate)
			c.AbortWithStatusJSON(400, "Bad state")
			return
		}

		params = map[string]string{}
	}

	zipCode := model.ZipCode{
		Identity:    ident.Code,
		Client:      client.ClientId,
		RedirectUri: redirect.String(),
		TTL:         300, // five minutes
		Params:      params,
	}

	if zipErr := zipCode.Save(); zipErr != nil {
		c.Error(zipErr).SetType(gin.ErrorTypePrivate)
		c.AbortWithStatusJSON(500, "system error")
		return
	}

	emailErr := util.SendMail(ident.PreferredName, ident.Email, "login", util.TemplateContext{
		"Code":        zipCode.Code,
		"Client_Name": client.DisplayName,
	})
	if emailErr != nil {
		c.Error(emailErr).SetType(gin.ErrorTypePrivate)
		c.AbortWithStatusJSON(500, "system error")
		return
	}

	c.Status(204)
}

func processZipCode(c *gin.Context) {
	code := c.Param("code")

	res := attemptIdentifyUser(c, model.IdentificationRequest{
		Strategy: model.ZIPCODE,
		ZipCode:  &code,
	})

	if res.Success {
		if res.ZipCode.Signup {
			// If it's a signup, activate the user since they just verified
			// This is also going to be used for inviting a new user.  User is inactive until they click the link.
			identity := res.Identity
			identity.Active = true

			saveErr := identity.SaveChanges()
			if saveErr != nil {
				c.AbortWithError(400, saveErr)
				return
			}

		}
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

func listRevocation(c *gin.Context) {
	revoked, revErr := model.ListRevocations()
	if revErr != nil {
		c.Error(revErr).SetType(gin.ErrorTypePrivate)
		c.AbortWithStatus(500)
	}

	c.JSON(200, revoked)
}

func establishSession(c *gin.Context, context string, identData model.IdentificationResult) (*model.UserAuthDetails, error) {
	client, clientErr := model.FindClientById(viper.GetString("identity.issuer_id"))
	if clientErr != nil {
		log.Error("Error with own client?")
		return nil, clientErr
	}

	user := identData.Identity
	perms, permsErr := model.ActionsForIdentity(user.Code, context)
	if permsErr != nil {
		return nil, permsErr
	}

	var sessionCode string
	if identData.Session == nil {
		token := user.IdentityToken(map[string]bool{})

		token.ClientID = client.ClientId
		token.Context = context

		encToken, err := kcutil.EncodeJWTOpen(token)
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
					if err := kcutil.DecodeJWTOpen(cookieVal, &encData); err == nil {
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
			MaxAge:   int(time.Until(time.Unix(token.Expiration, 0)).Seconds()),
		})
	} else {
		// TODO: make sure that the session we have is valid, don't just trust the token
		sessionCode = *identData.Session
	}

	accessContext := model.AccessContext{
		Session:   &sessionCode,
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
		Code:      user.Code,
		Browser:   sessionCode,
		Context:   accessContext.Code,
		Strength:  identData.Strength,
		Method:    identData.Method,
		Permitted: perms,
		// TODO: Can this be populated with at least the set of clients the user can hit?
		// Maybe intersect that with the set of resources that this client might direct a user to hit?
		// That would require tracking that, which wouldn't be the worst idea...
		//		ValidResource: []string{client.ClientId},
	}, nil
}

func clearSessionCookie(c *gin.Context, cookieName, domain string) {
	http.SetCookie(c.Writer, &http.Cookie{
		Domain:   domain,
		Name:     cookieName,
		Value:    "",
		Path:     "/",
		Secure:   !viper.GetBool("development.insecure"),
		HttpOnly: true,
		MaxAge:   -1,
	})
}
