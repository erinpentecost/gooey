# Gooey

Gooey is a Vulkan demo with HTML-defined UI.

* [Lorca](https://github.com/zserge/lorca)
* [Webview](https://github.com/zserge/webview) run a local browser that's hooked in to your executable. communicate over java.
  * [machinebox](https://github.com/machinebox/desktop)
* [go compiles to wasm](https://github.com/golang/go/wiki/WebAssembly) compile go code to browser native code.
* [goki/gi])(https://github.com/goki/gi) meh
* [mozilla webrender](https://github.com/servo/webrender/wiki) Webrender is an experimental renderer for Servo that aims to draw web content like a modern game engine.

Nothing here is performant. I need to be able to draw on a screen from vulkan directly, too. Vulkan needs to be the window owner, but I should be able to pass in the output of a headless browser (built on webkit) as a layer to pass into vulkan. However, webkit uses GTK internally as well as a browser engine. Seems super overkill and slow.

### Setup

Install **Go**, then install Mage:

```sh
cd ~
git clone https://github.com/magefile/mage
cd mage
go run bootstrap.go
```

Try to build with `mage -v build`.

If you get *Package gtk+-3.0 was not found in the pkg-config search path*, then run `sudo apt-get install libwebkit2gtk-4.0-dev`

Once built, you can run inside your browser with `goexec 'http.ListenAndServe(":8080", http.FileServer(http.Dir("./www")))'`