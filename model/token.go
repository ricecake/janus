package model

import (
	"time"

	"github.com/openshift/osin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/ricecake/janus/util"
)

type SessionToken struct {
	Code      string    `gorm:"column:code;not null"`
	Identity  string    `gorm:"column:identity;not null"`
	UserAgent string    `gorm:"column:user_agent;not null"`
	IpAddress string    `gorm:"column:ip_address;not null"`
	CreatedAt time.Time `gorm:"column:created_at;not null"`
	ExpiresIn int       `gorm:"column:expires_in;not null"`
}

func (this SessionToken) TableName() string {
	return "session_token"
}

func CreateSessionToken(tok *SessionToken) error {
	db := util.GetDb()

	if tok.Code == "" {
		tok.Code = util.CompactUUID()
	}
	log.Info("Creating session ", tok.Code)

	return db.Create(tok).Error
}

func InvalidateSessionToken(sessid string) error {
	db := util.GetDb()
	db.Where("code = ?", sessid).Delete(&SessionToken{})
	return InsertRevocation(sessid, int((time.Duration(viper.GetInt("identity.ttl")) * time.Hour).Seconds()))
}

type AccessContext struct {
	Code      string    `gorm:"column:code;not null"`
	Session   string    `gorm:"column:session;not null"`
	Client    string    `gorm:"column:client;not null"`
	CreatedAt time.Time `gorm:"column:created_at;not null"`
}

func (this AccessContext) TableName() string {
	return "access_context"
}

func EnsureAccessContext(con *AccessContext) error {
	db := util.GetDb()

	if db.Where("session = ? AND client = ?", con.Session, con.Client).Find(&con).RecordNotFound() {
		log.Info("Creating context ", con.Code)
		con.Code = util.CompactUUID()
		return db.Create(&con).Error
	}

	return nil
}

type RevocationEntry struct {
	EntityCode string    `gorm:"column:entity_code;not null"`
	CreatedAt  time.Time `gorm:"column:created_at;not null"`
	ExpiresIn  int       `gorm:"column:expires_in;not null"`
}

func (this RevocationEntry) TableName() string {
	return "revocation"
}

func InsertRevocation(entity string, ttl int) error {
	db := util.GetDb()

	revocationEntry := RevocationEntry{
		EntityCode: entity,
		CreatedAt:  time.Now(),
		ExpiresIn:  ttl,
	}

	log.Info("Revoking entity ", revocationEntry.EntityCode)

	return db.Create(&revocationEntry).Error
}

func EntityRevoked(entity string) bool {
	db := util.GetDb()
	return !db.Where("entity_code = ?", entity).First(&RevocationEntry{}).RecordNotFound()
}

type StashToken struct{}

type TokenGenerator struct{}

type UserAuthDetails struct {
	Code     string
	Nonce    string
	Browser  string
	Context  string
	Strength int
	Method   string
}
type AuthCodeData struct {
	Code                string
	ClientId            string
	ExpiresIn           int32
	Scope               string
	RedirectUri         string
	State               string
	CreatedAt           time.Time
	UserData            *UserAuthDetails
	CodeChallenge       string
	CodeChallengeMethod string
}

func (a *TokenGenerator) GenerateAuthorizeToken(data *osin.AuthorizeData) (ret string, err error) {
	codeData := AuthCodeData{
		Code:                util.CompactUUID(),
		ClientId:            data.Client.GetId(),
		ExpiresIn:           data.ExpiresIn,
		Scope:               data.Scope,
		RedirectUri:         data.RedirectUri,
		State:               data.State,
		CreatedAt:           data.CreatedAt,
		CodeChallenge:       data.CodeChallenge,
		CodeChallengeMethod: data.CodeChallengeMethod,
	}
	if data.UserData != nil {
		codeData.UserData = data.UserData.(*UserAuthDetails)
	}

	return util.EncodeJWTClose(codeData, viper.GetString("security.passphrase"))
}

type AccessToken struct {
	Issuer     string `json:"iss"`
	UserCode   string `json:"sub"`
	Expiration int64  `json:"exp"`
	IssuedAt   int64  `json:"iat"`
	TokenCode  string `json:"jti"`

	Nonce         string `json:"nonce,omitempty"` // Non-manditory fields MUST be "omitempty"
	ValidResource string `json:"aud,omitempty"`
	ContextCode   string `json:"ctx,omitempty"`
}
type RefreshToken struct {
	TokenCode   string `json:"jti"`
	AccessCode  string `json:"ati"`
	ContextCode string `json:"ctx"`
}

func (a *TokenGenerator) GenerateAccessToken(data *osin.AccessData, generaterefresh bool) (accessToken string, refreshToken string, err error) {
	accessTokenData := AccessToken{
		Issuer:     viper.GetString("identity.issuer"),
		Expiration: data.CreatedAt.Add(time.Duration(data.ExpiresIn) * time.Second).Unix(),
		IssuedAt:   data.CreatedAt.Unix(),
		TokenCode:  util.CompactUUID(),
	}

	if data.UserData != nil {
		authDetails := data.UserData.(*UserAuthDetails)
		accessTokenData.Nonce = authDetails.Nonce
		accessTokenData.ContextCode = authDetails.Context
		accessTokenData.UserCode = authDetails.Code
	}

	accessToken, err = util.EncodeJWTOpen(accessTokenData)
	if err != nil {
		return
	}

	if generaterefresh {
		refreshToken, err = util.EncodeJWTClose(RefreshToken{
			TokenCode:   util.CompactUUID(),
			AccessCode:  accessTokenData.TokenCode,
			ContextCode: "",
		}, viper.GetString("security.passphrase"))
	}
	return
}

type IDToken struct {
	Issuer     string `json:"iss"`
	UserCode   string `json:"sub"`
	ClientID   string `json:"aud"`
	IssuedAt   int64  `json:"iat"`
	Expiration int64  `json:"exp"`
	TokenId    string `json:"jti"`

	Nonce    string `json:"nonce,omitempty"` // Non-manditory fields MUST be "omitempty"
	Strength string `json:"acr,omitempty"`
	Methods  string `json:"amr,omitempty"`

	// Custom claims supported by this server.
	//
	// See: https://openid.net/specs/openid-connect-core-1_0.html#StandardClaims

	Email         string `json:"email,omitempty"`
	FamilyName    string `json:"family_name,omitempty"`
	GivenName     string `json:"given_name,omitempty"`
	PreferredName string `json:"preferred_name,omitempty"`
}
