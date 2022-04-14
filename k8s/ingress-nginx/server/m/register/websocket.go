package register

import (
	"net/http"
	"sync/atomic"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

type Error struct {
	Code    codes.Code `json:"code,omitempty"`
	Message string     `json:"message,omitempty"`
}
type Websocket struct {
	*websocket.Conn
}

func NewWebsocket(c *gin.Context, responseHeader http.Header) (conn Websocket, e error) {
	if !c.IsWebsocket() {
		e = status.Error(codes.InvalidArgument, `expect websocket`)
		rError(c, e)
		return
	}
	ws, e := upgrader.Upgrade(c.Writer, c.Request, responseHeader)
	if e != nil {
		e = status.Error(codes.Unknown, e.Error())
		rError(c, e)
		return
	}
	conn = Websocket{ws}
	return
}
func rError(c *gin.Context, e error) {
	if e == nil {
		ResponseError(c, Error{
			Code:    codes.OK,
			Message: codes.OK.String(),
		})
		return
	}
	ResponseError(c, Error{
		Code:    status.Code(e),
		Message: e.Error(),
	})
}
func ResponseError(c *gin.Context, err Error) {
	code := http.StatusOK
	switch err.Code {
	case codes.OK:
		// code = http.StatusOK
	case codes.Canceled:
		code = http.StatusRequestTimeout
	case codes.Unknown:
		code = http.StatusInternalServerError
	case codes.InvalidArgument:
		code = http.StatusBadRequest
	case codes.DeadlineExceeded:
		code = http.StatusGatewayTimeout
	case codes.NotFound:
		code = http.StatusNotFound
	case codes.AlreadyExists:
		code = http.StatusConflict
	case codes.PermissionDenied:
		code = http.StatusForbidden
	case codes.ResourceExhausted:
		code = http.StatusTooManyRequests
	case codes.FailedPrecondition:
		code = http.StatusBadRequest
	case codes.Aborted:
		code = http.StatusConflict
	case codes.OutOfRange:
		code = http.StatusBadRequest
	case codes.Unimplemented:
		code = http.StatusNotImplemented
	case codes.Internal:
		code = http.StatusInternalServerError
	case codes.Unavailable:
		code = http.StatusServiceUnavailable
	case codes.DataLoss:
		code = http.StatusInternalServerError
	case codes.Unauthenticated:
		code = http.StatusUnauthorized
	default:
		code = http.StatusInternalServerError
	}
	c.JSON(code, err)
}
func (w Websocket) SendMessage(m proto.Message) error {
	b, e := Marshal(m)
	if e != nil {
		return e
	}
	return w.WriteMessage(websocket.TextMessage, b)
}
func (w Websocket) SendBinary(b []byte) error {
	return w.WriteMessage(websocket.BinaryMessage, b)
}
func (w Websocket) Send(v interface{}) error {
	return w.WriteJSON(v)
}
func (w Websocket) Success() error {
	return w.Send(Error{
		Code:    codes.OK,
		Message: codes.OK.String(),
	})
}
func (w Websocket) Error(e error) error {
	if e == nil {
		return w.Send(Error{
			Code:    codes.OK,
			Message: codes.OK.String(),
		})
	} else {
		return w.Send(Error{
			Code:    status.Code(e),
			Message: e.Error(),
		})
	}
}
func (w Websocket) Forward(f Forward) {
	work := newWebsocketForward(w, f)
	work.Serve()
}

type websocketForward struct {
	w      Websocket
	f      Forward
	closed int32
	cancel chan struct{}
}

func newWebsocketForward(w Websocket, f Forward) *websocketForward {
	return &websocketForward{
		w:      w,
		f:      f,
		cancel: make(chan struct{}),
	}
}
func (wf *websocketForward) Serve() {
	go wf.request()
	go wf.response()
	<-wf.cancel
	wf.w.Close()
	wf.f.CloseSend()
}
func (wf *websocketForward) request() {
	var counted uint64
	for {
		t, p, e := wf.w.ReadMessage()
		if e != nil {
			break
		}
		e = wf.f.Request(counted, t, p)
		if e != nil {
			wf.w.Error(e)
			break
		}
		counted++
	}

	if wf.closed == 0 &&
		atomic.SwapInt32(&wf.closed, 1) == 0 {
		close(wf.cancel)
	}
}
func (wf *websocketForward) response() {
	var counted uint64
	for {
		e := wf.f.Response(counted)
		if e != nil {
			wf.w.Error(e)
			break
		}
		counted++
	}
	if wf.closed == 0 &&
		atomic.SwapInt32(&wf.closed, 1) == 0 {
		close(wf.cancel)
	}
}
