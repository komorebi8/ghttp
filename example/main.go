package main

import (
	"fmt"
	"net/http"
	"github.com/kunnpuu/ghttp"
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