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
			log.Printf("use CORS for %v", host)
			log.Printf("allowed Headers: %v", hostCfg.CORS.Headers)
			log.Printf("allowed Origins: %v", hostCfg.CORS.Origins)
			log.Printf("allowed Methods: %v", hostCfg.CORS.Methods)
			headers := handlers.AllowedHeaders(hostCfg.CORS.Headers)
			origins := handlers.AllowedOrigins(hostCfg.CORS.Origins)
			methods := handlers.AllowedMethods(hostCfg.CORS.Methods)
			corsWrapper := handlers.CORS(headers, origins, methods)
			handler = corsWrapper(handler)
		}
		r.Host(host).Handler(handler)
	}
	proxy.router = r
	return proxy, nil
}

func (proxy *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	proxy.router.ServeHTTP(w, r)
}
