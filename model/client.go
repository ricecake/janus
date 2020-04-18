package model

import (
	"fmt"
	"strings"

	"github.com/ricecake/karma_chameleon/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Client struct {
	Context     string `gorm:"column:context;not null"`
	DisplayName string `gorm:"column:display_name;not null"`
	ClientId    string `gorm:"column:client_id;not null"`
	Secret      []byte `gorm:"column:secret;not null" json:"-"`
	BaseUri     string `gorm:"column:base_uri;not null"`
}

func (this Client) TableName() string {
	return "client"
}

func (this Client) GetId() string {
	return this.ClientId
}

func CreateClient(client *Client) error {
	db := util.GetDb()

	client.ClientId = util.CompactUUID()

	log.Info("Creating client ", client.ClientId)

	return db.Create(client).Error
}

func (this *Client) SetSecret(plainSecret string) error {
	hash, err := util.PasswordHash([]byte(plainSecret))
	if err != nil {
		return fmt.Errorf("hash failed")
	}

	this.Secret = hash
	return nil
}

func (this Client) GetSecret() (secret string) {
	log.Fatal("Insecure password access attempt")
	return
}

func (this Client) GetRedirectUri() string {
	return strings.Join([]string{this.BaseUri, viper.GetString("identity.issuer")}, "|")
}

func (this Client) GetUserData() interface{} {
	return nil
}

func (this Client) ClientSecretMatches(plainSecret string) bool {
	return util.PasswordHashValid([]byte(plainSecret), this.Secret)
}

func FindClientById(id string) (client Client, err error) {
	db := util.GetDb()
	if db.Where("client_id = ?", id).Find(&client).RecordNotFound() {
		err = fmt.Errorf("Invalid client id")
	}
	return client, err
}
