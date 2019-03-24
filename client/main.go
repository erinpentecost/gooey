// +build js,wasm
package main

import (
	"fmt"
	"syscall/js"
)

func main() {
	fmt.Println("Hello, WebAssembly!")
	app := js.Global().Get("document").Call("getElementById", "app")
	app.Set("innerHTML", "WASM running!")
}
