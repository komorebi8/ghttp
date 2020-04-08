package ghttp

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"github.com/panjf2000/gnet"
	"github.com/panjf2000/gnet/pool/goroutine"
	"log"
	"strconv"
	"testing"
	"time"
)

type echoServer struct {
	*gnet.EventServer
	pool *goroutine.Pool
}

func (es *echoServer) OnInitComplete(srv gnet.Server) (action gnet.Action) {
	log.Printf("Echo server is listening on %s (multi-cores: %t, loops: %d)\n",
		srv.Addr.String(), srv.Multicore, srv.NumEventLoop)
	return
}
func (es *echoServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	log.Println("new frame\n", string(frame))
	req, err := ReadRequest(bufio.NewReader(bytes.NewBuffer(frame)))
	if err != nil {
		log.Println(err.Error())
	} else {
		log.Println("new request", req)
	}

	out = appendResp(out, "200 OK", "", "hello world")
	return
}

type GResponse struct {
	handlerHeader Header

}

func (res *GResponse) Header() Header {
	if res.handlerHeader == nil {
		res.handlerHeader = make(Header)
	}
	return res.handlerHeader
}

func (res *GResponse) Write(b[]byte) (n int, err error) {
	return 0, nil
}

func (res *GResponse) WriteHeader(statusCode int)  {

}


func appendResp(b []byte, status, head, body string) []byte {
	b = append(b, "HTTP/1.1"...)
	b = append(b, ' ')
	b = append(b, status...)
	b = append(b, '\r', '\n')
	b = append(b, "Server: gnet\r\n"...)
	b = append(b, "Date: "...)
	b = time.Now().AppendFormat(b, "Mon, 02 Jan 2006 15:04:05 GMT")
	b = append(b, '\r', '\n')
	if len(body) > 0 {
		b = append(b, "Content-Length: "...)
		b = strconv.AppendInt(b, int64(len(body)), 10)
		b = append(b, '\r', '\n')
	}
	b = append(b, head...)
	b = append(b, '\r', '\n')
	if len(body) > 0 {
		b = append(b, body...)
	}
	return b
}

func TestNewGServer(t *testing.T) {
	var port int
	var multicore bool

	// Example command: go run echo.go --port 9000 --multicore=true
	flag.IntVar(&port, "port", 9000, "--port 9000")
	flag.BoolVar(&multicore, "multicore", false, "--multicore true")
	flag.Parse()
	echo := new(echoServer)
	echo.pool = goroutine.Default()
	log.Fatal(gnet.Serve(echo, fmt.Sprintf("tcp://:%d", port), gnet.WithMulticore(multicore)))
}
