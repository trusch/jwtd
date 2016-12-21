package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Listen string
	Cert   string
	Hosts  map[string]*Host
}

type Host struct {
	Backend string
	TLS     *TLSConfig
	Routes  []*Route
}

type TLSConfig struct {
	Cert string
	Key  string
}

type Route struct {
	Path    string
	Require map[string]string
}

func NewConfigFromFile(path string) (*Config, error) {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cfg := &Config{}
	err = yaml.Unmarshal(bs, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func (cfg *Config) String() string {
	bs, _ := yaml.Marshal(cfg)
	return string(bs)
}

func (cfg *Config) GetCerts() []*TLSConfig {
	res := make([]*TLSConfig, 0, len(cfg.Hosts))
	for _, hostCfg := range cfg.Hosts {
		res = append(res, hostCfg.TLS)
	}
	return res
}
