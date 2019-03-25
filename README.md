# Gooey

Gooey is a template project for creating a single-file [WebKit](https://github.com/zserge/webview) application that runs go compiled to [WASM](https://github.com/golang/go/wiki/WebAssembly).

Use build tags `// +build js,wasm` to identify the code you want compiled and run on the WebKit. I'm guessing most of your code is going to use those tags. The magefile builds the `js,wasm` target first, and sticks the resulting *main.wasm* binary in `./www`. Afterward, `./www/*` is packaged with [vfsgen](https://github.com/shurcooL/vfsgen) to create *wwwData.go*. Finally, native targets are executed and you end up with something that looks like a native application but is built on HTML/WASM. You can distribute the resulting binary, or serve the files in `./www` directly over a normal web server.

## Setup

Install Mage:

```sh
cd ~
git clone https://github.com/magefile/mage
cd mage
go run bootstrap.go
```

**Build** with `mage -v build`. If you get *Package gtk+-3.0 was not found in the pkg-config search path* on Linux, then run `sudo apt-get install libwebkit2gtk-4.0-dev`.