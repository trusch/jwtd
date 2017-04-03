package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/trusch/jwtd/validator"
)

// SingleProxy wraps a single reverse proxy and guards it by the jwtd validator
type SingleProxy struct {
	service   string
	router    *mux.Router
	proxy     *httputil.ReverseProxy
	validator *validator.Validator
}

// NewSingleProxy returns a new SingleProxy instance
func NewSingleProxy(serviceName, backend string, routes []*Route, jwtdCrt interface{}) (*SingleProxy, error) {
	backendURL, err := url.Parse(backend)
	if err != nil {
		return nil, err
	}
	validator, err := validator.New(jwtdCrt)
	if err != nil {
		return nil, err
	}
	proxy := &SingleProxy{
		service:   serviceName,
		proxy:     httputil.NewSingleHostReverseProxy(backendURL),
		validator: validator,
	}
	proxy.constructRouter(routes)
	return proxy, nil
}

func (proxy *SingleProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v: incoming request: %v", proxy.service, r.URL)
	proxy.router.ServeHTTP(w, r)
}

func (proxy *SingleProxy) constructRouter(routes []*Route) {
	r := mux.NewRouter()
	for _, route := range routes {
		sub := r.PathPrefix(route.Path)
		if len(route.Methods) > 0 {
			sub = sub.Methods(route.Methods...)
		}
		sub.HandlerFunc(proxy.buildHandler(route.Require))
	}
	proxy.router = r
}

func (proxy *SingleProxy) sendMessage(w http.ResponseWriter, msg string) {
	data := map[string]string{
		"message": msg,
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(data)
}

func (proxy *SingleProxy) buildHandler(required map[string]string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := proxy.validator.Validate(r, proxy.service, required); err != nil {
			proxy.sendMessage(w, err.Error())
			return
		}
		proxy.proxy.ServeHTTP(w, r)
	}
}
