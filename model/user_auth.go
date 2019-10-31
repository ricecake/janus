package model

type UserAuth interface{}

type UserPassword struct {
	User         int
	PasswordHash []byte
	TotpSecret   *[]byte
}
