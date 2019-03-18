// +build mage

package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/gobuffalo/packr/v2/jam"

	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
)

var Default = Build

const exeName = "gooey"

// Build makes the executable.
func Build() error {
	mg.Deps(Packr)

	//run("go", map[string]string{}, "get", "github.com/zserge/webview")
	//run("go", map[string]string{}, "get", "-u", "github.com/gobuffalo/packr/v2/packr2")

	if er := run("go",
		map[string]string{"GOOS": "linux"}, "build", "-o", exeName); er != nil {
		return er
	}

	return nil
}

func Packr() error {
	mg.Deps(BuildWASM)
	os.Chdir("./www")
	defer os.Chdir("..")

	opts := jam.PackOptions{
		IgnoreImports: false,
		Legacy:        false,
		StoreCmd:      "",
		Roots:         []string{},
	}
	return jam.Pack(opts)
}

func BuildWASM() error {
	// Copy the wasm js support file
	download("./www/wasm_exec.js", "https://raw.githubusercontent.com/golang/go/master/misc/wasm/wasm_exec.js")

	wasmExec := "./www/main.wasm"
	os.Remove(wasmExec)
	// Build WASM stuff.
	return run("go", map[string]string{"GOOS": "js", "GOARCH": "wasm"}, "build", "-o", wasmExec, "./wasm/main.go")
}

// run starts a shell cmd.
func run(name string, override map[string]string, arg ...string) error {
	log.Printf("Calling '%s %s'...", name, strings.Join(arg, " "))
	cmd := exec.Command(name, arg...)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	cmd.Stdout = os.Stdout

	cmd.Env = os.Environ()
	for k, v := range override {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}
	if err := cmd.Run(); err != nil {
		defer log.Println("...Failed.")
		log.Println(string(stderr.Bytes()))
		log.Println(err)
		return err
	}
	log.Printf("...Returned ok.")
	return nil
}

// copy the src file to dst. Any existing file will be overwritten and will not
// copy file attributes.
func copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

// download will pull the file at url into filepath
func download(filepath, url string) error {
	log.Printf("Downloading '%s'...\n", url)
	defer log.Println("...Done.")
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		os.Remove(filepath)
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		os.Remove(filepath)
		return err
	}
	return nil
}
