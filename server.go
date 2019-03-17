// +build !js,!wasm
package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gobuffalo/packr/v2"
	"github.com/zserge/webview"
)

func main() {
	box := packr.New("www", "./www")

	log.Printf("Box contents:")
	for _, f := range box.List() {
		log.Printf("    %s", f)
	}

	addr := "127.0.0.1:8622"

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	go func() {
		log.Printf("Hosting on %s ...\n", addr)
		http.Handle("/", http.FileServer(box))
		log.Fatal(http.Serve(ln, nil))
	}()

	log.Print("Opening window...\n", addr)
	webview.Open("Hello", fmt.Sprintf("http://%s/index.html", addr), 400, 300, false)
}
