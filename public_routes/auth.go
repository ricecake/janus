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

	res, err := attemptIdentifyUser(c, model.PASSWORD)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePrivate)
		c.AbortWithError(500, fmt.Errorf("System Error")).SetType(gin.ErrorTypePublic)
	}

	ar.Authorized = res.Success

	if res.Success {
		userDetails, err := establishSession(c, *res.Identity)
		if err != nil {
			c.Error(err).SetType(gin.ErrorTypePrivate)
			c.AbortWithError(500, fmt.Errorf("System Error")).SetType(gin.ErrorTypePublic)
		}
		ar.UserData = userDetails
	}

	server.FinishAuthorizeRequest(response, c.Request, ar)

	if response.IsError && response.InternalError != nil {
		log.Printf("internal error: %v", response.InternalError)
	}
	osin.OutputJSON(response, c.Writer, c.Request)
}

func signupPage(c *gin.Context)   {}
func signupSubmit(c *gin.Context) {}
func logoutPage(c *gin.Context)   {}
func logoutSubmit(c *gin.Context) {}

type AuthParams struct {
	Email    *string `form:"email" json:"email"`
	Password *string `form:"password" json:"password"`
	Totp     *string `form:"totp" json:"totp"`
}

func attemptIdentifyUser(c *gin.Context, preference model.IdentificationStrategy) (*model.IdentificationResult, error) {
	var authData model.IdentificationRequest

	authData.Strategy = preference

	if authData.Strategy == model.NONE || authData.Strategy == model.PASSWORD {
		var authParams AuthParams
		if c.ShouldBind(&authParams) == nil {
			if authParams.Email != nil {
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
	return model.IdentifyFromCredentials(authData)
}
