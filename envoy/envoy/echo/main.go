package main

import (
	"flag"
	"log"
)

func main() {
	var (
		help, client bool
		addr         string
	)
	flag.BoolVar(&help, `help`, false, `disapley help`)
	flag.BoolVar(&client, `client`, false, `run as client`)
	flag.StringVar(&addr, `addr`, `127.0.0.1:9000`, `server listen address; or client connect address`)
	flag.Parse()
	if help {
		flag.PrintDefaults()
		return
	}
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if client {
		runClient(addr)
	} else {
		runServer(addr)
	}
}
