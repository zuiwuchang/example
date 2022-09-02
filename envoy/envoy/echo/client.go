package main

import (
	"context"
	"echo/server"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func runClient(addr string) {
	cc, e := grpc.Dial(addr,
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)
	if e != nil {
		log.Fatalln(e)
	}
	defer cc.Close()
	NewClient(cc).Serve()
}

type Client struct {
	client server.ServerClient
	ch     chan string
}

func NewClient(cc *grpc.ClientConn) *Client {
	return &Client{
		client: server.NewServerClient(cc),
		ch:     make(chan string),
	}
}

func (c *Client) Serve() {
	stream, e := c.client.Echo(context.Background())
	if e != nil {
		log.Fatalln(e)
	}
	go c.recv(stream)

	var msg string
	for {
		fmt.Print(`input send message: `)
		fmt.Scan(&msg)
		if msg == `quit` {
			stream.CloseSend()
			break
		}
		log.Println(`send:`, msg)
		e = stream.Send(&server.EchoRequest{Message: msg})
		if e != nil {
			log.Fatalln(e)
		}
		msg = <-c.ch
		log.Println(`recv:`, msg)
	}
}
func (c *Client) recv(stream server.Server_EchoClient) {
	for {
		resp, e := stream.Recv()
		if e != nil {
			if e == context.Canceled {
				break
			}
			log.Fatalln(e)
		}
		c.ch <- resp.Message
	}
}
