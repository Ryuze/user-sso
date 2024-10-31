package util

import (
	"crypto/x509"
	"encoding/pem"
	"os"

	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwe"
	"github.com/sirupsen/logrus"
)

func EncryptJwe(message string) (*string, error) {
	var result string

	data, err := os.ReadFile("/home/ryuze/projects/sso/secret/pubkey.pem")
	if err != nil {
		logrus.Fatalf("failed to read pubkey with error: %v", err)
		return nil, err
	}

	block, _ := pem.Decode(data)

	ecKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		logrus.Fatalf("failed to parse key with error: %v", err)
		return nil, err
	}

	encrypted, err := jwe.Encrypt([]byte(message), jwe.WithKey(jwa.ECDH_ES(), ecKey))
	if err != nil {
		logrus.Fatalf("failed to encrypt message with error: %v", err)
		return nil, err
	}

	result = string(encrypted)

	return &result, nil
}

func DecryptJwe(message string) (*string, error) {
	var result string

	data, err := os.ReadFile("/home/ryuze/projects/sso/secret/privkey.pem")
	if err != nil {
		logrus.Fatalf("failed to read privkey with error: %v", err)
		return nil, err
	}

	block, _ := pem.Decode(data)

	ecKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		logrus.Fatalf("failed to parse key with error: %v", err)
		return nil, err
	}

	decrypted, err := jwe.Decrypt([]byte(message), jwe.WithKey(jwa.ECDH_ES(), ecKey))
	if err != nil {
		logrus.Fatalf("failed to decrypt message with error: %v", err)
		return nil, err
	}

	result = string(decrypted)

	return &result, nil
}
