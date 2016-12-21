package db

import "errors"

type Group struct {
	Name    string                       `json:"name"`
	Project string                       `json:"project"`
	Rights  map[string]map[string]string `json:"rights"`
}

func (db *DB) CreateGroup(project, name string, rights map[string]map[string]string) error {
	group := &Group{
		Name:    name,
		Project: project,
		Rights:  rights,
	}
	for _, g := range db.Config.Groups {
		if g.Name == name && g.Project == project {
			return errors.New("group already exists")
		}
	}
	db.Config.Groups = append(db.Config.Groups, group)
	return db.Config.Save(db.ConfigPath)
}

func (db *DB) GetGroup(project, name string) (*Group, error) {
	for _, g := range db.Config.Groups {
		if g.Project == project && g.Name == name {
			return g, nil
		}
	}
	return nil, errors.New("group not found")
}

func (db *DB) DelGroup(project, name string) error {
	for idx, g := range db.Config.Groups {
		if g.Project == project && g.Name == name {
			db.Config.Groups = append(db.Config.Groups[:idx], db.Config.Groups[idx+1:]...)
			return db.Config.Save(db.ConfigPath)
		}
	}
	return errors.New("group not found")
}

func (db *DB) UpdateGroup(group *Group) error {
	err := db.DelGroup(group.Project, group.Name)
	if err != nil {
		return err
	}
	db.Config.Groups = append(db.Config.Groups, group)
	return db.Config.Save(db.ConfigPath)
}

func (db *DB) ListGroups(project string) ([]*Group, error) {
	groups := []*Group{}
	for _, g := range db.Config.Groups {
		if g.Project == project {
			groups = append(groups, g)
		}
	}
	return groups, nil
}
