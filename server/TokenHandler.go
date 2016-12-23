package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/trusch/jwtd/jwt"
)

type TokenHandler struct{}

func NewTokenHandler() *TokenHandler {
	return &TokenHandler{}
}

func (h *TokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	request := &TokenRequest{}
	err := decoder.Decode(request)
	if err != nil {
		log.Print("failed request...")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if request.Project == "" {
		request.Project = "default"
	}

	user, err := storage.GetUser(request.Project, request.Username)
	if err != nil {
		log.Printf("failed request: no such user (%v)", request.Username)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if ok, e := user.CheckPassword(request.Password); e != nil || !ok {
		log.Printf("failed request: wrong password (user: %v)", request.Username)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if ok, e := user.CheckRights(storage, request.Project, request.Service, request.Labels); e != nil || !ok {
		log.Printf("failed request: no rights (user: %v service: %v, labels: %v)", request.Username, request.Service, request.Labels)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	claims := jwt.Claims{
		"user":    request.Username,
		"service": request.Service,
		"project": request.Project,
		"labels":  request.Labels,
		"nbf":     time.Now(),
		"exp":     time.Now().Add(10 * time.Minute),
	}
	token, err := jwt.CreateToken(claims, key)
	if err != nil {
		log.Print("failed request: can not generate token (wtf?!)")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("successfully created token (user: %v service: %v, labels: %v)", request.Username, request.Service, request.Labels)
	w.Write([]byte(token))
}
