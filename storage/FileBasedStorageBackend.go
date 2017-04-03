package storage

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type FileBasedStorageBackend struct {
	ConfigFile string
}

func (backend *FileBasedStorageBackend) Load() (*ProjectConfig, error) {
	bs, err := ioutil.ReadFile(backend.ConfigFile)
	if err != nil {
		return &ProjectConfig{}, nil
	}
	cfg := &ProjectConfig{}
	err = yaml.Unmarshal(bs, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func (backend *FileBasedStorageBackend) Save(cfg *ProjectConfig) error {
	bs, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(backend.ConfigFile, bs, 0644)
}
