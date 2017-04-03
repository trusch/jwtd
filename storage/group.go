package storage

import "errors"

func (storage *Storage) CreateGroup(name string, rights map[string]map[string]string) error {
	projectConfig, err := storage.GetProjectConfig()
	if err != nil {
		return err
	}
	group := &Group{
		Name:   name,
		Rights: rights,
	}
	for _, g := range projectConfig.Groups {
		if g.Name == name {
			return errors.New("group already exists")
		}
	}
	projectConfig.Groups = append(projectConfig.Groups, group)
	return storage.backend.Save(projectConfig)
}

func (storage *Storage) GetGroup(name string) (*Group, error) {
	projectConfig, err := storage.GetProjectConfig()
	if err != nil {
		return nil, err
	}
	for _, g := range projectConfig.Groups {
		if g.Name == name {
			return g, nil
		}
	}
	return nil, errors.New("group not found")
}

func (storage *Storage) DelGroup(name string) error {
	projectConfig, err := storage.GetProjectConfig()
	if err != nil {
		return err
	}
	for idx, g := range projectConfig.Groups {
		if g.Name == name {
			projectConfig.Groups = append(projectConfig.Groups[:idx], projectConfig.Groups[idx+1:]...)
			return storage.backend.Save(projectConfig)
		}
	}
	return errors.New("group not found")
}

func (storage *Storage) UpdateGroup(group *Group) error {
	projectConfig, err := storage.GetProjectConfig()
	if err != nil {
		return err
	}
	err = storage.DelGroup(group.Name)
	if err != nil {
		return err
	}
	projectConfig.Groups = append(projectConfig.Groups, group)
	return storage.backend.Save(projectConfig)
}

func (storage *Storage) ListGroups() ([]*Group, error) {
	projectConfig, err := storage.GetProjectConfig()
	if err != nil {
		return nil, err
	}
	return projectConfig.Groups, nil
}
