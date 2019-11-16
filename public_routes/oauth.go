package public_routes

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
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
			"username",
			"name",
			"family_name",
			"given_name",
			"preferred_name",
			"groups",
			"roles",
			"grants",
		},
	})

}

func publicKeys(c *gin.Context)  {}
func userInfo(c *gin.Context)    {}
func accessToken(c *gin.Context) {}
