package main

import (
	"flag"
	"log"

	"github.com/trusch/jwtd/server"
)

var certFile = flag.String("cert", "/etc/jwtd/jwtd.crt", "certfile")
var keyFile = flag.String("key", "/etc/jwtd/jwtd.key", "keyfile")
var configFile = flag.String("config", "/etc/jwtd/config.yaml", "config path")
var listen = flag.String("listen", ":443", "listen address")

func main() {
	log.SetFlags(log.Lshortfile)
	flag.Parse()
	err := server.Init(*configFile, *keyFile)
	if err != nil {
		log.Fatal(err)
	}
	server.Serve(*listen, *certFile, *keyFile)
}
