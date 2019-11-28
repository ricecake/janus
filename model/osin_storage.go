package model

import (
	"github.com/openshift/osin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/ricecake/janus/util"
)

type DbStorage struct{}

func NewDbStorage() *DbStorage {
	return &DbStorage{}
}

func (s *DbStorage) Clone() osin.Storage {
	return s
}

func (s *DbStorage) Close() {
}

func (s *DbStorage) GetClient(id string) (osin.Client, error) {
	client, clientErr := FindClientById(id)
	if clientErr != nil {
		return nil, osin.ErrNotFound
	}
	return client, nil
}

func (s *DbStorage) LoadAuthorize(code string) (*osin.AuthorizeData, error) {
	var encData AuthCodeData
	if err := util.DecodeJWTClose(code, viper.GetString("security.passphrase"), &encData); err != nil {
		return nil, osin.ErrNotFound
	}

	//TODO: validate code not expired

	client, clientErr := FindClientById(encData.ClientId)
	if clientErr != nil {
		return nil, osin.ErrNotFound
	}

	if EntityRevoked(encData.Code) {
		return nil, osin.ErrNotFound
	}

	return &osin.AuthorizeData{
		Code:                code,
		Client:              client,
		ExpiresIn:           encData.ExpiresIn,
		Scope:               encData.Scope,
		RedirectUri:         encData.RedirectUri,
		State:               encData.State,
		CreatedAt:           encData.CreatedAt,
		CodeChallenge:       encData.CodeChallenge,
		CodeChallengeMethod: encData.CodeChallengeMethod,
		UserData:            encData.UserData,
	}, nil
}

func (s *DbStorage) RemoveAuthorize(code string) error {
	var encData AuthCodeData
	if err := util.DecodeJWTClose(code, viper.GetString("security.passphrase"), &encData); err != nil {
		log.Error(err)
		return err
	}

	return InsertRevocation(encData.Code, int(encData.ExpiresIn))
}

func (s *DbStorage) LoadAccess(code string) (*osin.AccessData, error) {
	return nil, osin.ErrNotFound
}

func (s *DbStorage) RemoveAccess(code string) error {
	var encData AccessToken
	if err := util.DecodeJWTOpen(code, &encData); err != nil {
		log.Error(err)
		return err
	}

	return InsertRevocation(encData.TokenCode, int(encData.Expiration-encData.IssuedAt))
}

func (s *DbStorage) LoadRefresh(code string) (*osin.AccessData, error) {
	return nil, osin.ErrNotFound
}

func (s *DbStorage) RemoveRefresh(code string) error {
	return nil
}

func (s *DbStorage) SaveAccess(data *osin.AccessData) error {
	return nil
}

func (s *DbStorage) SaveAuthorize(data *osin.AuthorizeData) error {
	return nil
}
