// +build !js,!wasm
package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/zserge/webview"
)

func main() {
	// Reserve port to act as server
	ln, err := net.Listen("tcp", ":0")
	addr := ln.Addr()
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	// Serve up http requests on the port
	go func() {
		log.Printf("Hosting on %s ...\n", addr)
		http.Handle("/", http.FileServer(Assets))
		log.Fatal(http.Serve(ln, nil))
	}()

	targetURL := fmt.Sprintf("http://%s/main.html", addr)

	// Open up a window
	log.Printf("Opening window to %s ...\n", targetURL)
	webview.Open("Hello", targetURL, 400, 300, false)
}
