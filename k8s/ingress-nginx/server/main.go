package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"server/s"
	"time"

	grpc_api "server/protocol/api"

	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	addr := os.Getenv("ExampleAddr")
	if addr == "" {
		addr = "127.0.0.1:9000"
	}
	mode := os.Getenv("ExampleMode")
	if mode == "server" || mode == "" {
		s.Run(addr)
	} else if mode == "client" {
		runClient(addr)
	} else {
		log.Fatalln("env 'ExampleMode' only supported 'server' or 'client'")
	}
}
func printHelp(http bool) {
	if http {
		fmt.Println(`http client`)
	} else {
		fmt.Println(`grpc client`)
	}
	fmt.Println(` * 1 /api/v1/get
 * 2 /api/v1/post
 * 3 /api/v1/stream
 * 4 /api/v2/get
 * 5 /api/v2/post
 * 6 /api/v2/stream
 * grpc set grpc client
 * http set http client
 * h print help`)
}
func runClient(addr string) {
	c := os.Getenv("ExampleClient")
	var cc *grpc.ClientConn
	var client *HttpClient
	if c == "grpc" {
		var e error
		cc, e = grpc.Dial(addr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if e != nil {
			log.Fatalln(e)
		}
	} else if c == "http" {
		client = &HttpClient{
			addr:   addr,
			client: &http.Client{},
		}
	} else {
		log.Fatalln("env 'ExampleClient' only supported 'grpc' or 'http'")
	}
	printHelp(cc == nil)

	var cmd string
	for {
		fmt.Scan(&cmd)
		switch cmd {
		case "1":
			v1get(client, cc)
		case "2":
			v1post(client, cc)
		case "3":
			v1stream(client, cc)
		case "4":
			v2get(client, cc)
		case "5":
			v2post(client, cc)
		case "6":
			v2stream(client, cc)
		case "h":
			printHelp(cc == nil)
		case "grpc":
			if cc == nil {
				var e error
				cc, e = grpc.Dial(addr,
					grpc.WithTransportCredentials(insecure.NewCredentials()),
				)
				if e == nil {
					client = nil
				} else {
					log.Println(e)
				}
			}
		case "http":
			if client == nil {
				cc.Close()
				cc = nil
				client = &HttpClient{
					addr:   addr,
					client: &http.Client{},
				}
			}
		default:
			fmt.Println("not support commnad")
		}
	}
}

type HttpClient struct {
	addr   string
	client *http.Client
}

func (h *HttpClient) Get(path, value string) {
	resp, e := h.client.Get(`http://` + h.addr + path + "?" + url.Values{
		`value`: {value},
	}.Encode())
	if e != nil {
		log.Println(e)
		return
	}
	defer resp.Body.Close()
	b, e := io.ReadAll(io.LimitReader(resp.Body, 1024))
	if e != nil {
		log.Println(e)
		return
	}
	fmt.Println("success", string(b))
}
func (h *HttpClient) Post(path, value string) {
	req := grpc_api.PostRequest{
		Value: value,
	}
	b, e := json.Marshal(&req)
	if e != nil {
		log.Println(e)
		return
	}
	resp, e := h.client.Post(`http://`+h.addr+path, `application/json; charset=utf-8`, bytes.NewReader(b))
	if e != nil {
		log.Println(e)
		return
	}
	defer resp.Body.Close()
	b, e = io.ReadAll(io.LimitReader(resp.Body, 1024))
	if e != nil {
		log.Println(e)
		return
	}
	fmt.Println("success", string(b))
}
func (h *HttpClient) Stream(path string) {
	ws, _, e := websocket.DefaultDialer.Dial(`ws://`+h.addr+path, nil)
	if e != nil {
		log.Println(e)
		return
	}
	defer ws.Close()
	for i := 0; i < 5; i++ {
		value := time.Now().String()
		log.Println("send", i, value)
		e = ws.WriteMessage(websocket.BinaryMessage, []byte(value))
		if e != nil {
			log.Println(e)
			return
		}
		_, b, e := ws.ReadMessage()
		if e != nil {
			log.Println(e)
			return
		}
		log.Println("recv", i, string(b))
		time.Sleep(time.Second)
	}
	fmt.Println("success")
}
func v1get(client *HttpClient, cc *grpc.ClientConn) {
	value := time.Now().String()
	url := "/api/v1/get"
	log.Println("request", url, value)
	if cc == nil {
		client.Get(url, value)
	} else {
		client := grpc_api.NewV1Client(cc)
		req := &grpc_api.GetRequest{
			Value: value,
		}
		resp, e := client.Get(context.Background(), req)
		if e != nil {
			log.Println(e)
			return
		}
		fmt.Println("success", resp)
	}
}
func v2get(client *HttpClient, cc *grpc.ClientConn) {
	value := time.Now().String()
	url := "/api/v2/get"
	log.Println("request", url, value)
	if cc == nil {
		client.Get(url, value)
	} else {
		client := grpc_api.NewV2Client(cc)
		req := &grpc_api.GetRequest{
			Value: value,
		}
		resp, e := client.Get(context.Background(), req)
		if e != nil {
			log.Println(e)
			return
		}
		fmt.Println("success", resp)
	}
}
func v1post(client *HttpClient, cc *grpc.ClientConn) {
	value := time.Now().String()
	url := "/api/v1/post"
	log.Println("request", url, value)
	if cc == nil {
		client.Post(url, value)
	} else {
		client := grpc_api.NewV1Client(cc)
		req := &grpc_api.PostRequest{
			Value: value,
		}
		resp, e := client.Post(context.Background(), req)
		if e != nil {
			log.Println(e)
			return
		}
		fmt.Println("success", resp)
	}
}
func v2post(client *HttpClient, cc *grpc.ClientConn) {
	value := time.Now().String()
	url := "/api/v2/post"
	log.Println("request", url, value)
	if cc == nil {
		client.Post(url, value)
	} else {
		client := grpc_api.NewV2Client(cc)
		req := &grpc_api.PostRequest{
			Value: value,
		}
		resp, e := client.Post(context.Background(), req)
		if e != nil {
			log.Println(e)
			return
		}
		fmt.Println("success", resp)
	}
}
func v1stream(client *HttpClient, cc *grpc.ClientConn) {
	url := "/api/v1/stream"
	log.Println("request", url)
	if cc == nil {
		client.Stream(url)
	} else {
		client := grpc_api.NewV1Client(cc)
		stream, e := client.Stream(context.Background())
		if e != nil {
			log.Println(e)
			return
		}
		grpcStream(stream)
		stream.CloseSend()
	}
}
func v2stream(client *HttpClient, cc *grpc.ClientConn) {
	url := "/api/v2/stream"
	log.Println("request", url)
	if cc == nil {
		client.Stream(url)
	} else {
		client := grpc_api.NewV2Client(cc)
		stream, e := client.Stream(context.Background())
		if e != nil {
			log.Println(e)
			return
		}
		grpcStream(stream)
		stream.CloseSend()
	}
}
func grpcStream(stream grpc_api.V2_StreamClient) {
	for i := 0; i < 5; i++ {
		value := time.Now().String()
		log.Println("send", i, value)
		e := stream.Send(&grpc_api.StreamRequest{
			Value: value,
		})
		if e != nil {
			log.Println(e)
			return
		}
		resp, e := stream.Recv()
		if e != nil {
			log.Println(e)
			return
		}
		log.Println("recv", i, resp)
		time.Sleep(time.Second)
	}
	fmt.Println("success")
}
