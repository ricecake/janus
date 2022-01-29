package model

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/ricecake/karma_chameleon/util"
)

type Identity struct {
	Code          string `gorm:"column:code;not null;primary_key"`
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
	Identity   string    `gorm:"column:identity;not null"`
	Hash       []byte    `gorm:"column:hash;not null"`
	Totp       *[]byte   `gorm:"column:totp;not null"`
	TotpActive bool      `gorm:"column:totp_active;not null"`
	CreatedAt  time.Time `gorm:"column:created_at;not null;autoCreateTime"`
}

func (this AuthPassword) TableName() string {
	return "auth_password"
}

type WebauthnCredential struct {
	Identity               string
	Id                     string
	PublicKey              string
	AttestationType        string
	AuthenticatorGUID      string
	AuthenticatorSignCount int
}

func (wc WebauthnCredential) TableName() string {
	return "webauthn_credential"
}

func CreateIdentity(ident *Identity) error {
	db := util.GetDb()

	ident.Code = util.CompactUUID()

	log.Info("Creating user ", ident.Code)

	return db.Create(ident).Error
}

func (this *Identity) SaveChanges() error {
	db := util.GetDb()

	err := db.Save(this).Error
	if err != nil {
		log.Errorf("Error while updating identity: %s", err)
	}

	return err
}

func FindIdentityById(id string) (ident Identity, err error) {
	db := util.GetDb()
	if db.Where("code = ?", id).Find(&ident).RecordNotFound() {
		err = fmt.Errorf("Invalid user id")
	}
	return ident, err
}

func FindIdentityByEmail(id string) (ident Identity, err error) {
	db := util.GetDb()
	if db.Where("email = ? and active", id).Find(&ident).RecordNotFound() {
		err = fmt.Errorf("Invalid email")
	}
	return ident, err
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

func (this *Identity) AddWebauthnCredential(wc *webauthn.Credential) (err error) {
	db := util.GetDb()
	cred := WebauthnCredential{
		Identity:               this.Code,
		Id:                     base64.StdEncoding.EncodeToString(wc.ID),
		PublicKey:              base64.StdEncoding.EncodeToString(wc.PublicKey),
		AttestationType:        wc.AttestationType,
		AuthenticatorGUID:      base64.StdEncoding.EncodeToString(wc.Authenticator.AAGUID),
		AuthenticatorSignCount: int(wc.Authenticator.SignCount),
	}
	return db.Create(&cred).Error
}

func (this Identity) IdentityToken(claims map[string]bool) IDToken {
	issued := time.Now()
	expires := issued.Add(time.Duration(viper.GetInt("identity.ttl")) * time.Hour)

	token := IDToken{
		UserCode:   this.Code,
		IssuedAt:   int64(issued.Unix()),
		Expiration: int64(expires.Unix()),
		TokenId:    util.CompactUUID(),
		Issuer:     viper.GetString("identity.issuer"),
	}

	if claims["openid"] {
		if claims["profile"] {
			token.Email = this.Email
			token.PreferredName = this.PreferredName
			token.FamilyName = this.FamilyName
			token.GivenName = this.GivenName
		}
		if claims["roles"] {
			roles, err := IdentityRoles(this.Code)
			if err != nil {
				log.Error(err)
			} else {
				token.Roles = roles
			}
		}

		if claims["cliques"] {
		}
		if claims["actions"] {
		}
		if claims["auth_info"] {
		}
	}

	return token
}

type AvailableAuthMethods struct {
	Email    bool
	Password bool
	Totp     bool
	Webauthn bool
}

func (user Identity) ValidAuthMethods() (methods AvailableAuthMethods, err error) {
	db := util.GetDb()

	methods.Email = true

	var pwCount int64
	var wauthnCount int64

	db.Model(&AuthPassword{}).Where("identity = ?", user.Code).Count(&pwCount)
	db.Model(&WebauthnCredential{}).Where("identity = ?", user.Code).Count(&wauthnCount)

	methods.Password = pwCount > 0
	methods.Webauthn = wauthnCount > 0

	return
}

func (user Identity) WebAuthnID() []byte {
	return []byte(user.Code)
}

func (user Identity) WebAuthnName() string {
	return user.Email
}

func (user Identity) WebAuthnDisplayName() string {
	return user.PreferredName
}

func (user Identity) WebAuthnIcon() string {
	return ""
}

func (user Identity) WebAuthnCredentials() (creds []webauthn.Credential) {
	db := util.GetDb()
	var webCreds []WebauthnCredential
	err := db.Where("identity = ?", user.Code).Find(&webCreds).Error
	if err != nil {
		log.Error(err)
	} else {
		for _, cred := range webCreds {
			id, err := base64.StdEncoding.DecodeString(cred.Id)
			if err != nil {
				log.Error(err)
				continue
			}
			pubkey, err := base64.StdEncoding.DecodeString(cred.PublicKey)
			if err != nil {
				log.Error(err)
				continue
			}
			guid, err := base64.StdEncoding.DecodeString(cred.AuthenticatorGUID)
			if err != nil {
				log.Error(err)
				continue
			}

			creds = append(creds, webauthn.Credential{
				ID:              id,
				PublicKey:       pubkey,
				AttestationType: cred.AttestationType,
				Authenticator: webauthn.Authenticator{
					AAGUID:    guid,
					SignCount: uint32(cred.AuthenticatorSignCount),
				},
			})

		}
	}

	return
}

// CredentialExcludeList returns a CredentialDescriptor array filled
// with all the user's credentials
func (user Identity) CredentialExcludeList() []protocol.CredentialDescriptor {

	credentialExcludeList := []protocol.CredentialDescriptor{}
	for _, cred := range user.WebAuthnCredentials() {
		descriptor := protocol.CredentialDescriptor{
			Type:         protocol.PublicKeyCredentialType,
			CredentialID: cred.ID,
		}
		credentialExcludeList = append(credentialExcludeList, descriptor)
	}

	return credentialExcludeList
}

type IdentificationStrategy int

const (
	NONE IdentificationStrategy = iota
	PASSWORD
	SESSION_TOKEN
	WEBAUTHN
	ZIPCODE
)

type IdentificationRequest struct {
	Strategy     IdentificationStrategy
	Context      *string
	Email        *string
	Password     *string
	Totp         *string
	SessionToken *[]string
	ZipCode      *string
	Credential   *webauthn.Credential
}
type IdentificationResult struct {
	Success       bool
	Identity      *Identity
	Session       *string
	ZipCode       *ZipCode
	Strategy      IdentificationStrategy
	Strength      string
	Method        string
	FailureCode   int
	FailureReason string
}

func IdentifyFromCredentials(req IdentificationRequest) *IdentificationResult {
	switch req.Strategy {
	case NONE:
		return &IdentificationResult{
			FailureCode:   401,
			FailureReason: "Bad auth attempt",
		}
	case PASSWORD:
		db := util.GetDb()
		var ident Identity
		if db.Where("email = ?", *req.Email).Find(&ident).RecordNotFound() {
			return &IdentificationResult{
				FailureCode:   401,
				FailureReason: "Bad user",
			}
		}
		//TODO: verify user active
		var auth AuthPassword
		if db.Where("identity = ?", ident.Code).Find(&auth).RecordNotFound() {
			return &IdentificationResult{
				FailureCode:   401,
				FailureReason: "Bad auth method",
			}
		}
		if !util.PasswordHashValid([]byte(*req.Password), auth.Hash) {
			return &IdentificationResult{
				FailureCode:   401,
				FailureReason: "Bad password",
			}
		}
		// TODO: verify mfa
		return &IdentificationResult{
			Success:  true,
			Strategy: req.Strategy,
			Identity: &ident,
			Strength: "1",
			Method:   "password",
		}
	case SESSION_TOKEN:
		if req.SessionToken == nil || len(*req.SessionToken) == 0 || req.Context == nil {
			return &IdentificationResult{
				FailureCode:   401,
				FailureReason: "No token",
			}
		}
		results := []*IdentificationResult{}
		for _, token := range *req.SessionToken {
			var encData IDToken
			if err := util.DecodeJWTOpen(token, &encData); err != nil {
				results = append(results, &IdentificationResult{
					FailureCode:   401,
					FailureReason: err.Error(),
				})
				continue
			}

			now := time.Now()
			if now.Unix() >= encData.Expiration {
				results = append(results, &IdentificationResult{
					FailureCode:   401,
					FailureReason: "Expired",
				})
				continue
			}

			clientId := viper.GetString("identity.issuer_id")
			if encData.ClientID != clientId {
				results = append(results, &IdentificationResult{
					FailureCode:   401,
					FailureReason: "Bad token",
				})
				continue
			}

			if encData.Context != *req.Context {
				results = append(results, &IdentificationResult{
					FailureCode:   401,
					FailureReason: "Bad token",
				})
				continue
			}

			db := util.GetDb()
			var ident Identity
			if db.Where("code = ?", encData.UserCode).Find(&ident).RecordNotFound() {
				results = append(results, &IdentificationResult{
					FailureCode:   401,
					FailureReason: "Bad user",
				})
				continue
			}

			if EntityRevoked(encData.TokenId) {
				results = append(results, &IdentificationResult{
					FailureCode:   401,
					FailureReason: "Bad session",
				})
				continue
			}

			return &IdentificationResult{
				Success:  true,
				Strategy: req.Strategy,
				Identity: &ident,
				Strength: "0",
				Method:   "session",
				Session:  &encData.TokenId,
			}
		}
		return results[0]

	case ZIPCODE:
		if req.ZipCode == nil {
			return &IdentificationResult{
				FailureCode:   401,
				FailureReason: "missing code",
			}
		}

		zipCode, zipErr := FetchZipCode(*req.ZipCode)
		if zipErr != nil {
			return &IdentificationResult{
				FailureCode:   401,
				FailureReason: "missing code",
			}
		}

		db := util.GetDb()
		var ident Identity
		if db.Where("code = ?", zipCode.Identity).Find(&ident).RecordNotFound() {
			return &IdentificationResult{
				FailureCode:   401,
				FailureReason: "Bad user",
			}
		}
		return &IdentificationResult{
			Success:  true,
			Strategy: req.Strategy,
			Identity: &ident,
			ZipCode:  &zipCode,
			Strength: "1",
			Method:   "email possession",
		}
	case WEBAUTHN:
		db := util.GetDb()

		var webCred WebauthnCredential
		credId := base64.StdEncoding.EncodeToString(req.Credential.ID)
		if db.Where("id = ?", credId).Find(&webCred).RecordNotFound() {
			return &IdentificationResult{
				FailureCode:   401,
				FailureReason: "Bad credential",
			}
		}

		var ident Identity
		if db.Where("code = ?", webCred.Identity).Find(&ident).RecordNotFound() {
			return &IdentificationResult{
				FailureCode:   401,
				FailureReason: "Bad user",
			}
		}

		return &IdentificationResult{
			Success:  true,
			Strategy: req.Strategy,
			Identity: &ident,
			Strength: "3",
			Method:   "webauthn",
		}

	default:
		return &IdentificationResult{
			FailureCode:   500,
			FailureReason: "Unknown auth method",
		}
	}
}

type AclCheckRequest struct {
	Identity string  `gorm:"column:identity;not null"`
	Context  string  `gorm:"column:context;not null"`
	Clique   *string `gorm:"column:clique;"`
	Role     string  `gorm:"column:role;not null"`
	Action   string  `gorm:"column:action;not null"`
}

func (this AclCheckRequest) TableName() string {
	return "identity_access_summary"
}

func AclCheck(req AclCheckRequest) (allowed bool, err error) {
	if req.Identity == "" {
		err = fmt.Errorf("No Identity passed")
		return
	}
	if req.Context == "" {
		err = fmt.Errorf("No Context passed")
		return
	}
	if req.Action == "" {
		err = fmt.Errorf("No Action passed")
		return
	}

	db := util.GetDb()
	var count int

	model := db.Model(req)
	if req.Clique == nil {
		model = model.Where(req).Where("clique is null")
	} else {
		clique := req.Clique
		req.Clique = nil
		model = model.Where(req).Where("clique = ? or clique is null", clique)
		req.Clique = clique
	}
	err = model.Count(&count).Error
	allowed = count > 0

	return
}

func ActionsForIdentity(identCode, context string) (allowed []string, err error) {
	if identCode == "" {
		err = fmt.Errorf("No Identity passed")
		return
	}

	db := util.GetDb()
	var results []AclCheckRequest

	err = db.Where("identity = ? AND context = ?", identCode, context).Find(&results).Error

	for _, acl := range results {
		action := acl.Action
		if acl.Clique != nil {
			action = strings.Join([]string{*acl.Clique, action}, ".")
		}
		allowed = append(allowed, action)
	}

	sort.Strings(allowed) // dedupe
	return
}

type IdentitySummary struct {
	Identity string `gorm:"column:identity"`
	Email    string `gorm:"column:email"`
	Roles    []byte `gorm:"column:roles"`
	Actions  []byte `gorm:"column:actions"`
}

func (view IdentitySummary) TableName() string {
	return "identity_summary"
}

func IdentityRoles(identCode string) (map[string][]string, error) {
	db := util.GetDb()

	result := make(map[string][]string)
	var summary IdentitySummary

	err := db.Where("identity = ?", identCode).Find(&summary).Error
	if err != nil {
		return result, err
	}

	if unmarshalError := json.Unmarshal(summary.Roles, &result); unmarshalError != nil {
		return result, unmarshalError
	}

	return result, nil
}

type IdentityAllowedClient struct {
	Identity string `gorm:"column:identity"`
	Email    string `gorm:"column:email"`
	Details  []byte `gorm:"column:details"`
}

func (view IdentityAllowedClient) TableName() string {
	return "identity_allowed_clients"
}

type ClientDisplayDetails struct {
	BaseUri     string `json:"base_uri"`
	ClientId    string `json:"client_id"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
}
type ContextClientDetails struct {
	Context     string                 `json:"context"`
	DisplayName string                 `json:"display_name"`
	Description string                 `json:"description"`
	Clients     []ClientDisplayDetails `json:"clients"`
}
type AllowedClientList []ContextClientDetails

func IdentityAllowedClients(identCode string) (AllowedClientList, error) {
	db := util.GetDb()

	var result AllowedClientList
	var allowedList IdentityAllowedClient

	err := db.Where("identity = ?", identCode).Find(&allowedList).Error
	if err != nil {
		return result, err
	}

	if unmarshalError := json.Unmarshal(allowedList.Details, &result); unmarshalError != nil {
		return result, unmarshalError
	}

	return result, nil
}
