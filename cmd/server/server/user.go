package main

import (
	"errors"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	bson "gopkg.in/mgo.v2/bson"
)

//User type
type User struct {
	ID       bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Username string        `json:"username" bson:"username"`
	Password string        `json:"password" bson:"password"`
}

func (u *User) getJWTToken() (string, error) {
	if u == nil {
		return "", errors.New("user type not initialized")
	}

	var (
		tokenString string
		err         error
	)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": map[string]interface{}{
			"username": u.Username,
			"id":       u.ID,
		},
		"iat": time.Now().Add(time.Hour * 0).Unix(),  // iat = issued at
		"exp": time.Now().Add(time.Hour * 24).Unix(), // exp = expiration time
	})

	if tokenString, err = token.SignedString(tokenSecret); err != nil {
		return "", err
	}

	return tokenString, nil
}

func isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header[tokenName] != nil {
			if _, err := getToken(r); err != nil {
				jsonError(w, http.StatusUnauthorized, err.Error())
				return
			}

			endpoint(w, r)
		} else {
			jsonError(w, http.StatusUnauthorized, "not authorized")
			return
		}
	})
}
