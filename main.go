package main

import (
	"flag"
	"log"

	"github.com/quicsocks/quicsocks/client"
	"github.com/quicsocks/quicsocks/server"
)

var hostAddr = flag.String("host", "", "host ip:port")
var clientAddr = flag.String("client", "", "client ip:port")
var serverMod = flag.Bool("server", false, "server mode")

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	flag.Parse()
	if *serverMod {
		err := server.NewServer(*hostAddr)
		log.Fatal(err)
	}
	err := client.NewClient(*clientAddr, *hostAddr)
	log.Fatal(err)
}
