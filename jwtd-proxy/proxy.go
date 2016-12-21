package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/trusch/jwtd/jwt"
)

type Proxy struct {
	cfg     *Config
	proxies map[string]*SingleProxy
	router  *mux.Router
	key     interface{}
}

func NewProxy(cfg *Config) (*Proxy, error) {
	r := mux.NewRouter()
	key, err := jwt.LoadPublicKey(cfg.Cert)
	if err != nil {
		return nil, err
	}
	proxy := &Proxy{cfg: cfg, key: key}
	for host, hostCfg := range cfg.Hosts {
		singleProxy, err := NewSingleProxy(cfg.Project, host, hostCfg.Backend, hostCfg.Routes, proxy.key)
		if err != nil {
			return nil, err
		}
		r.Host(host).Handler(singleProxy)
	}
	proxy.router = r
	return proxy, nil
}

func (proxy *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	proxy.router.ServeHTTP(w, r)
}
