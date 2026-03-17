package services

import (
	"crypto/pbkdf2"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"hash"
)

type PasswordManager struct {
	passwordIterations int
	passwordKeyLength  int
	hashAlgorithm func() hash.Hash
}

func NewPasswordManager(passwordIterations int, passwordKeyLength int, hashAlgorithm func() hash.Hash) *PasswordManager {
	return &PasswordManager{
		passwordIterations: passwordIterations,
		passwordKeyLength:  passwordKeyLength,
		hashAlgorithm: hashAlgorithm,
	}
}

func (pm *PasswordManager) GenerateSalt(length int) (string, error) {
	salt := make([]byte, length)
	
	_, err := rand.Read(salt)
	if err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	return hex.EncodeToString(salt), nil
}

func (pm *PasswordManager) HashPassword(password string, salt string) (string, error) {
	derivedKey, err := pbkdf2.Key(pm.hashAlgorithm, password, []byte(salt), pm.passwordIterations, pm.passwordKeyLength)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return hex.EncodeToString(derivedKey), nil
}