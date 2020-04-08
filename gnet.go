package ghttp

import (
	"bufio"
	"bytes"
	"github.com/panjf2000/gnet"
	"github.com/panjf2000/gnet/pool/goroutine"
	"log"
)

type GServer struct {
	*gnet.EventServer
	pool *goroutine.Pool
}

func NewGServer() GServer {
	return GServer{
		pool:goroutine.Default(),
	}
}

func (s *GServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	req, _ := ReadRequest(bufio.NewReader(bytes.NewBuffer(frame)))
	resp := &response{
		req:req,
		w:bufio.NewWriter(bytes.NewBuffer(out)),
		handlerHeader:Header{},
	}
	resp.Write([]byte("hello world"))
	log.Println(req)
	return
}