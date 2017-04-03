package main

import (
	"flag"
	"log"

	"github.com/trusch/jwtd/server"
)

var keyFile = flag.String("key", "/etc/jwtd/jwtd.key", "keyfile")
var configFile = flag.String("config", "/etc/jwtd/config.yml", "config path")
var listen = flag.String("listen", ":80", "listen address")

func main() {
	log.SetFlags(log.Lshortfile)
	flag.Parse()
	err := server.Init(*configFile, *keyFile)
	if err != nil {
		log.Fatal(err)
	}
	server.Serve(*listen)
}
