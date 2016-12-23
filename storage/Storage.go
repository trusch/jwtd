package storage

import "errors"

type Storage struct {
	projects map[string]*ProjectConfig
	backend  StorageBackend
}

func New(backend StorageBackend) *Storage {
	return &Storage{make(map[string]*ProjectConfig), backend}
}

func (storage *Storage) GetProjectConfig(project string) (*ProjectConfig, error) {
	if cfg, ok := storage.projects[project]; ok {
		return cfg, nil
	}
	if cfg, err := storage.backend.Load(project); err == nil {
		storage.projects[project] = cfg
		return cfg, nil
	}
	return nil, errors.New("no such project")
}

func (storage *Storage) Reset() {
	storage.projects = make(map[string]*ProjectConfig)
}

func (storage *Storage) CreateProject(project string) {
	storage.projects[project] = &ProjectConfig{}
}
