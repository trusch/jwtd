package db

import (
	"errors"

	"gopkg.in/mgo.v2/bson"
)

type Group struct {
	ID     bson.ObjectId `bson:"_id,omitempty"`
	Name   string
	Rights []*AccessRight
}

type AccessRight struct {
	Service   string
	Subject string
}

func (db *DB) CreateGroup(name string, rights []*AccessRight) error {
	group := &Group{
		Name:   name,
		Rights: rights,
	}
	c := db.session.DB("jwtd").C("groups")
	n, err := c.Find(bson.M{"name": name}).Count()
	if n != 0 || err != nil {
		return errors.New("group already exists")
	}
	return c.Insert(group)
}

func (db *DB) GetGroup(name string) (*Group, error) {
	c := db.session.DB("jwtd").C("groups")
	group := &Group{}
	err := c.Find(bson.M{"name": name}).One(group)
	if err != nil {
		return nil, err
	}
	return group, nil
}

func (db *DB) DelGroup(name string) error {
	c := db.session.DB("jwtd").C("groups")
	return c.Remove(bson.M{"name": name})
}

func (db *DB) UpdateGroup(group *Group) error {
	c := db.session.DB("jwtd").C("groups")
	_, err := c.UpsertId(group.ID, group)
	return err
}

func (db *DB) ListGroups() ([]*Group, error) {
	c := db.session.DB("jwtd").C("groups")
	var groups []*Group
	err := c.Find(nil).All(&groups)
	return groups, err
}
