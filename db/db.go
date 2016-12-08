package db

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type DB struct {
	ConfigPath string
	Config     *ConfigFile
}

func New(path string) (*DB, error) {
	config := &ConfigFile{}
	err := config.Load(path)
	if err != nil {
		return nil, err
	}
	return &DB{path, config}, nil
}

type ConfigFile struct {
	Users  []*User
	Groups []*Group
}

func (config *ConfigFile) Load(path string) error {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(bs, config)
}

func (config *ConfigFile) Save(path string) error {
	bs, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, bs, 0655)
}
