package model

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/ricecake/osin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/square/go-jose.v2"

	"github.com/ricecake/karma_chameleon/util"
)

type SessionToken struct {
	Code      string    `gorm:"column:code;not null"`
	Identity  string    `gorm:"column:identity;not null"`
	UserAgent string    `gorm:"column:user_agent;not null"`
	IpAddress string    `gorm:"column:ip_address;not null"`
	CreatedAt time.Time `gorm:"column:created_at;not null"`
	ExpiresIn int       `gorm:"column:expires_in;not null"`
}

func (this SessionToken) TableName() string {
	return "session_token"
}

func CreateSessionToken(tok *SessionToken) error {
	db := util.GetDb()

	if tok.Code == "" {
		tok.Code = util.CompactUUID()
	}
	log.Info("Creating session ", tok.Code)

	return db.Create(tok).Error
}

func ReplaceSessionToken(sessid, newSessid string) error {
	db := util.GetDb()
	tx := db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := db.Model(&AccessContext{}).Where("session = ?", sessid).Update("session", newSessid).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := db.Where("code = ?", sessid).Delete(&SessionToken{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	revocationEntry := RevocationEntry{
		EntityCode: sessid,
		CreatedAt:  time.Now(),
		ExpiresIn:  int((time.Duration(viper.GetInt("identity.ttl")) * time.Hour).Seconds()),
	}

	if err := db.Create(&revocationEntry).Error; err != nil {
		tx.Rollback()
		return err
	}

	log.Info("Revoking session ", revocationEntry.EntityCode)
	return tx.Commit().Error
}

type AccessContext struct {
	Code      string    `gorm:"column:code;not null"`
	Session   *string   `gorm:"column:session;not null"`
	Client    string    `gorm:"column:client;not null"`
	CreatedAt time.Time `gorm:"column:created_at;not null"`
}

func (this AccessContext) TableName() string {
	return "access_context"
}

func EnsureAccessContext(con *AccessContext) error {
	db := util.GetDb()

	if db.Where("session = ? AND client = ?", con.Session, con.Client).Find(&con).RecordNotFound() {
		log.Info("Creating context ", con.Code)
		con.Code = util.CompactUUID()
		return db.Create(&con).Error
	}

	return nil
}

type RevocationEntry struct {
	EntityCode string    `gorm:"column:entity_code;not null"`
	Field      string    `gorm:"column:field;not null; default:'jti'"`
	CreatedAt  time.Time `gorm:"column:created_at;not null"`
	ExpiresIn  int       `gorm:"column:expires_in;not null"`
}

func (this RevocationEntry) TableName() string {
	return "revocation"
}

func InsertRevocation(entity string, ttl int) error {
	db := util.GetDb()

	revocationEntry := RevocationEntry{
		EntityCode: entity,
		CreatedAt:  time.Now(),
		ExpiresIn:  ttl,
	}

	log.Info("Revoking entity ", revocationEntry.EntityCode)

	return db.Set("gorm:insert_option", "ON CONFLICT DO NOTHING").Create(&revocationEntry).Error
}

func EntityRevoked(entity string) bool {
	db := util.GetDb()
	return !db.Where("entity_code = ?", entity).First(&RevocationEntry{}).RecordNotFound()
}

func ListRevocations() (*util.RevMap, error) {
	db := util.GetDb()
	var results []RevocationEntry

	err := db.Find(&results).Error
	if err != nil {
		return nil, err
	}

	revMap := util.NewRevMap()
	for _, rev := range results {
		revMap.Add(rev.Field, rev.EntityCode, int(rev.CreatedAt.Unix()), rev.ExpiresIn)
	}

	return revMap, nil
}

type StashToken struct {
	UUID      string    `gorm:"column:uuid;       not null" json:"-"`
	Data      []byte    `gorm:"column:data;       not null" json:"-"`
	CreatedAt time.Time `gorm:"column:created_at; not null" json:"-"`
	ExpiresIn int       `gorm:"column:expires_in; not null" json:"-"`
}

func (StashToken) TableName() string {
	return "stash_data"
}

func Stash(data interface{}) (string, error) {
	return StashTTL(data, viper.GetInt("stash_ttl"))
}

func StashTTL(data interface{}, ttl int) (string, error) {
	db := util.GetDb()

	encData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	stash := StashToken{
		UUID:      util.CompactUUID(),
		Data:      encData,
		CreatedAt: time.Now(),
		ExpiresIn: ttl,
	}

	if stashErr := db.Create(&stash).Error; stashErr != nil {
		return "", err
	}
	return stash.UUID, nil
}

func StashFetch(uuid string, data interface{}) error {
	db := util.GetDb()
	tx := db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return err
	}

	var stash StashToken

	if tx.Where("uuid = ?", uuid).Find(&stash).RecordNotFound() {
		tx.Rollback()
		return fmt.Errorf("Invalid State Code")
	}

	deleteErr := tx.Where("uuid = ?", uuid).Delete(StashToken{}).Error
	if deleteErr != nil {
		tx.Rollback()
		return fmt.Errorf("stash Error: %s", deleteErr)
	}

	if stash.CreatedAt.Add(time.Second * time.Duration(stash.ExpiresIn)).Before(time.Now()) {
		tx.Rollback()
		return fmt.Errorf("bad stash")
	}

	if unmarshalError := json.Unmarshal(stash.Data, data); unmarshalError != nil {
		tx.Rollback()
		return unmarshalError
	}

	return tx.Commit().Error
}

type ZipCode struct {
	Identity    string
	Client      string
	Code        string
	TTL         int
	RedirectUri string
	Params      map[string]string
}

func (zip *ZipCode) Save() error {
	idp, clientErr := FindClientById(zip.Client)
	if clientErr != nil {
		return clientErr
	}

	redirectURL := idp.BaseUri
	if zip.RedirectUri != "" {
		redirectBase, err := url.Parse(idp.BaseUri)
		if err != nil {
			return err
		}

		redirect, redirErr := url.Parse(zip.RedirectUri)
		if redirErr != nil {
			return redirErr
		}

		redirect.Scheme = redirectBase.Scheme
		redirect.Host = redirectBase.Host

		baseQuery := redirect.Query()
		for key, value := range zip.Params {
			baseQuery.Add(key, value)
		}
		redirect.RawQuery = baseQuery.Encode()

		redirectURL = redirect.String()
	}

	zip.RedirectUri = redirectURL

	code, zipErr := StashTTL(zip, zip.TTL)
	zip.Code = code
	return zipErr
}

func FetchZipCode(code string) (zip ZipCode, zipErr error) {
	zipErr = StashFetch(code, &zip)
	return
}

type TokenGenerator struct{}

type UserAuthDetails struct {
	Code          string
	Nonce         string
	Browser       string
	Context       string
	Strength      string
	Method        string
	ValidResource []string
	Permitted     []string
}
type AuthCodeData struct {
	Code                string
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
	codeData := AuthCodeData{
		Code:                util.CompactUUID(),
		ClientId:            data.Client.GetId(),
		Scope:               data.Scope,
		RedirectUri:         data.RedirectUri,
		State:               data.State,
		CreatedAt:           data.CreatedAt,
		ExpiresIn:           data.ExpiresIn,
		CodeChallenge:       data.CodeChallenge,
		CodeChallengeMethod: data.CodeChallengeMethod,
	}
	if data.UserData != nil {
		codeData.UserData = data.UserData.(*UserAuthDetails)
	}

	return util.EncodeJWTClose(codeData, viper.GetString("security.passphrase"))
}

type AccessToken struct {
	Issuer     string `json:"iss"`
	UserCode   string `json:"sub"`
	Expiration int64  `json:"exp"`
	IssuedAt   int64  `json:"iat"`
	Code       string `json:"jti"`
	ClientId   string `json:"azp"`

	Nonce         string   `json:"nonce,omitempty"` // Non-manditory fields MUST be "omitempty"
	ValidResource []string `json:"aud,omitempty"`
	ContextCode   string   `json:"ctx,omitempty"`
	Scope         string   `json:"scope,omitempty"`
	Permitted     []string `json:"perm,omitempty"`

	Browser  string `json:"bro,omitempty"`
	Strength string `json:"acr,omitempty"`
	Method   string `json:"amr,omitempty"`
}
type RefreshToken struct {
	Code        string      `json:"jti"`
	AccessToken AccessToken `json:"ati"`

	ExpiresIn   int32
	Scope       string
	RedirectUri string
	CreatedAt   time.Time

	UserData *UserAuthDetails
}

func (a *TokenGenerator) GenerateAccessToken(data *osin.AccessData, generaterefresh bool) (accessToken string, refreshToken string, err error) {
	accessTokenData := AccessToken{
		Issuer:     viper.GetString("identity.issuer"),
		Expiration: data.CreatedAt.Add(time.Duration(data.ExpiresIn) * time.Second).Unix(),
		IssuedAt:   data.CreatedAt.Unix(),
		Code:       util.CompactUUID(),
		Scope:      data.Scope,
		ClientId:   data.Client.GetId(),
	}

	if data.UserData != nil {
		authDetails := data.UserData.(*UserAuthDetails)
		accessTokenData.Nonce = authDetails.Nonce
		accessTokenData.ContextCode = authDetails.Context
		accessTokenData.UserCode = authDetails.Code
		accessTokenData.Permitted = authDetails.Permitted
		accessTokenData.Browser = authDetails.Browser
		accessTokenData.Strength = authDetails.Strength
		accessTokenData.Method = authDetails.Method
		accessTokenData.ValidResource = authDetails.ValidResource
	}

	accessToken, err = util.EncodeJWTOpen(accessTokenData)
	if err != nil {
		return
	}

	if generaterefresh {
		plainRefreshToken := RefreshToken{
			Code:        util.CompactUUID(),
			AccessToken: accessTokenData,
			ExpiresIn:   data.ExpiresIn,
			Scope:       data.Scope,
			RedirectUri: data.RedirectUri,
			CreatedAt:   data.CreatedAt,
		}

		if data.UserData != nil {
			plainRefreshToken.UserData = data.UserData.(*UserAuthDetails)
		}

		refreshToken, err = util.EncodeJWTClose(plainRefreshToken, viper.GetString("security.passphrase"))
	}
	return
}

type IDToken struct {
	Issuer     string `json:"iss"`
	UserCode   string `json:"sub"`
	ClientID   string `json:"aud"`
	IssuedAt   int64  `json:"iat"`
	Expiration int64  `json:"exp"`
	TokenId    string `json:"jti"`

	Nonce    string `json:"nonce,omitempty"` // Non-manditory fields MUST be "omitempty"
	Strength string `json:"acr,omitempty"`
	Methods  string `json:"amr,omitempty"`
	Context  string `json:"ctx,omitempty"`

	// Custom claims supported by this server.
	//
	// See: https://openid.net/specs/openid-connect-core-1_0.html#StandardClaims

	Email         string `json:"email,omitempty"`
	FamilyName    string `json:"family_name,omitempty"`
	GivenName     string `json:"given_name,omitempty"`
	PreferredName string `json:"preferred_name,omitempty"`
}

func Cleanup() {
	db := util.GetDb()

	simpleRecordTypes := []interface{}{
		StashToken{},
		SessionToken{},
		RevocationEntry{},
	}
	for _, recType := range simpleRecordTypes {
		db.Where("created_at + expires_in * interval '1 second' < now()").Delete(recType)
	}
}

type LocalVerifierCache struct {
}

func NewLocalVerifierCache() (newCacher *LocalVerifierCache) {
	newCacher = &LocalVerifierCache{}
	return newCacher
}

func (verifier *LocalVerifierCache) Fetch() (*jose.JSONWebKeySet, *util.RevMap, error) {
	revoked, revErr := ListRevocations()
	if revErr != nil {
		return nil, nil, revErr
	}
	return util.Keys, revoked, nil
}
