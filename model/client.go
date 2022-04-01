package model

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ricecake/karma_chameleon/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Client struct {
	Context     string `gorm:"column:context;not null"`
	DisplayName string `gorm:"column:display_name;not null"`
	ClientId    string `gorm:"column:client_id;not null; primary_key"`
	Secret      []byte `gorm:"column:secret;not null" json:"-"`
	BaseUri     string `gorm:"column:base_uri;not null"`
	Description string `gorm:"column:description; not null"`
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

	// TODO: insert the client code as an action for the context.

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

func ActionsForClient(identCode, context string) (allowed []string, err error) {
	/**
	For Now, this is just going to return all actions in a context.
	Long term, should have tables that can put clients into groups,
	and also assign actions to specific clients.
	**/
	if identCode == "" {
		err = fmt.Errorf("No Identity passed")
		return
	}

	db := util.GetDb()
	var results []Action

	err = db.Where("context = ?", context).Find(&results).Error

	for _, act := range results {
		allowed = append(allowed, act.Name)
	}

	sort.Strings(allowed)
	return
}

func ListClients() ([]Client, error) {
	db := util.GetDb()
	var clients []Client
	err := db.Find(&clients).Error
	return clients, err
}

func (this *Client) SaveChanges() error {
	db := util.GetDb()

	err := db.Save(this).Error
	if err != nil {
		log.Errorf("Error while updating client: %s", err)
	}

	return err
}
