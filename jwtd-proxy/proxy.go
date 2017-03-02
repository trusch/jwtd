package main

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"
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
		singleProxy, err := NewSingleProxy(hostCfg.Project, host, hostCfg.Backend, hostCfg.Routes, proxy.key)
		if err != nil {
			return nil, err
		}
		var handler http.Handler = singleProxy
		if hostCfg.CORS != nil {
			headers := handlers.AllowedHeaders(hostCfg.CORS.AllowedHeaders)
			origins := handlers.AllowedOrigins(hostCfg.CORS.AllowedOrigins)
			methods := handlers.AllowedMethods(hostCfg.CORS.AllowedMethods)
			corsWrapper := handlers.CORS(headers, origins, methods)
			handler = corsWrapper(handler)
		}
		r.Host(host).Handler(handler)
	}
	proxy.router = r
	return proxy, nil
}

func (proxy *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Print("incoming request for host: ", r.Host)
	proxy.router.ServeHTTP(w, r)
}
