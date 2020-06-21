package main

import (
	"flag"
	"github.com/mosout/diof/rpc"
	"log"
)

var (
	port string
)

func init() {
	flag.StringVar(&port, "port", "7365", "The port which server will listen on.")
}

func main() {
	flag.Parse()
	s, err := rpc.NewServer("localhost", port)
	if err != nil {
		log.Fatal(err)
	}
	err = s.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
