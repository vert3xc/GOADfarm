package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateSecretKey() (string, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(key), nil
}

func InitSecrets() (string, string, error){
	sessionSecret, err := GenerateSecretKey()
	if err != nil {
		return "", "", err
	}
	apiKey, err := GenerateSecretKey()
	if err != nil {
		return "", "", err
	}
	return sessionSecret, apiKey, nil
}