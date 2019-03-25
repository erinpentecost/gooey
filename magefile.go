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

// exeName is the name of the executable file.
const exeName = "gooey"

const clientWASMMain = "./"
const staticDir = "./www"

// Build makes the standalone executable.
func Build() error {
	mg.Deps(embedWWW)

	binPath := path.Join("dist", os.ExpandEnv("${GOOS}"), exeName)

	return run("go", map[string]string{}, "build", "-o", binPath)
}

// embedWWW packages up the static client/www directory into a binary.
func embedWWW() error {
	mg.Deps(BuildWASM)

	fs := http.Dir(staticDir)

	return vfsgen.Generate(fs, vfsgen.Options{
		Filename:     "./wwwData.go",
		PackageName:  "main",
		BuildTags:    "!js,!wasm",
		VariableName: "Assets",
	})
}

// Run starts up a static local server for quick testing.
// This skips binary building and relies on an external browser.
func Run() error {
	mg.Deps(BuildWASM)
	port := 8089
	log.Printf("Running on http://localhost:%v/main.html", port)
	log.Println("Ctrl-C to exit.")
	return http.ListenAndServe(fmt.Sprintf(":%v", port), http.FileServer(http.Dir(staticDir)))
}

// BuildWASM builds the js/wasm packages.
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

	// Build WASM stuff.
	wasmExec := path.Join(staticDir, "main.wasm")
	return run("go", map[string]string{"GOOS": "js", "GOARCH": "wasm"}, "build", "-o", wasmExec)
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
