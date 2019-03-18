// +build !js,!wasm
package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gobuffalo/packr/v2"
	"github.com/zserge/webview"
)

func main() {
	box := packr.New("www", "./www")

	log.Printf("Box contents:")
	for _, f := range box.List() {
		log.Printf("    %s", f)
	}

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
		http.Handle("/", http.FileServer(box))
		http.HandleFunc("/time", handleTime)
		log.Fatal(http.Serve(ln, nil))
	}()

	// Open up a window
	log.Print("Opening window...\n", addr)
	webview.Open("Hello", fmt.Sprintf("http://%s/main.html", addr), 400, 300, false)
}

// handleTime handles the /time endpoint.
// It is used to demonstrate an API call from app.js.
func handleTime(w http.ResponseWriter, r *http.Request) {
	if _, err := io.WriteString(w, time.Now().Format(time.RFC822)); err != nil {
		log.Println(err)
	}
}
