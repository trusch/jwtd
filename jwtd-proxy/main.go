package main

import (
	"flag"
	"log"
	"net/http"
)

var config = flag.String("config", "/etc/jwtd-proxy/config.yaml", "config file")

func main() {
	log.SetFlags(log.Lshortfile)
	flag.Parse()
	cfg, err := NewConfigFromFile(*config)
	if err != nil {
		log.Fatal(err)
	}
	proxy, err := NewProxy(cfg)
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/", proxy)
	http.ListenAndServe(cfg.Listen, nil)
	log.Print(cfg)
}
