package main

import (
	"errors"
	"net/http"

	mgo "gopkg.in/mgo.v2"
	bson "gopkg.in/mgo.v2/bson"
)

//UserModel type
type UserModel struct {
	DB *mgo.Database
}

func (userModel *UserModel) register(user *User) (int, error) {
	if userModel == nil || userModel.DB == nil {
		return http.StatusInternalServerError, errors.New("user model type not initialized")
	}

	if user == nil {
		return http.StatusInternalServerError, errors.New("user not initialized")
	}

	if count, err := userModel.DB.C("users").Find(
		bson.M{"username": user.Username}).Count(); count > 0 || err != nil {
		if count > 0 {
			return http.StatusInternalServerError, errors.New("already exists")
		}
		return http.StatusInternalServerError, err
	}

	user.ID = bson.NewObjectId()
	user.Password, _ = hashPassword(user.Password)

	if err := userModel.DB.C("users").Insert(user); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (userModel *UserModel) login(user *User) (int, error) {
	if userModel == nil || userModel.DB == nil {
		return http.StatusInternalServerError, errors.New("user model type not initialized")
	}

	if user == nil {
		return http.StatusInternalServerError, errors.New("user not initialized")
	}

	if count, err := userModel.DB.C("users").Find(bson.M{
		"username": user.Username,
	}).Count(); count == 0 || err != nil {
		if count == 0 {
			return http.StatusNotFound, errors.New("user not found")
		} else if count > 1 {
			return http.StatusInternalServerError, errors.New("user many found")
		}
		return http.StatusInternalServerError, err
	}

	var u User
	if err := userModel.DB.C("users").Find(
		bson.M{"username": user.Username}).One(&u); err != nil {
		return http.StatusInternalServerError, err
	}

	if match := checkPasswordHash(user.Password, u.Password); !match {
		return http.StatusInternalServerError, errors.New("invalid password")
	}

	*user = *&u

	return http.StatusOK, nil
}
