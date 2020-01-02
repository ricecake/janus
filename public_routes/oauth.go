package public_routes

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ricecake/osin"
	"github.com/spf13/viper"

	"github.com/ricecake/janus/model"
	"github.com/ricecake/janus/util"
)

func discovery(c *gin.Context) {
	issuer := viper.GetString("identity.issuer")

	// For other example see: https://accounts.google.com/.well-known/openid-configuration
	c.JSON(200, map[string]interface{}{
		"issuer":                                issuer,
		"authorization_endpoint":                issuer + "/login",
		"token_endpoint":                        issuer + "/token",
		"userinfo_endpoint":                     issuer + "/userinfo",
		"jwks_uri":                              issuer + "/publickeys",
		"response_types_supported":              []string{"code", "token"},
		"subject_types_supported":               []string{"public"},
		"id_token_signing_alg_values_supported": []string{"RS256"},
		"scopes_supported":                      []string{"openid", "profile", "groups", "roles", "grants"},
		"token_endpoint_auth_methods_supported": []string{"client_secret_basic"},
		"claims_supported": []string{
			"iss",
			"sub",
			"aud",
			"exp",
			"iat",
			"jti",
			"nonce",
			"email",
			"family_name",
			"given_name",
			"preferred_name",
			"groups",
			"roles",
			"grants",
		},
	})
}

func publicKeys(c *gin.Context) {
	c.JSON(200, util.Keys)
}

func userInfo(c *gin.Context) {}

func accessToken(c *gin.Context) {
	response := server.NewResponse()
	defer response.Close()

	if ar := server.HandleAccessRequest(response, c.Request); ar != nil {
		switch ar.Type {
		case osin.AUTHORIZATION_CODE:
			ar.Authorized = true
		case osin.REFRESH_TOKEN:
			ar.Authorized = true
		}

		if ar.UserData != nil {
			authDetails := ar.UserData.(*model.UserAuthDetails)
			ident, err := model.FindIdentityById(authDetails.Code)
			if err != nil {
				response.InternalError = err
			} else {
				scopes := make(map[string]bool)
				for _, s := range strings.Fields(ar.Scope) {
					scopes[s] = true
				}
				token := ident.IdentityToken(scopes)

				token.ClientID = ar.Client.GetId()
				token.Nonce = authDetails.Nonce

				encToken, err := util.EncodeJWTOpen(token)
				if err != nil {
					response.InternalError = err
				} else {
					response.Output["id_token"] = encToken
				}
			}

			// Record errors as internal server errors.
			if response.InternalError != nil {
				response.IsError = true
				response.ErrorId = osin.E_SERVER_ERROR
			}
		}

		server.FinishAccessRequest(response, c.Request, ar)
	}

	if response.IsError && response.InternalError != nil {
		c.Error(response.InternalError).SetType(gin.ErrorTypePrivate)
	}

	osin.OutputJSON(response, c.Writer, c.Request)

}
