# Gooey

Gooey is a template project for creating a single-file [WebKit](https://github.com/zserge/webview) application that runs go compiled to [WASM](https://github.com/golang/go/wiki/WebAssembly).

The `main.go` file in `./client/` is the entrypoint for your application. This go source is compiled to `./client/www/main.wasm`. After compilation, `./client/www/*` is packaged with [vfsgen](https://github.com/shurcooL/vfsgen). Once that's done, the second compilation occurs to create a single binary to serve up the application. You can either distribute the resulting binary, or serve the files in `./client/www` over a normal web server.

## Setup

Install Mage:

```sh
cd ~
git clone https://github.com/magefile/mage
cd mage
go run bootstrap.go
```

Build with `mage -v build`. If you get *Package gtk+-3.0 was not found in the pkg-config search path* on Linux, then run `sudo apt-get install libwebkit2gtk-4.0-dev`.
