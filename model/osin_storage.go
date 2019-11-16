package model

import (
	"github.com/openshift/osin"
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
	return nil, osin.ErrNotFound
}

func (s *DbStorage) RemoveAuthorize(code string) error {
	return nil
}

func (s *DbStorage) SaveAccess(data *osin.AccessData) error {
	return nil
}

func (s *DbStorage) LoadAccess(code string) (*osin.AccessData, error) {
	return nil, osin.ErrNotFound
}

func (s *DbStorage) RemoveAccess(code string) error {
	return nil
}

func (s *DbStorage) LoadRefresh(code string) (*osin.AccessData, error) {
	return nil, osin.ErrNotFound
}

func (s *DbStorage) RemoveRefresh(code string) error {
	return nil
}

func (s *DbStorage) SaveAuthorize(data *osin.AuthorizeData) error {
	return nil
}
