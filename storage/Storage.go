package storage

import "errors"

type Storage struct {
	cfg     *ProjectConfig
	backend StorageBackend
}

func New(backend StorageBackend) *Storage {
	return &Storage{nil, backend}
}

func (storage *Storage) GetProjectConfig() (*ProjectConfig, error) {
	if storage.cfg != nil {
		return storage.cfg, nil
	}
	if cfg, err := storage.backend.Load(); err == nil {
		storage.cfg = cfg
		return cfg, nil
	}
	return nil, errors.New("config not found")
}

func (storage *Storage) Reset() {
	storage.cfg = nil
}
