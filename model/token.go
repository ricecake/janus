package model

import (
	"time"

	"github.com/openshift/osin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/ricecake/janus/util"
)

type SessionToken struct{}

type AccessToken struct{}

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

func (a *TokenGenerator) GenerateAccessToken(data *osin.AccessData, generaterefresh bool) (accesstoken string, refreshtoken string, err error) {
	log.Printf("REQUEST %+v", data)
	accesstoken = ""

	if generaterefresh {
		refreshtoken = ""
	}
	return
}
