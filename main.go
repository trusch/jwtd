package main

import (
	"flag"

	"github.com/trusch/jwtd/server"
)

var certFile = flag.String("cert", "cert.pem", "certfile")
var keyFile = flag.String("key", "key.pem", "keyfile")
var dbAddr = flag.String("db", "localhost", "db uri")
var listen = flag.String("listen", ":443", "listen address")

func main() {
	flag.Parse()
	server.Init(*dbAddr, *keyFile)
	server.Serve(*listen, *certFile, *keyFile)
}
