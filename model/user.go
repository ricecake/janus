package model

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/ricecake/janus/util"
)

type Identity struct {
	Code          string `gorm:"column:code;not null"`
	Active        bool   `gorm:"column:active;not null"`
	Email         string `gorm:"column:email;not null"`
	PreferredName string `gorm:"column:preferred_name;not null"`
	GivenName     string `gorm:"column:given_name"`
	FamilyName    string `gorm:"column:family_name"`
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

func (this Identity) IdentityToken(claims map[string]bool) IDToken {
	issued := time.Now()
	expires := issued.Add(time.Duration(viper.GetInt("identity.ttl")) * time.Hour)

	token := IDToken{
		UserCode:   this.Code,
		IssuedAt:   int64(issued.Unix()),
		Expiration: int64(expires.Unix()),
		TokenId:    util.CompactUUID(),
	}

	if claims["profile"] {
		token.Email = this.Email
		token.PreferredName = this.PreferredName
		token.FamilyName = this.FamilyName
		token.GivenName = this.GivenName
	}

	return token
}

type IdentificationStrategy int

const (
	NONE IdentificationStrategy = iota
	PASSWORD
	SESSION_TOKEN
	WEBAUTHN
)

type IdentificationRequest struct {
	Strategy     IdentificationStrategy
	Email        *string
	Password     *string
	Totp         *string
	SessionToken *string
}
type IdentificationResult struct {
	Success       bool
	Identity      *Identity
	Strength      string
	Method        string
	FailureCode   int
	FailureReason string
}

func IdentifyFromCredentials(req IdentificationRequest) (*IdentificationResult, error) {
	switch req.Strategy {
	case NONE:
		return &IdentificationResult{
			FailureCode:   500,
			FailureReason: "Bad auth attempt",
		}, nil
	case PASSWORD:
		db := util.GetDb()
		var ident Identity
		if db.Where("email = ?", *req.Email).Find(&ident).RecordNotFound() {
			return &IdentificationResult{
				FailureCode:   401,
				FailureReason: "Bad user",
			}, nil
		}
		var auth AuthPassword
		if db.Where("identity = ?", ident.Code).Find(&auth).RecordNotFound() {
			return &IdentificationResult{
				FailureCode:   401,
				FailureReason: "Bad auth method",
			}, nil
		}
		if !util.PasswordHashValid([]byte(*req.Password), auth.Hash) {
			return &IdentificationResult{
				FailureCode:   401,
				FailureReason: "Bad password",
			}, nil
		}
		return &IdentificationResult{
			Success:  true,
			Identity: &ident,
			Strength: "1",
			Method:   "password",
		}, nil
	case SESSION_TOKEN:
		return &IdentificationResult{
			FailureCode:   500,
			FailureReason: "Unknown auth method",
		}, nil
	default:
		return &IdentificationResult{
			FailureCode:   500,
			FailureReason: "Unknown auth method",
		}, nil
	}
}
