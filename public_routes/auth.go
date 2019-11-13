package public_routes

import (
	"github.com/openshift/osin"
	"github.com/ricecake/janus/model"
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
