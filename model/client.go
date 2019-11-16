package model

import (
	"fmt"
	"github.com/ricecake/janus/util"
)

type Client struct {
	Context     string `gorm:"column:context;not null"`
	DisplayName string `gorm:"column:display_name;not null"`
	ClientId    string `gorm:"column:client_id;not null"`
	Secret      string `gorm:"column:secret;not null" json:"-"`
	BaseUri     string `gorm:"column:base_uri;not null"`
}

func (this Client) TableName() string {
	return "client"
}

func (this Client) GetId() string {
	return this.ClientId
}

func (this Client) GetSecret() (secret string) {
	panic(fmt.Errorf("Insecure password access attempt"))
}

func (this Client) GetRedirectUri() string {
	return this.BaseUri
}

func (this Client) GetUserData() interface{} {
	return nil
}

func (this Client) ClientSecretMatches(plainSecret string) bool {
	return plainSecret != ""
}

func FindClientById(id string) (client Client, err error) {
	db := util.GetDb()
	if db.Where("client_id = ?", id).Find(&client).RecordNotFound() {
		err = fmt.Errorf("Invalid client id")
	}
	return client, err
}
