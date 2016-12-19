package server

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/trusch/jwtd/db"
	"github.com/trusch/jwtd/jwt"
)

var (
	database *db.DB
	key      interface{}
)

type TokenRequest struct {
	Project  string            `json:"project"`
	Username string            `json:"username"`
	Password string            `json:"password"`
	Service  string            `json:"service"`
	Labels   map[string]string `json:"labels"`
}

func Init(path, keyfile string) error {
	d, err := db.New(path)
	if err != nil {
		d = &db.DB{ConfigPath: path, Config: &db.ConfigFile{}}
		e := d.CreateUser("default", "admin", "admin", []string{"admin"})
		if e != nil {
			return e
		}
		e = d.CreateGroup("default", "admin", map[string]map[string]string{
			"jwtd": map[string]string{
				"scope": "admin",
			},
		})
		if e != nil {
			return err
		}
	}
	go func() {
		stat, _ := os.Stat(path)
		modtime := stat.ModTime()
		for {
			stat, _ = os.Stat(path)
			newModtime := stat.ModTime()
			if modtime.Unix() != newModtime.Unix() {
				dNew, e := db.New(path)
				if e == nil {
					database = dNew
				}
			}
			time.Sleep(5 * time.Second)
		}
	}()
	database = d
	k, err := jwt.LoadPrivateKey(keyfile)
	if err != nil {
		return err
	}
	key = k
	http.HandleFunc("/", handleRequest)
	return nil
}

func Serve(uri, certfile, keyfile string) {
	log.Fatal(http.ListenAndServeTLS(uri, certfile, keyfile, nil))
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
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

	user, err := database.GetUser(request.Project, request.Username)
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
	if ok, e := user.CheckRights(database, request.Service, request.Labels); e != nil || !ok {
		log.Printf("failed request: no rights (user: %v service: %v, subject: %v)", request.Username, request.Service, request.Labels)
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
