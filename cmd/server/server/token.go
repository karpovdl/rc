package main

import (
	"errors"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
)

func getToken(r *http.Request) (*jwt.Token, error) {
	var (
		token *jwt.Token
		err   error
	)

	hashSecretGetter := func(token *jwt.Token) (interface{}, error) {
		if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok || method.Alg() != "HS256" {
			return nil, errors.New("bad sign method")
		}
		return tokenSecret, nil
	}

	inToken := r.Header.Get(tokenName)
	if token, err = jwt.Parse(inToken, hashSecretGetter); err != nil || !token.Valid {
		return nil, err
	}

	return token, nil
}
