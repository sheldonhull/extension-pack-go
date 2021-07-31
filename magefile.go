//+build mage

package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// allow user to override go executable by running as GOEXE=xxx make ... on unix-like systems
var goexe = "go"

func init() {
	if exe := os.Getenv("GOEXE"); exe != "" {
		goexe = exe
	}

	// We want to use Go 1.11 modules even if the source lives inside GOPATH.
	// The default is "auto".
	os.Setenv("GO111MODULE", "on")
}

// Install all dependencies to local directory
func Build() error {
	if err := sh.Run("go", "mod", "download"); err != nil {
		return err
	}
	return sh.Run("go", "install", "./...")
}

func runCmd(env map[string]string, cmd string, args ...interface{}) error {
	if mg.Verbose() {
		return runWith(env, cmd, args...)
	}
	output, err := sh.OutputWith(env, cmd, argsToStrings(args...)...)
	if err != nil {
		fmt.Fprint(os.Stderr, output)
	}

	return err
}

func runWith(env map[string]string, cmd string, inArgs ...interface{}) error {
	s := argsToStrings(inArgs...)
	return sh.RunWith(env, cmd, s...)
}

func isGoLatest() bool {
	return strings.Contains(runtime.Version(), "1.16")
}

func isCI() bool {
	return os.Getenv("CI") != ""
}

func argsToStrings(v ...interface{}) []string {
	var args []string
	for _, arg := range v {
		switch v := arg.(type) {
		case string:
			if v != "" {
				args = append(args, v)
			}
		case []string:
			if v != nil {
				args = append(args, v...)
			}
		default:
			panic("invalid type")
		}
	}

	return args
}

// Remove the temporarily generated files from Release.
func Clean() error {
	return sh.Rm("dist")
}

// InstallTools installs all the tooling for the project
func InstallTools() error {
	cmds := []string{
		"go install github.com/evilmartians/lefthook@latest",
		"go install github.com/owenrumney/squealer/cmd/squealer@latest",
		"GO111MODULE=on go get github.com/zricethezav/gitleaks/v7",
		"go install github.com/hekike/unchain@master",
	}

	for _, cmd := range cmds {
		if err := sh.Run(cmd); err != nil {
			return err
		}
	}

	return nil
}
