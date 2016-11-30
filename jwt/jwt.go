package jwt

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/dgrijalva/jwt-go"
)

type Claims map[string]interface{}

func CreateToken(claims Claims, key *rsa.PrivateKey) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.MapClaims(claims))
	return token.SignedString(key)
}

func ValidateToken(tokenString string, key *rsa.PublicKey) (Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return Claims(claims), nil
	}
	return nil, errors.New("invalid token")
}

func LoadPublicKey(keyFile string) (*rsa.PublicKey, error) {
	bs, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}
	return jwt.ParseRSAPublicKeyFromPEM(bs)
}

func LoadPrivateKey(keyFile string) (*rsa.PrivateKey, error) {
	bs, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}
	return jwt.ParseRSAPrivateKeyFromPEM(bs)
}
