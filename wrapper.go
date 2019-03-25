// +build !js,!wasm

package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"flag"

	"github.com/zserge/webview"
)

func main() {
	// Parse input flags
	var debug = flag.Bool("debug", false, "turn on webkit debug")
	flag.Parse()

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

	w := webview.New(webview.Settings{
		Width:  400,
		Height: 300,
		Title:  "gooey",
		URL:    targetURL,
		Resizable: true,
		Debug: *debug,
		ExternalInvokeCallback: handleRPC,
	})
	defer w.Exit()
	w.Run()
}

func handleRPC(w webview.WebView, data string) {
	// If you want to call non-WASM native Go code,
	// you'd handle it with this function (or something
	// like it).
	log.Printf("RPC: %s\n", data)
}