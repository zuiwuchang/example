package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"server/s"
	"time"

	grpc_api "server/protocol/api"
	_version "server/version"

	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	var (
		help    bool
		version bool
	)
	flag.BoolVar(&help, `help`, false, `dsiplay help`)
	flag.BoolVar(&version, `version`, false, `dsiplay version`)
	flag.Parse()
	if help {
		flag.PrintDefaults()
		return
	} else if version {
		fmt.Println(_version.Version)
		return
	}

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
func newGrpcClient(addr string, h2 bool) (*grpc.ClientConn, error) {
	if h2 {
		return grpc.Dial(addr,
			grpc.WithTransportCredentials(
				credentials.NewTLS(&tls.Config{
					InsecureSkipVerify: true,
				}),
			),
		)
	} else {
		return grpc.Dial(addr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
	}

}
func runClient(addr string) {
	c := os.Getenv("ExampleClient")
	val := os.Getenv("ExampleH2")
	h2 := val != "" && val != "false" && val != "0"
	var cc *grpc.ClientConn
	var client *HttpClient
	if c == "grpc" {
		var e error
		cc, e = newGrpcClient(addr, h2)
		if e != nil {
			log.Fatalln(e)
		}
	} else if c == "http" {
		client = newHttpClient(addr, h2)
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
				cc, e = newGrpcClient(addr, h2)
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
				client = newHttpClient(addr, h2)
			}
		default:
			fmt.Println("not support commnad")
		}
	}
}

type HttpClient struct {
	h2     bool
	addr   string
	client *http.Client
	dialer *websocket.Dialer
}

func newHttpClient(addr string, h2 bool) *HttpClient {
	var transport http.RoundTripper
	var dialer *websocket.Dialer
	if h2 {
		conf := &tls.Config{InsecureSkipVerify: true}
		transport = &http.Transport{
			TLSClientConfig: conf,
		}
		dialer = &websocket.Dialer{
			Proxy:            http.ProxyFromEnvironment,
			HandshakeTimeout: 45 * time.Second,
			TLSClientConfig:  conf,
		}
	} else {
		dialer = websocket.DefaultDialer
	}
	return &HttpClient{
		addr: addr,
		h2:   h2,
		client: &http.Client{
			Transport: transport,
		},
		dialer: dialer,
	}
}
func (h *HttpClient) url(path string) string {
	var scheme string
	if h.h2 {
		scheme = `https://`
	} else {
		scheme = `http://`
	}
	return scheme + h.addr + path
}
func (h *HttpClient) Get(path, value string) {
	resp, e := h.client.Get(h.url(path) + "?" + url.Values{
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
	resp, e := h.client.Post(h.url(path), `application/json; charset=utf-8`, bytes.NewReader(b))
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
	var scheme string
	if h.h2 {
		scheme = `wss://`
	} else {
		scheme = `ws://`
	}
	ws, _, e := h.dialer.Dial(scheme+h.addr+path, nil)
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
