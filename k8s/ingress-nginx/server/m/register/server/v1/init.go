package v1

import (
	"context"
	"os"

	grpc_api "server/protocol/api"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Module int

func (Module) RegisterGRPC(srv *grpc.Server) {
	grpc_api.RegisterV1Server(srv, server{})
}

var hostname string

func (Module) RegisterGateway(gateway *runtime.ServeMux, cc *grpc.ClientConn) error {
	hostname = os.Getenv("HOSTNAME")
	return grpc_api.RegisterV1Handler(context.Background(), gateway, cc)
}
