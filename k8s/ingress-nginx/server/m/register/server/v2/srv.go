package v1

import (
	"context"
	"io"
	"log"
	grpc_api "server/protocol/api"
)

type server struct {
	grpc_api.UnimplementedV2Server
}

func (s server) Get(ctx context.Context, req *grpc_api.GetRequest) (resp *grpc_api.GetResponse, e error) {
	resp = &grpc_api.GetResponse{
		Value: req.Value + " from " + hostname + " v2",
	}
	log.Println(req, resp)
	return
}
func (s server) Post(ctx context.Context, req *grpc_api.PostRequest) (resp *grpc_api.PostResponse, e error) {
	resp = &grpc_api.PostResponse{
		Value: req.Value + " from " + hostname + " v2",
	}
	log.Println(req, resp)
	return
}
func (s server) Stream(stream grpc_api.V2_StreamServer) (e error) {
	for {
		req, e := stream.Recv()
		if e != nil {
			if e != io.EOF {
				log.Println(e)
			}
			break
		}
		log.Println("recv", req)
		resp := &grpc_api.StreamResponse{
			Value: req.Value + " from " + hostname + " v2",
		}
		e = stream.Send(resp)
		if e != nil {
			log.Println(e)
			break
		}
		log.Println("send", resp)
	}
	return
}
