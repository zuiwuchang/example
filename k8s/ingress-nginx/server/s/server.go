package s

import (
	"context"
	"log"
	"net"
	"net/http"
	"strings"

	"server/m/register"
	"server/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)

type Server struct {
	pipe  *PipeListener
	gpipe *grpc.Server

	tcp  net.Listener
	gtcp *grpc.Server

	mux *gin.Engine
}

func newServer(l net.Listener) (s *Server) {
	pipe := ListenPipe()
	clientConn, e := grpc.Dial(`pipe`,
		grpc.WithInsecure(),
		grpc.WithContextDialer(func(c context.Context, s string) (net.Conn, error) {
			return pipe.DialContext(c, `pipe`, s)
		}),
	)
	if e != nil {
		log.Fatalln(`pipe`, e)
	}

	gateway := utils.NewGateway()
	mux := gin.Default()
	mux.RedirectTrailingSlash = false
	register.HTTP(clientConn, mux, gateway)

	gpipe := newGRPC(gateway, clientConn)
	s = &Server{
		pipe:  pipe,
		gpipe: gpipe,
		tcp:   l,
		gtcp:  newGRPC(nil, nil),
		mux:   mux,
	}
	return
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	contextType := r.Header.Get(`Content-Type`)
	if r.ProtoMajor == 2 && strings.Contains(contextType, `application/grpc`) {
		s.gtcp.ServeHTTP(w, r)
	} else {
		s.mux.ServeHTTP(w, r)
	}
}
func (s *Server) Serve() (e error) {
	go s.gpipe.Serve(s.pipe)

	// h2c
	var httpServer http.Server
	var http2Server http2.Server
	e = http2.ConfigureServer(&httpServer, &http2Server)
	if e != nil {
		return
	}
	httpServer.Handler = h2c.NewHandler(s, &http2Server)
	// h2c Serve
	e = httpServer.Serve(s.tcp)
	return
}
func (s *Server) ServeTLS(certFile, keyFile string) (e error) {
	go s.gpipe.Serve(s.pipe)

	e = http.ServeTLS(s.tcp, s, certFile, keyFile)
	return
}
