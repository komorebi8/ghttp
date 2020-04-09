package ghttp

import (
	"fmt"
	"net/http"
	"testing"
)

func TestListenAndServe(t *testing.T) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "hello http server")
	})
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "pong")
	})
	ListenAndServe(":8080", nil)
}
