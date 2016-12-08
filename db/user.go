package db

import (
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Name         string
	PasswordHash string
	Groups       []string
}

func (db *DB) CreateUser(name, password string, groups []string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := &User{
		Name:         name,
		PasswordHash: string(hash),
		Groups:       groups,
	}
	for _, u := range db.Config.Users {
		if u.Name == name {
			return errors.New("user already exists")
		}
	}
	db.Config.Users = append(db.Config.Users, user)
	return db.Config.Save(db.ConfigPath)
}

func (db *DB) GetUser(name string) (*User, error) {
	for _, u := range db.Config.Users {
		if u.Name == name {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (db *DB) DelUser(name string) error {
	for idx, user := range db.Config.Users {
		if user.Name == name {
			db.Config.Users = append(db.Config.Users[:idx], db.Config.Users[idx+1:]...)
			return db.Config.Save(db.ConfigPath)
		}
	}
	return errors.New("group not found")
}

func (db *DB) UpdateUser(user *User) error {
	err := db.DelUser(user.Name)
	if err != nil {
		return err
	}
	db.Config.Users = append(db.Config.Users, user)
	return db.Config.Save(db.ConfigPath)
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
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return false, err
	}
	return true, nil
}

func (db *DB) ListUsers() ([]*User, error) {
	return db.Config.Users, nil
}
