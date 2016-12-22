package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/trusch/jwtd/db"
)

type GroupHandler struct {
	router *mux.Router
}

func NewGroupHandler() *GroupHandler {
	handler := &GroupHandler{mux.NewRouter()}
	handler.router.Path("/project/{project}/group").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.handleGetGroups(w, r)
	})
	handler.router.Path("/project/{project}/group").Methods("POST").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.handleCreateGroup(w, r)
	})
	handler.router.Path("/project/{project}/group/{group}").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.handleGetGroup(w, r)
	})
	handler.router.Path("/project/{project}/group/{group}").Methods("DELETE").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.handleDeleteGroup(w, r)
	})
	handler.router.Path("/project/{project}/group/{group}").Methods("PATCH").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.handleUpdateGroup(w, r)
	})
	return handler
}

func (h *GroupHandler) handleGetGroups(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groups, err := database.ListGroups(vars["project"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(groups)
}

func (h *GroupHandler) handleGetGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	group, err := database.GetGroup(vars["project"], vars["group"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(group)
}

func (h *GroupHandler) handleDeleteGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	err := database.DelGroup(vars["project"], vars["group"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte("delete ok"))
}

func (h *GroupHandler) handleCreateGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	project := vars["project"]
	decoder := json.NewDecoder(r.Body)
	group := &db.Group{}
	err := decoder.Decode(group)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	err = database.CreateGroup(project, group.Name, group.Rights)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte("create ok"))
}

func (h *GroupHandler) handleUpdateGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var (
		project   = vars["project"]
		groupname = vars["group"]
	)
	decoder := json.NewDecoder(r.Body)
	group := &db.Group{}
	err := decoder.Decode(group)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}
	group.Name = groupname
	group.Project = project
	err = database.UpdateGroup(group)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte("update ok"))
}

func (h *GroupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Print("Group handler is called!")
	h.router.ServeHTTP(w, r)
}
