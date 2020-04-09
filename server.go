// Copyright 2019 Andy Pan. All rights reserved.
// Copyright 2017 Joshua J Baker. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package ghttp

import (
	"bufio"
	"bytes"
	"log"
	"net/http"
	"time"

	"github.com/panjf2000/gnet"
	"github.com/panjf2000/gnet/pool/goroutine"
)

type GServer struct {
	*gnet.EventServer
	pool *goroutine.Pool
	handler http.Handler
}

func (s *GServer) OnInitComplete(srv gnet.Server) (action gnet.Action) {
	log.Printf("HTTP server is listening on %s (multi-cores: %t, loops: %d)\n",
		srv.Addr.String(), srv.Multicore, srv.NumEventLoop)
	return
}

func (s *GServer) OnOpened(c gnet.Conn) (out []byte, action gnet.Action) {
	return
}

func (s *GServer) PreWrite() {

}

func (s *GServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	if c.Context() != nil {
		// bad thing happened
		out = InternalErrorServerResponseBytes
		action = gnet.Close
		return
	}
	// handle the request
	req, err := http.ReadRequest(bufio.NewReader(bytes.NewBuffer(frame)))
	if err != nil {
		out = InternalErrorServerResponseBytes
		action = gnet.Close
		return
	}
	resp := NewGResponse()
	handler := s.handler
	if handler == nil {
		handler = http.DefaultServeMux
	}
	if req.RequestURI == "*" && req.Method == "OPTIONS" {
		handler = globalOptionsHandler{}
	}
	handler.ServeHTTP(resp, req)
	out = resp.Bytes()
	action = gnet.None
	//s.pool.Submit(func() {
	//	resp := NewGResponse()
	//	var handler http.Handler
	//	handler = http.DefaultServeMux
	//	if req.RequestURI == "*" && req.Method == "OPTIONS" {
	//		handler = globalOptionsHandler{}
	//	}
	//	handler.ServeHTTP(resp, req)
	//	err :=c.AsyncWrite(resp.Bytes())
	//	if err != nil {
	//		log.Println(err.Error())
	//	}
	//})
	return
}

func (s *GServer) OnClosed(c gnet.Conn, err error) gnet.Action {
	return gnet.None
}

func (s *GServer) Tick() (delay time.Duration, action gnet.Action) {
	return
}

func ListenAndServe(addr string, handler http.Handler) error {
	server := new(GServer)
	server.handler = handler
	server.pool = goroutine.Default()
	defer server.pool.Release()
	return gnet.Serve(server, "tcp://" + addr, gnet.WithMulticore(true))
}


