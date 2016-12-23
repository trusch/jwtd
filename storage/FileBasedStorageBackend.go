package storage

import (
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type FileBasedStorageBackend struct {
	ConfigDir string
}

func (backend *FileBasedStorageBackend) Load(project string) (*ProjectConfig, error) {
	bs, err := ioutil.ReadFile(filepath.Join(backend.ConfigDir, project+".yaml"))
	if err != nil {
		return nil, err
	}
	cfg := &ProjectConfig{}
	err = yaml.Unmarshal(bs, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func (backend *FileBasedStorageBackend) Save(project string, cfg *ProjectConfig) error {
	bs, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(backend.ConfigDir, project+".yaml"), bs, 0644)
}
