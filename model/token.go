package model

import (
	"github.com/openshift/osin"
)

type SessionToken struct{}

type AccessToken struct{}

type RevocationEntry struct{}

type StashToken struct{}

type TokenGenerator struct{}

func (a *TokenGenerator) GenerateAuthorizeToken(data *osin.AuthorizeData) (ret string, err error) {
	return "", nil
}

func (a *TokenGenerator) GenerateAccessToken(data *osin.AccessData, generaterefresh bool) (accesstoken string, refreshtoken string, err error) {
	accesstoken = ""

	if generaterefresh {
		refreshtoken = ""
	}
	return
}
