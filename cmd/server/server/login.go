package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	mgo "gopkg.in/mgo.v2"
)

func loginUser(w http.ResponseWriter, r *http.Request) {
	var (
		body        []byte
		resp        []byte
		db          *mgo.Database
		tokenString string
		err         error
	)

	if r.Header.Get("Content-Type") != "application/json" {
		jsonError(w, http.StatusBadRequest, "unknown payload")
		return
	}

	if body, err = ioutil.ReadAll(r.Body); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer r.Body.Close()

	user := &User{}
	if err = json.Unmarshal(body, user); err != nil {
		jsonError(w, http.StatusBadRequest, "cant unpack payload")
		return
	}

	if user.Username == "" || user.Password == "" {
		jsonError(w, http.StatusUnauthorized, "bad login or password")
		return
	}

	if db, err = connect(); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer db.Session.Close()

	userModel := UserModel{
		DB: db,
	}

	if status, err := userModel.login(user); status != http.StatusOK {
		jsonMessage(w, status, err.Error())
		return
	}

	if tokenString, err = user.getJWTToken(); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if resp, err = json.Marshal(map[string]interface{}{
		"token": tokenString,
	}); err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
	w.Write([]byte("\n\n"))
}
