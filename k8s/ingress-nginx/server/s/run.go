package s

import (
	"log"
	"net"
)

func Run(addr string) {
	l, e := net.Listen(`tcp`, addr)
	if e != nil {
		log.Fatalln(e)
		return
	}
	log.Println("listen success", addr)

	// serve
	s := newServer(l)
	e = s.Serve()
	if e != nil {
		log.Fatalln(e)
		return
	}
}
