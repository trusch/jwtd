package main

import (
	"log"

	"github.com/trusch/jwtd/jwtd-ctl/cmd"
)

func main() {
	log.SetFlags(log.Lshortfile)
	cmd.Execute()
}
