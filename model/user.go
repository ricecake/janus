package model

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ricecake/janus/util"
)

type Identity struct {
	Code          string  `gorm:"column:code;not null"`
	Active        bool    `gorm:"column:active;not null"`
	Email         string  `gorm:"column:email;not null"`
	PreferredName string  `gorm:"column:preferred_name;not null"`
	GivenName     *string `gorm:"column:given_name"`
	FamilyName    *string `gorm:"column:family_name"`
}

func (this Identity) TableName() string {
	return "identity"
}

type AuthPassword struct {
	Identity  string    `gorm:"column:identity;not null"`
	Hash      []byte    `gorm:"column:hash;not null"`
	Totp      *[]byte   `gorm:"column:totp;not null"`
	CreatedAt time.Time `gorm:"column:created_at;not null"`
}

func (this AuthPassword) TableName() string {
	return "auth_password"
}

func CreateIdentity(ident *Identity) error {
	db := util.GetDb()

	ident.Code = util.CompactUUID()

	log.Info("Creating user ", ident.Code)

	return db.Create(ident).Error
}

func (this *Identity) SetPassword(password string) (err error) {
	db := util.GetDb()

	hash, err := util.PasswordHash([]byte(password))
	if err != nil {
		return fmt.Errorf("hash failed")
	}

	auth := AuthPassword{
		Identity:  this.Code,
		CreatedAt: time.Now(),
	}
	if db.Where("identity = ?", this.Code).Find(&auth).RecordNotFound() {
		auth.Hash = hash
		err = db.Create(&auth).Error
	} else {
		auth.Hash = hash
		err = db.Save(&auth).Error
	}

	if err != nil {
		log.Errorf("Error while synchronizing browser token: %s", err)
	}

	return err
}
