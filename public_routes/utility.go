package public_routes

import (
	"fmt"
	"net/http"
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
						// TODO: should do something to move access context entries to the new token
						if err := model.InvalidateSessionToken(encData.TokenId); err != nil {
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
