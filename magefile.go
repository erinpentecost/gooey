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
	"path"
	"strings"

	"github.com/shurcooL/vfsgen"

	"github.com/magefile/mage/mg"
)

var Default = Build

const exeName = "gooey"
const clientWASMMain = "./client/main.go"
const staticDir = "./client/www"

// Build makes the standalone executable.
func Build() error {
	mg.Deps(EmbedWWW)

	if er := run("go",
		map[string]string{"GOOS": "linux"}, "build", "-o", exeName); er != nil {
		return er
	}

	return nil
}

// EmbedWWW packages up the static client/www directory into a binary.
func EmbedWWW() error {
	mg.Deps(BuildWASM)

	fs := http.Dir(staticDir)

	return vfsgen.Generate(fs, vfsgen.Options{
		PackageName:  "main",
		BuildTags:    "!wasm",
		VariableName: "Assets",
	})
}

// Run starts up a static local server for quick testing.
func Run() error {
	mg.Deps(BuildWASM)
	port := 8089
	log.Printf("Running on http://localhost:%v/main.html", port)
	log.Println("Ctrl-C to exit.")
	return http.ListenAndServe(fmt.Sprintf(":%v", port), http.FileServer(http.Dir(staticDir)))
}

// BuildWASM builds the client/main.go file into client/www/main.wasm.
func BuildWASM() error {
	// Copy the wasm js support file
	// download("./www/wasm_exec.js", "https://raw.githubusercontent.com/golang/go/master/misc/wasm/wasm_exec.js")
	targetJS := path.Join(staticDir, "wasm_exec.js")
	sourceJS := path.Join(os.ExpandEnv("${GOROOT}"), "misc", "wasm", "wasm_exec.js")
	if _, err := os.Stat(sourceJS); !os.IsNotExist(err) {
		if er := copy(sourceJS, targetJS); er != nil {
			return er
		}
	}

	wasmExec := path.Join(staticDir, "main.wasm")
	os.Remove(wasmExec)
	// Build WASM stuff.
	return run("go", map[string]string{"GOOS": "js", "GOARCH": "wasm"}, "build", "-o", wasmExec, clientWASMMain)
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
	log.Printf("Copying '%s' to '%s'...", src, dst)
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
	if err = out.Close(); err != nil {
		return err
	}
	log.Printf("...Returned ok.")
	return nil
}
