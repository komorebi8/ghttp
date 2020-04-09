# ghttp

A http server based on github.com/panjf2000/gnet.

It can be used like net/http.

Please run the code on Linux or MacOS. It will fail on Windows and wsl1.

An example:

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
    	// replace http.ListenAndServe() with ghttp.ListenAndServe()
    	ghttp.ListenAndServe(":8080", nil)
    }