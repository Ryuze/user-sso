package util

import (
	"crypto/x509"
	"encoding/pem"
	"os"
	"time"

	database "github.com/ideal-tekno-solusi/sso/database/main"
	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwe"
	"github.com/lestrrat-go/jwx/v3/jwt"
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

func BuildUserJwt(user database.GetUserRow) (*jwt.Token, *time.Duration, error) {
	expiryTime := time.Minute * 20

	token, err := jwt.NewBuilder().
		Expiration(time.Now().Add(expiryTime)).
		Build()
	if err != nil {
		logrus.Fatalf("failed to build token with error: %v", err)
		return nil, nil, err
	}

	token.Set("id", user.ID)
	token.Set("username", user.Username)
	token.Set("name", user.Name)
	token.Set("services", user.AllowedServices.String)

	return &token, &expiryTime, nil
}

func SignJwt(token jwt.Token) (*string, error) {
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

	sign, err := jwt.Sign(token, jwt.WithKey(jwa.ES256(), ecKey))
	if err != nil {
		logrus.Fatalf("failed to sign token with error: %v", err)
		return nil, err
	}

	result = string(sign)

	return &result, nil
}

func VerifyJwt(token, pubKey string) (jwt.Token, error) {
	block, _ := pem.Decode([]byte(pubKey))

	ecKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		logrus.Fatalf("failed to parse key with error: %v", err)
		return nil, err
	}

	verifiedToken, err := jwt.Parse([]byte(token), jwt.WithKey(jwa.ES256(), ecKey))
	if err != nil {
		logrus.Fatalf("failed to parse token with error: %v", err)
		return nil, err
	}

	return verifiedToken, nil
}
