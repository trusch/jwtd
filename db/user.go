package db

import (
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	ID           bson.ObjectId `bson:"_id,omitempty"`
	Name         string
	PasswordHash []byte `yaml:"-"`
	Groups       []string
}

func (db *DB) CreateUser(name, password string, groups []string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := &User{
		Name:         name,
		PasswordHash: hash,
		Groups:       groups,
	}
	c := db.session.DB("jwtd").C("users")
	n, err := c.Find(bson.M{"name": name}).Count()
	if n != 0 || err != nil {
		return errors.New("user already exists")
	}
	return c.Insert(user)
}

func (db *DB) GetUser(name string) (*User, error) {
	c := db.session.DB("jwtd").C("users")
	user := &User{}
	err := c.Find(bson.M{"name": name}).One(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (db *DB) DelUser(name string) error {
	c := db.session.DB("jwtd").C("users")
	return c.Remove(bson.M{"name": name})
}

func (db *DB) UpdateUser(user *User) error {
	c := db.session.DB("jwtd").C("users")
	_, err := c.UpsertId(user.ID, user)
	if err != nil {
		log.Print("update fail")
	}
	return err
}

func (user *User) CheckRights(db *DB, service string, subject string) (bool, error) {
	for _, groupName := range user.Groups {
		group, err := db.GetGroup(groupName)
		if err != nil {
			log.Printf("Error loading group %v: %v", groupName, err)
			continue
		}
		for _, right := range group.Rights {
			if right.Service == service && right.Subject == subject {
				return true, nil
			}
		}
	}
	return false, nil
}

func (user *User) CheckPassword(password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password))
	if err != nil {
		return false, err
	}
	return true, nil
}

func (db *DB) ListUsers() ([]*User, error) {
	c := db.session.DB("jwtd").C("users")
	var users []*User
	err := c.Find(nil).All(&users)
	return users, err
}
