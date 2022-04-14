package register

import (
	"net/http"

	grpc_api "server/protocol/api"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024 * 32,
	WriteBufferSize: 1024 * 32,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HTTP(cc *grpc.ClientConn, engine *gin.Engine, gateway *runtime.ServeMux) {
	engine.NoRoute(func(c *gin.Context) {
		c.Status(http.StatusOK)
		if c.Request.Method == `GET` || c.Request.Method == `HEAD` {
			c.Request.Header.Set(`Method`, c.Request.Method)
		}
		gateway.ServeHTTP(c.Writer, c.Request)
	})

	api := engine.Group(`api`)

	api.GET(`v1/stream`, func(c *gin.Context) {
		ws, e := NewWebsocket(c, nil)
		if e != nil {
			return
		}
		defer ws.Close()
		client := grpc_api.NewV1Client(cc)
		stream, e := client.Stream(c.Request.Context())
		if e != nil {
			ws.Error(e)
			return
		}
		f := NewForward(func(counted uint64, messageType int, p []byte) error {
			return stream.Send(&grpc_api.StreamRequest{
				Value: string(p),
			})
		}, func(counted uint64) (e error) {
			resp, e := stream.Recv()
			if e != nil {
				return
			}
			return ws.SendMessage(resp)
		}, func() error {
			return stream.CloseSend()
		})
		ws.Forward(f)
	})
	api.GET(`v2/stream`, func(c *gin.Context) {
		ws, e := NewWebsocket(c, nil)
		if e != nil {
			return
		}
		defer ws.Close()
		client := grpc_api.NewV2Client(cc)
		stream, e := client.Stream(c.Request.Context())
		if e != nil {
			ws.Error(e)
			return
		}
		f := NewForward(func(counted uint64, messageType int, p []byte) error {
			return stream.Send(&grpc_api.StreamRequest{
				Value: string(p),
			})
		}, func(counted uint64) (e error) {
			resp, e := stream.Recv()
			if e != nil {
				return
			}
			return ws.SendMessage(resp)
		}, func() error {
			return stream.CloseSend()
		})
		ws.Forward(f)
	})
}
