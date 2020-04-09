package main

import (
	"fmt"
	"github.com/kunnpuu/ghttp"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "hello http server")
	})
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "pong")
	})
	ghttp.ListenAndServe(":8080", nil)
}