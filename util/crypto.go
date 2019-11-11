package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"github.com/pborman/uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/sha3"
	"io"
)

// Hashing, encrypting, and strong random generation

func PasswordHash(plaintext []byte) ([]byte, error) {
	ciphertext, err := bcrypt.GenerateFromPassword(plaintext, 10)
	if err != nil {
		return nil, err
	}

	return ciphertext, nil
}

func PasswordHashValid(plaintext, hash []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, plaintext)

	return err == nil
}

func CompactUUID() string {
	return base64.RawURLEncoding.EncodeToString([]byte(uuid.NewRandom()))
}

func Encrypt(passphrase string, plaintext []byte) ([]byte, error) {
	key := DeriveKey(passphrase)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)
	return append(nonce, ciphertext...), nil
}

func Decrypt(passphrase string, ciphertext []byte) ([]byte, error) {
	key := DeriveKey(passphrase)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err := aesgcm.Open(nil, ciphertext[:aesgcm.NonceSize()], ciphertext[aesgcm.NonceSize():], nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

func DeriveKey(passphrase string) []byte {
	return pbkdf2.Key([]byte(passphrase), []byte("0b131f2dc50d836216e2f25453cd012b131cd773"), 100, 32, sha3.New256)
}
