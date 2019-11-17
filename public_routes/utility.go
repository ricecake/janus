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
func checkAuth(c *gin.Context)     {}

func establishSession(c *gin.Context, user model.Identity) (*model.UserAuthDetails, error) {
	client, clientErr := model.FindClientById(viper.GetString("identity.issuer_id"))
	if clientErr != nil {
		log.Error("Error with own client?")
		return nil, clientErr
	}

	token := user.IdentityToken(map[string]bool{})

	encToken, err := util.EncodeJWTOpen(token)
	if err != nil {
		return nil, err
	}

	sessionToken := model.SessionToken{
		Identity:  user.Code,
		UserAgent: "test",
		IpAddress: "127.0.0.1",
		CreatedAt: time.Now(),
		ExpiresIn: int(token.Expiration - token.IssuedAt),
	}
	if err := model.CreateSessionToken(&sessionToken); err != nil {
		return nil, err
	}

	accessContext := model.AccessContext{
		Session:   sessionToken.Code,
		Client:    client.ClientId,
		CreatedAt: time.Now(),
	}
	if err := model.CreateAccessContext(&accessContext); err != nil {
		return nil, err
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "janus.user.email",
		Value:    user.Email,
		Path:     "/",
		Secure:   !viper.GetBool("development.insecure"),
		HttpOnly: false,
	})

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     fmt.Sprintf("janus.auth.%s", client.Context),
		Value:    encToken,
		Path:     "/",
		Secure:   !viper.GetBool("development.insecure"),
		HttpOnly: true,
	})

	return &model.UserAuthDetails{
		Code:    user.Code,
		Browser: sessionToken.Code,
		Context: accessContext.Code,
	}, nil
}
