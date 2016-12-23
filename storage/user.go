package storage

import (
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func (storage *Storage) CreateUser(project, name, password string, groups []string) error {
	projectConfig, err := storage.GetProjectConfig(project)
	if err != nil {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := &User{
		Name:         name,
		PasswordHash: string(hash),
		Groups:       groups,
	}
	for _, u := range projectConfig.Users {
		if u.Name == name {
			return errors.New("user already exists")
		}
	}
	projectConfig.Users = append(projectConfig.Users, user)
	return storage.backend.Save(project, projectConfig)
}

func (storage *Storage) GetUser(project, name string) (*User, error) {
	projectConfig, err := storage.GetProjectConfig(project)
	if err != nil {
		return nil, err
	}
	for _, u := range projectConfig.Users {
		if u.Name == name {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (storage *Storage) DelUser(project, name string) error {
	projectConfig, err := storage.GetProjectConfig(project)
	if err != nil {
		return err
	}
	for idx, user := range projectConfig.Users {
		if user.Name == name {
			projectConfig.Users = append(projectConfig.Users[:idx], projectConfig.Users[idx+1:]...)
			return storage.backend.Save(project, projectConfig)
		}
	}
	return errors.New("user not found")
}

func (storage *Storage) UpdateUser(project string, user *User) error {
	projectConfig, err := storage.GetProjectConfig(project)
	if err != nil {
		return err
	}
	err = storage.DelUser(project, user.Name)
	if err != nil {
		return err
	}
	projectConfig.Users = append(projectConfig.Users, user)
	return storage.backend.Save(project, projectConfig)
}

func (user *User) CheckRights(storage *Storage, project, service string, labels map[string]string) (bool, error) {
	toAck := len(labels)
	for _, groupName := range user.Groups {
		group, err := storage.GetGroup(project, groupName)
		if err != nil {
			log.Printf("Error loading group %v: %v", groupName, err)
			continue
		}
		for requestedKey, requestedValue := range labels {
			for key, value := range group.Rights[service] {
				if (key == "*" || key == requestedKey) && (value == "*" || value == requestedValue) {
					toAck -= 1
					break
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

func (storage *Storage) ListUsers(project string) ([]*User, error) {
	projectConfig, err := storage.GetProjectConfig(project)
	if err != nil {
		return nil, err
	}
	return projectConfig.Users, nil
}
