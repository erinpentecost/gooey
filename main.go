// +build js,wasm

package main

import (
	"fmt"
	"syscall/js"
	"time"
)

func main() {
	// This is code that is compiled to WASM and run on WebKit.
	// It's controlled by the build tags on the first line.
	fmt.Println("Hello, WebAssembly!")
	app := js.Global().Get("document").Call("getElementById", "app")
	
	for {
		app.Set("innerHTML", fmt.Sprintf("WASM running!\n<hr/>%s", time.Now().String()))
		time.Sleep(time.Second)
	}
}
