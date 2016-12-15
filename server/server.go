package server

import (
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

func Init(path, keyfile string) error {
	d, err := db.New(path)
	if err != nil {
		d = &db.DB{ConfigPath: path, Config: &db.ConfigFile{}}
		e := d.CreateUser("admin", "admin", []string{"admin"})
		if e != nil {
			return e
		}
		e = d.CreateGroup("admin", []*db.AccessRight{&db.AccessRight{Service: "jwtd", Subject: "admin"}})
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
				dNew, err := db.New(path)
				if err == nil {
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
	var (
		username string
		password string
		service  string
		subject  string
	)
	err := r.ParseForm()
	if err != nil {
		log.Print("failed request...")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	username = r.Form.Get("username")
	password = r.Form.Get("password")
	service = r.Form.Get("service")
	subject = r.Form.Get("subject")

	user, err := database.GetUser(username)
	if err != nil {
		log.Printf("failed request: no such user (%v)", username)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if ok, e := user.CheckPassword(password); e != nil || !ok {
		log.Printf("failed request: wrong password (user: %v)", username)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if ok, e := user.CheckRights(database, service, subject); e != nil || !ok {
		log.Printf("failed request: no rights (user: %v service: %v, subject: %v)", username, service, subject)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	claims := jwt.Claims{
		"user":    username,
		"service": service,
		"subject": subject,
		"nbf":     time.Now(),
		"exp":     time.Now().Add(10 * time.Minute),
	}
	token, err := jwt.CreateToken(claims, key)
	if err != nil {
		log.Print("failed request: can not generate token (wtf?!)")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("successfully created token (user: %v service: %v, subject: %v)", username, service, subject)
	w.Write([]byte(token))
}
