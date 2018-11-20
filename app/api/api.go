package api

import (
	"io/ioutil"
	"github.com/prometheus/common/log"
	"crypto/rsa"
	"github.com/dgrijalva/jwt-go"

	"time"
)

type APIClaims struct {
	Username string `json:"username"`
	UserType string `json:"user_type"`
}

func SetupAPIKeys() (*rsa.PrivateKey, *rsa.PublicKey) {
	var err error

	signBytes, err := ioutil.ReadFile("app/keys/app.rsa")
	if err != nil {
		log.Fatal(err.Error())
	}

	signKey, _ := jwt.ParseRSAPrivateKeyFromPEM(signBytes)

	verifyBytes, err := ioutil.ReadFile("app/keys/app.rsa.pub")
	if err != nil {
		log.Fatal(err.Error())
	}

	verifyKey, _ := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)

	return signKey, verifyKey
}

func GenerateToken(username, userType string, validity int, signkey *rsa.PrivateKey) (string, error) {
	claims := &jwt.MapClaims{"exp": time.Now().Add(time.Minute*time.Duration(validity)), "username": username, "user_type": userType}
	tok := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	tokString, err := tok.SignedString(signkey)
	if err != nil {
		return "", err
	}
	return tokString, nil
}