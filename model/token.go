package model

import (
	"time"

	"github.com/openshift/osin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/ricecake/janus/util"
)

type SessionToken struct{}

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
	Code  string
	Nonce string
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
	SessionCode   string `json:"sess,omitempty"`
}
type RefreshToken struct {
	TokenCode   string `json:"jti"`
	AccessCode  string `json:"ati"`
	SessionCode string `json:"sti"`
}

func (a *TokenGenerator) GenerateAccessToken(data *osin.AccessData, generaterefresh bool) (accessToken string, refreshToken string, err error) {
	accessTokenData := AccessToken{
		Issuer:      viper.GetString("identity.issuer"),
		UserCode:    "potato",
		Expiration:  data.CreatedAt.Add(time.Duration(data.ExpiresIn) * time.Second).Unix(),
		IssuedAt:    data.CreatedAt.Unix(),
		TokenCode:   util.CompactUUID(),
		SessionCode: "",
		// Nonce:  ,
	}

	accessToken, err = util.EncodeJWTOpen(accessTokenData)
	if err != nil {
		return
	}

	if generaterefresh {
		refreshToken, err = util.EncodeJWTClose(RefreshToken{
			TokenCode:   util.CompactUUID(),
			AccessCode:  accessTokenData.TokenCode,
			SessionCode: "",
		}, viper.GetString("security.passphrase"))
	}
	return
}
