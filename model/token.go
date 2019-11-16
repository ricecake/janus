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

type RevocationEntry struct{}

type StashToken struct{}

type TokenGenerator struct{}

type UserAuthDetails struct {
	Code  string
	Nonce string
}
type AuthCodeData struct {
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
	log.Printf("REQUEST %+v", data)

	codeData := AuthCodeData{
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
