package ghttp

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

var InternalErrorServerResponseBytes []byte

func init() {
	s := NewGResponse()
	s.WriteHeader(500)
	InternalErrorServerResponseBytes = s.Bytes()
}

type GResponse struct {
	handlerHeader http.Header
	body *bytes.Buffer
	contentLength int64
	status int
}

func NewGResponse() *GResponse {
	return &GResponse{
		handlerHeader: make(http.Header),
		body:          bytes.NewBuffer(make([]byte, 0)),
		contentLength: 0,
		status:        200,
	}
}

func header2Bytes(header *http.Header) []byte {
	b := make([]byte, 0)
	buf := bytes.NewBuffer(b)
	header.Write(buf)
	return buf.Bytes()
}

func (resp *GResponse) Bytes() []byte {
	timeBytes := make([]byte, 0)
	timeBytes = time.Now().AppendFormat(timeBytes, "Mon, 02 Jan 2006 15:04:05 GMT")
	resp.handlerHeader.Set("Date", string(timeBytes))
	if ct := resp.handlerHeader.Get("Content-Type"); ct == "" {
		resp.handlerHeader.Set("Content-Type", "application/plain; charset=utf-8")
	}
	bodyBytes := resp.body.Bytes()
	if cl := resp.handlerHeader.Get("Content-Length"); cl == "" {
		contentLengthStr := fmt.Sprintf("%d", len(bodyBytes))
		resp.handlerHeader.Set("Content-Length", contentLengthStr)
	}
	// response line
	b := make([]byte, 0)
	b = append(b, "HTTP/1.1 "...)
	statusString := fmt.Sprintf("%d %s", resp.status, http.StatusText(resp.status))
	b = append(b, statusString...)
	b = append(b, '\r', '\n')
	// header
	b = append(b, header2Bytes(&resp.handlerHeader)...)
	b = append(b, '\r', '\n')
	// body
	b = append(b, bodyBytes...)
	return b
}

func (res *GResponse) Header() http.Header {
	if res.handlerHeader == nil {
		res.handlerHeader = make(http.Header)
	}
	return res.handlerHeader
}

func (res *GResponse) Write(b[]byte) (n int, err error) {
	return res.body.Write(b)
}

func (res *GResponse) WriteHeader(statusCode int)  {
	var cl string
	if v := res.handlerHeader["Content-Length"]; len(v) > 0 {
		cl = v[0]
	} else {
		cl = ""
	}
	if cl != "" {
		v, err := strconv.ParseInt(cl, 10, 64)
		if err == nil && v >= 0 {
			res.contentLength = v
		} else {
			res.handlerHeader.Del("Content-Length")
		}
	}
	res.status = statusCode
}

type globalOptionsHandler struct{}

func (globalOptionsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Length", "0")
	if r.ContentLength != 0 {
		// Read up to 4KB of OPTIONS body (as mentioned in the
		// spec as being reserved for future use), but anything
		// over that is considered a waste of server resources
		// (or an attack) and we abort and close the connection,
		// courtesy of MaxBytesReader's EOF behavior.
		mb := http.MaxBytesReader(w, r.Body, 4<<10)
		io.Copy(ioutil.Discard, mb)
	}
}