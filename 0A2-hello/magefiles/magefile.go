package main

import (
	"fmt"
	"runtime"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/magefile/mage/target"
)

var appExecutable = "bin/demo"
var appDir = "."

const goCompiler = "go"

var appGlobs = []string{
	"magefiles/magefile.go",
	"*.go",
	"*/*.go",
}

//goland:noinspection GoBoolExpressions
func init() {
	if runtime.GOOS == "windows" {
		appExecutable += ".exe"
	}
}

func buildApp() error {
	changes, err := target.Glob(appExecutable, appGlobs...)
	if err != nil {
		return err
	}
	if !changes {
		return nil
	}
	fmt.Println("> Building App...")
	return sh.RunV(goCompiler, "build", "-v", "-o", appExecutable, appDir)
}

func BuildWasm() error {
	changes, err := target.Glob("web/app.wasm", appGlobs...)
	if err != nil {
		return err
	}
	if !changes {
		return nil
	}
	fmt.Println("> Building WASM...")
	return sh.RunWithV(map[string]string{"GOOS": "js", "GOARCH": "wasm"}, goCompiler, "build", "-v", "-o",
		"web/app.wasm", appDir)
}

func Build() error {
	mg.Deps(BuildWasm)
	return buildApp()
}

func Run() error {
	mg.Deps(Build)
	return sh.RunV(appExecutable)
}
