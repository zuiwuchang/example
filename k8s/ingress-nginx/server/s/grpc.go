package s

import (
	"context"

	"server/m/register"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func auth(ctx context.Context) (context.Context, error) {
	return ctx, nil
}
func newGRPC(gateway *runtime.ServeMux, cc *grpc.ClientConn) (srv *grpc.Server) {
	opts := []grpc.ServerOption{}

	opts = append(opts,
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_recovery.StreamServerInterceptor(),
			grpc_auth.StreamServerInterceptor(auth),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(),
			grpc_auth.UnaryServerInterceptor(auth),
		)),
	)

	srv = grpc.NewServer(opts...)
	register.GRPC(srv, gateway, cc)
	return
}
