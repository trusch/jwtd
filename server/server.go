package server

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
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
	if err := initDB(path); err != nil {
		return err
	}
	initDBReloader(path)
	if err := initKey(keyfile); err != nil {
		return err
	}
	router := mux.NewRouter()
	router.Path("/token").Handler(NewTokenHandler())
	router.PathPrefix("/project/{project}/user").Handler(NewUserHandler())
	router.PathPrefix("/project/{project}/group").Handler(NewGroupHandler())
	http.Handle("/", router)
	return nil
}

func initDB(path string) error {
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
	database = d
	return nil
}

func initDBReloader(path string) {
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
}

func initKey(keyfile string) error {
	k, err := jwt.LoadPrivateKey(keyfile)
	if err != nil {
		return err
	}
	key = k
	return nil
}

func Serve(uri string) {
	log.Fatal(http.ListenAndServe(uri, nil))
}

func handleRequest(w http.ResponseWriter, r *http.Request) {

}
