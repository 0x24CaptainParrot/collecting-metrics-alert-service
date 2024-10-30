package main

import "flag"

var runServerAddrFlag string

func parseServerFlags() {
	flag.StringVar(&runServerAddrFlag, "a", ":8080", "server listens on this port")
	flag.Parse()
}
