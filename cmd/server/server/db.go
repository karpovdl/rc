package main

import (
	"fmt"

	mgo "gopkg.in/mgo.v2"
)

func connect() (*mgo.Database, error) {
	//host := "mongodb://localhost:27017"
	host := "mongo:27017"
	dbName := "rc"

	if session, err := mgo.Dial(host); err != nil {
		return nil, err
	} else {
		if err := session.Ping(); err != nil {
			return nil, err
		}

		db := session.DB(dbName)

		fmt.Println("Connected to MongoDB!")

		return db, nil
	}
}
