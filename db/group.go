package db

import "errors"

type Group struct {
	Name   string
	Rights []*AccessRight
}

type AccessRight struct {
	Service string
	Subject string
}

func (db *DB) CreateGroup(name string, rights []*AccessRight) error {
	group := &Group{
		Name:   name,
		Rights: rights,
	}
	for _, g := range db.Config.Groups {
		if g.Name == name {
			return errors.New("group already exists")
		}
	}
	db.Config.Groups = append(db.Config.Groups, group)
	return db.Config.Save(db.ConfigPath)
}

func (db *DB) GetGroup(name string) (*Group, error) {
	for _, g := range db.Config.Groups {
		if g.Name == name {
			return g, nil
		}
	}
	return nil, errors.New("group not found")
}

func (db *DB) DelGroup(name string) error {
	for idx, g := range db.Config.Groups {
		if g.Name == name {
			db.Config.Groups = append(db.Config.Groups[:idx], db.Config.Groups[idx+1:]...)
			return db.Config.Save(db.ConfigPath)
		}
	}
	return errors.New("group not found")
}

func (db *DB) UpdateGroup(group *Group) error {
	err := db.DelGroup(group.Name)
	if err != nil {
		return err
	}
	db.Config.Groups = append(db.Config.Groups, group)
	return db.Config.Save(db.ConfigPath)
}

func (db *DB) ListGroups() ([]*Group, error) {
	return db.Config.Groups, nil
}
