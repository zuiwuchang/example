package register

import (
	"log"
	v1 "server/m/register/server/v1"
	v2 "server/m/register/server/v2"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func GRPC(srv *grpc.Server, gateway *runtime.ServeMux, cc *grpc.ClientConn) {
	ms := []Module{
		v1.Module(0),
		v2.Module(0),
	}
	for _, m := range ms {
		m.RegisterGRPC(srv)
		if gateway != nil {
			e := m.RegisterGateway(gateway, cc)
			if e != nil {
				log.Fatalln(e)
			}
		}
	}
}

type Module interface {
	RegisterGRPC(srv *grpc.Server)
	RegisterGateway(gateway *runtime.ServeMux, cc *grpc.ClientConn) error
}
