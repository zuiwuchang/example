package main

import (
	"echo/server"
	"log"
	"net"

	"google.golang.org/grpc"
)

func runServer(addr string) {
	l, e := net.Listen(`tcp`, addr)
	if e != nil {
		log.Fatalln(e)
	}
	defer l.Close()
	log.Println(`h2c listen on`, addr)
	s := grpc.NewServer()
	server.RegisterServerServer(s, Server{})
	if e := s.Serve(l); e != nil {
		log.Fatalln(e)
	}
}

type Server struct {
	server.UnimplementedServerServer
}

func (Server) Echo(stream server.Server_EchoServer) error {
	for {
		req, e := stream.Recv()
		if e != nil {
			log.Println(`recv err:`, e)
			return e
		}
		log.Println(`recv:`, req.Message)
		e = stream.Send(&server.EchoResponse{
			Message: req.Message,
		})
		if e != nil {
			log.Println(`send err:`, e)
			return e
		}
		log.Println(`send:`, req.Message)
	}
}
