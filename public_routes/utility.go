package public_routes

import (
	"net/http"
	"fmt"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/ricecake/janus/model"
)

func defaultPage(c *gin.Context)   {}
func checkUsername(c *gin.Context) {}
func checkAuth(c *gin.Context)     {}

func establishSession(c *gin.Context, user model.User) (*model.UserAuthDetails, error) {
	client, clientErr := model.FindClientById(viper.GetString("identity.issuer_id"))
	if clientErr != nil {
		log.Error("Error with own client?")
		return nil, clientErr
	}

	log.Info("Generate an id token -- method should live off of user, since it's their id")
	log.Info("record record of it in db")
	log.Info("create session record referencing it")

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     fmt.Sprintf("janus.auth.%s", client.Context),
		Value:    "pancake",
		Path:     "/",
		Secure:   !viper.GetBool("development.insecure"),
		HttpOnly: true,
	})

	return &model.UserAuthDetails{}, nil
}
