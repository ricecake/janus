package util

import (
	"encoding/base64"
	"github.com/pborman/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Hashing, encrypting, and strong random generation

func PasswordHash(plaintext []byte) ([]byte, error) {
	ciphertext, err := bcrypt.GenerateFromPassword(plaintext, 10)
	if err != nil {
		return nil, err
	}

	return ciphertext, nil
}

func PasswordHashValid(plaintext, ciphertext []byte) bool {
	err := bcrypt.CompareHashAndPassword(ciphertext, plaintext)

	return err == nil
}

func CompactUUID() string {
	return base64.RawURLEncoding.EncodeToString([]byte(uuid.NewRandom()))
}
