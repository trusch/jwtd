package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	st0rage "github.com/trusch/jwtd/storage"
)

type UserHandler struct {
	router *mux.Router
}

type UserData struct {
	Username string   `json:"username"`
	Password string   `json:"password"`
	Groups   []string `json:"groups"`
}

func NewUserHandler() *UserHandler {
	handler := &UserHandler{mux.NewRouter()}
	handler.router.Path("/project/{project}/user").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.handleGetUsers(w, r)
	})
	handler.router.Path("/project/{project}/user").Methods("POST").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.handleCreateUser(w, r)
	})
	handler.router.Path("/project/{project}/user/{user}").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.handleGetUser(w, r)
	})
	handler.router.Path("/project/{project}/user/{user}").Methods("DELETE").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.handleDeleteUser(w, r)
	})
	handler.router.Path("/project/{project}/user/{user}").Methods("PATCH").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.handleUpdateUser(w, r)
	})
	return handler
}

func (h *UserHandler) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	users, err := storage.ListUsers(vars["project"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}
	res := make([]*st0rage.User, len(users))
	for i, user := range users {
		res[i] = &st0rage.User{Name: user.Name, Groups: user.Groups}
	}
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(res)
}

func (h *UserHandler) handleGetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user, err := storage.GetUser(vars["project"], vars["user"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}
	res := &st0rage.User{Name: user.Name, Groups: user.Groups}
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(res)
}

func (h *UserHandler) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	err := storage.DelUser(vars["project"], vars["user"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte("delete ok"))
}

func (h *UserHandler) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userData := &UserData{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(userData)
	if userData.Username == "" || userData.Password == "" {
		log.Print("invalid data in create user request")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("username and password needed"))
		return
	}
	err = storage.CreateUser(vars["project"], userData.Username, userData.Password, userData.Groups)
	if err != nil {
		log.Print("error in create user request storage call: ", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte("create ok"))
}

func (h *UserHandler) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var (
		project  = vars["project"]
		username = vars["user"]
		password string
		groups   []string
	)
	userData := &UserData{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(userData)
	password = userData.Password
	groups = userData.Groups
	user, err := storage.GetUser(project, username)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	if password != "" {
		groups = user.Groups
		err = storage.DelUser(project, username)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		err = storage.CreateUser(project, username, password, groups)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
	} else {
		user.Groups = groups
		err = storage.UpdateUser(project, user)
		log.Print(user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
	}
	w.Write([]byte("update ok"))
}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Print("User handler is called!")
	h.router.ServeHTTP(w, r)
}
