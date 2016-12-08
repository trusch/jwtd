package jwt

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/dgrijalva/jwt-go"
)

type Claims map[string]interface{}

func CreateToken(claims Claims, key interface{}) (string, error) {
	switch k := key.(type) {
	case *rsa.PrivateKey:
		{
			token := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.MapClaims(claims))
			return token.SignedString(k)
		}
	case *ecdsa.PrivateKey:
		{
			token := jwt.NewWithClaims(jwt.SigningMethodES512, jwt.MapClaims(claims))
			return token.SignedString(k)
		}
	}
	return "", errors.New("invalid private key")
}

func ValidateToken(tokenString string, key interface{}) (Claims, error) {
	var (
		token *jwt.Token
		err   error
	)
	switch k := key.(type) {
	case *rsa.PublicKey:
		{
			token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}
				return k, nil
			})
		}
	case *ecdsa.PublicKey:
		{
			token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}
				return k, nil
			})
		}
	}
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return Claims(claims), nil
	}
	return nil, errors.New("invalid token")
}

func LoadPublicKey(keyFile string) (interface{}, error) {
	bs, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}
	rsaKey, err := jwt.ParseRSAPublicKeyFromPEM(bs)
	if err != nil {
		ecKey, err := jwt.ParseECPublicKeyFromPEM(bs)
		if err != nil {
			return nil, errors.New("unknown public key type")
		}
		return ecKey, nil
	}
	return rsaKey, nil
}

func LoadPrivateKey(keyFile string) (interface{}, error) {
	bs, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}
	rsaKey, err := jwt.ParseRSAPrivateKeyFromPEM(bs)
	if err != nil {
		ecKey, err := jwt.ParseECPrivateKeyFromPEM(bs)
		if err != nil {
			return nil, errors.New("unknown public key type")
		}
		return ecKey, nil
	}
	return rsaKey, nil
}
