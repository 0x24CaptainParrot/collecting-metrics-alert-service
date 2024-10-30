package main

import (
	"flag"
	"fmt"
	"log"
)

var runServerAddrFlag string

func parseServerFlags() {
	flag.StringVar(&runServerAddrFlag, "a", "localhost:8080", "server listens on this port")
	flag.Parse()

	if len(flag.Args()) > 0 {
		fmt.Printf("Unknown flags %v\n", flag.Args())
		log.Fatal("Error: unknown flags were given")
	}
}
