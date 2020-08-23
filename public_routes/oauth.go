package public_routes

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ricecake/osin"
	"github.com/spf13/viper"

	"github.com/ricecake/janus/model"
	"github.com/ricecake/karma_chameleon/util"
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
		var authorized bool
		var auth_decided bool

		switch ar.Type {
		case osin.CLIENT_CREDENTIALS:
			client, clientErr := model.FindClientById(ar.Client.GetId())
			if clientErr != nil {
				response.InternalError = clientErr
				break
			}
			if client.ClientSecretMatches("") {
				// Don't allow public clients to be used this way
				if !auth_decided {
					authorized = false
					auth_decided = true
				}

			}

			perms, actErr := model.ActionsForClient(client.ClientId, client.Context)
			if actErr != nil {
				response.InternalError = actErr
				break
			}

			authorized = true
			ar.UserData = &model.UserAuthDetails{
				Code:      client.ClientId,
				Context:   client.Context,
				Strength:  "1",
				Method:    "client credentials",
				Permitted: perms,
			}

		case osin.PASSWORD:
			client, clientErr := model.FindClientById(ar.Client.GetId())
			if clientErr != nil {
				response.InternalError = clientErr
				break
			}
			context := client.Context

			identData := attemptIdentifyUser(c, model.IdentificationRequest{
				Strategy: model.PASSWORD,
				Context:  &client.Context,
			})

			permitted := identData.Success
			user := identData.Identity

			if permitted {
				allowed, err := model.AclCheck(model.AclCheckRequest{
					Identity: identData.Identity.Code,
					Context:  client.Context,
					Action:   client.ClientId,
				})

				if err != nil {
					response.InternalError = err
					break
				}

				permitted = permitted && allowed
			}

			if !auth_decided {
				authorized = permitted
				auth_decided = true
			}

			if !authorized {
				break
			}

			perms, permsErr := model.ActionsForIdentity(user.Code, context)
			if permsErr != nil {
				response.InternalError = permsErr
				break
			}

			accessContext := model.AccessContext{
				Client:    client.ClientId,
				CreatedAt: time.Now(),
			}
			if err := model.EnsureAccessContext(&accessContext); err != nil {
				response.InternalError = err
				break
			}

			ar.UserData = &model.UserAuthDetails{
				Code:          user.Code,
				Context:       accessContext.Code,
				Strength:      identData.Strength,
				Method:        identData.Method,
				ValidResource: []string{client.ClientId},
				Permitted:     perms,
			}

			fallthrough
		case osin.AUTHORIZATION_CODE:
			if !auth_decided {
				authorized = true
				auth_decided = true
			}
			fallthrough
		case osin.REFRESH_TOKEN:
			if !auth_decided {
				authorized = true
				auth_decided = true
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
			}
		}

		ar.Authorized = authorized

		// Record errors as internal server errors.
		if response.InternalError != nil {
			response.IsError = true
			response.ErrorId = osin.E_SERVER_ERROR
		}

		server.FinishAccessRequest(response, c.Request, ar)
	}

	if response.IsError && response.InternalError != nil {
		c.Error(response.InternalError).SetType(gin.ErrorTypePrivate)
	}

	osin.OutputJSON(response, c.Writer, c.Request)

}
