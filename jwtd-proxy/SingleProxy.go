package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/gorilla/mux"
	"github.com/trusch/jwtd/jwt"
)

type SingleProxy struct {
	service string
	project string
	router  *mux.Router
	proxy   *httputil.ReverseProxy
	jwtdCrt interface{}
}

func NewSingleProxy(project, serviceName, backend string, routes []*Route, jwtdCrt interface{}) (*SingleProxy, error) {
	backendUrl, err := url.Parse(backend)
	if err != nil {
		return nil, err
	}
	proxy := &SingleProxy{
		service: serviceName,
		project: project,
		proxy:   httputil.NewSingleHostReverseProxy(backendUrl),
		jwtdCrt: jwtdCrt,
	}
	proxy.constructRouter(routes)
	return proxy, nil
}

func (proxy *SingleProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("service: %v url: %v", proxy.service, r.URL)
	proxy.router.ServeHTTP(w, r)
}

func (proxy *SingleProxy) constructRouter(routes []*Route) {
	r := mux.NewRouter()
	for _, route := range routes {
		log.Printf("service: %v url: %v", proxy.service, route.Path)
		r.PathPrefix(route.Path).HandlerFunc(proxy.buildHandler(route.Require))
	}
	proxy.router = r
}

func (proxy *SingleProxy) buildHandler(required map[string]string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if len(required) == 0 {
			proxy.proxy.ServeHTTP(w, r)
			return
		}
		claims, err := jwt.GetClaimsFromRequest(r, proxy.jwtdCrt)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error() + "\n"))
			return
		}
		if service, ok := claims["service"].(string); ok {
			if service != proxy.service {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("service mismatch\n"))
				return
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("no valid service field in token\n"))
			return
		}
		if project, ok := claims["project"].(string); ok {
			if project != proxy.project {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("project mismatch\n"))
				return
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("no valid project field in token\n"))
			return
		}
		if err = proxy.validateNbf(claims); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error() + "\n"))
			return
		}
		if err = proxy.validateExp(claims); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error() + "\n"))
			return
		}
		if err = proxy.validateLabels(claims, proxy.resolveVariables(required, mux.Vars(r))); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error() + "\n"))
			return
		}
		proxy.proxy.ServeHTTP(w, r)
	}
}

func (proxy *SingleProxy) resolveVariables(reqs map[string]string, vars map[string]string) map[string]string {
	res := make(map[string]string)
	for key, value := range reqs {
		if len(key) > 0 && key[0] == '$' {
			varName := key[1:]
			if val, ok := vars[varName]; ok {
				key = val
			}
		}
		if len(value) > 0 && value[0] == '$' {
			varName := value[1:]
			if val, ok := vars[varName]; ok {
				value = val
			}
		}
		res[key] = value
	}
	log.Print(res)
	return res
}

func (proxy *SingleProxy) validateNbf(claims map[string]interface{}) error {
	if nbfStr, ok := claims["nbf"].(string); ok {
		nbf := time.Time{}
		err := nbf.UnmarshalText([]byte(nbfStr))
		if err != nil {
			return fmt.Errorf("failed parsing nbf string: %v", nbfStr)
		}
		if time.Now().Before(nbf) {
			return errors.New("nbf is in the future")
		}
		return nil
	}
	return errors.New("no nbf given")
}

func (proxy *SingleProxy) validateExp(claims map[string]interface{}) error {
	if expStr, ok := claims["exp"].(string); ok {
		exp := time.Time{}
		err := exp.UnmarshalText([]byte(expStr))
		if err != nil {
			return fmt.Errorf("failed parsing exp string: %v", expStr)
		}
		if !time.Now().Before(exp) {
			return errors.New("exp is in the past")
		}
		return nil
	}
	return errors.New("no exp given")
}

func (proxy *SingleProxy) validateLabels(claims map[string]interface{}, required map[string]string) error {
	if labels, ok := claims["labels"].(map[string]interface{}); ok {
		for rKey, rValue := range required {
			if uValue, ok := labels[rKey]; !ok || uValue != rValue {
				return fmt.Errorf("can not validate label %v:%v", rKey, rValue)
			}
		}
		return nil
	}
	return errors.New("no labels given")
}
