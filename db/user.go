package db

import (
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Name         string   `json:"name"`
	PasswordHash string   `json:"passwordhash,omitempty"`
	Project      string   `json:"project"`
	Groups       []string `json:"groups"`
}

func (db *DB) CreateUser(project, name, password string, groups []string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := &User{
		Name:         name,
		Project:      project,
		PasswordHash: string(hash),
		Groups:       groups,
	}
	for _, u := range db.Config.Users {
		if u.Name == name && u.Project == project {
			return errors.New("user already exists")
		}
	}
	db.Config.Users = append(db.Config.Users, user)
	return db.Config.Save(db.ConfigPath)
}

func (db *DB) GetUser(project, name string) (*User, error) {
	for _, u := range db.Config.Users {
		if u.Project == project && u.Name == name {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (db *DB) DelUser(project, name string) error {
	for idx, user := range db.Config.Users {
		if user.Project == project && user.Name == name {
			db.Config.Users = append(db.Config.Users[:idx], db.Config.Users[idx+1:]...)
			return db.Config.Save(db.ConfigPath)
		}
	}
	return errors.New("user not found")
}

func (db *DB) UpdateUser(user *User) error {
	err := db.DelUser(user.Project, user.Name)
	if err != nil {
		return err
	}
	db.Config.Users = append(db.Config.Users, user)
	return db.Config.Save(db.ConfigPath)
}

func (user *User) CheckRights(db *DB, service string, labels map[string]string) (bool, error) {
	toAck := len(labels)
	for _, groupName := range user.Groups {
		group, err := db.GetGroup(user.Project, groupName)
		if err != nil {
			log.Printf("Error loading group %v: %v", groupName, err)
			continue
		}
		for key, value := range group.Rights[service] {
			for k, v := range labels {
				if key == k && value == v {
					toAck -= 1
				}
			}
		}
	}
	if toAck == 0 {
		return true, nil
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

func (db *DB) ListUsers(project string) ([]*User, error) {
	users := []*User{}
	for _, u := range db.Config.Users {
		if u.Project == project {
			users = append(users, u)
		}
	}
	return users, nil
}
