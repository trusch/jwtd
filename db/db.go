package db

import "gopkg.in/mgo.v2"

type DB struct {
	session *mgo.Session
}

func New(url string) (*DB, error) {
	session, err := mgo.Dial(url)
	if err != nil {
		return nil, err
	}
	return &DB{session}, nil
}
