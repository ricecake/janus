package public_routes

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/openshift/osin"
	log "github.com/sirupsen/logrus"

	"github.com/ricecake/janus/model"
	"github.com/ricecake/janus/util"
)

var (
	server   *osin.Server
	tokenGen model.TokenGenerator
)

func init() {
	sconfig := osin.NewServerConfig()
	sconfig.AllowClientSecretInParams = true
	sconfig.RequirePKCEForPublicClients = true
	sconfig.AllowedAuthorizeTypes = osin.AllowedAuthorizeType{osin.CODE, osin.TOKEN}
	sconfig.AllowedAccessTypes = osin.AllowedAccessType{
		osin.AUTHORIZATION_CODE,
		osin.REFRESH_TOKEN,
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

	body, renderErr := util.RenderTemplate("login", util.TemplateContext{
		"Name": client.DisplayName,
	})

	if renderErr != nil {
		c.Error(renderErr).SetType(gin.ErrorTypePrivate)
		c.AbortWithError(500, fmt.Errorf("System Error")).SetType(gin.ErrorTypePublic)
		return
	}

	c.Data(200, "text/html", body)
}
func loginSubmit(c *gin.Context)  {}
func signupPage(c *gin.Context)   {}
func signupSubmit(c *gin.Context) {}
func logoutPage(c *gin.Context)   {}
func logoutSubmit(c *gin.Context) {}
