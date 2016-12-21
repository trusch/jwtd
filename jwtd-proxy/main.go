package main

import (
	"flag"
	"log"
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
	log.Fatal(ListenAndServeTLSSNI(cfg.Listen, cfg.GetCerts(), proxy))
}
