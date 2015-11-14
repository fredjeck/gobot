// The MIT License (MIT)
// Copyright (c) 2015 Frédéric Jecker

// Autobot is a dead simple run and forget daemon that will build, lint and test your go projects
package main

import (
	"flag"
	"fmt"
	"github.com/fredjeck/fswatcher"
	"github.com/fredjeck/status"
	"github.com/shiena/ansicolor"
	"go/build"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	version = "1.1.0"
)

var (

	// Path from which the golint executable will be loaded at runtime
	golint = ""

	// ---------
	//   Flags
	// ---------

	// FileSystem poll interval
	pollInterval int
	// Terminal width
	termWidth int

	stdout = ansicolor.NewAnsiColorWriter(os.Stdout)
	stderr = ansicolor.NewAnsiColorWriter(os.Stderr)
)

func init() {
	flag.IntVar(&pollInterval, "i", 500, "Elapsed time betwen checks for modified source code.")
	flag.IntVar(&termWidth, "w", 80, "Terminal width (default 80)")
}

func banner() {
	fmt.Fprint(stdout, "   ______      ____        __ \n")
	fmt.Fprint(stdout, "  / ____/___  / __ )____  / /_\n")
	fmt.Fprintf(stdout, " / / __/ __ \\/ __  / __ \\/ __/\n")
	fmt.Fprintf(stdout, "/ /_/ / /_/ / /_/ / /_/ / /_  version %s\n", version)
	fmt.Fprintf(stdout, "\\____/\\____/_____/\\____/\\__/  \x1b[1m\x1b[34m$GOPATH = %s\x1b[0m\n", build.Default.GOPATH)
}

func main() {
	flag.Parse()

	banner()

	cwd, err := os.Getwd()
	if err != nil {
		die("An error occured while getting the working directory's path, please check your credentials")
	}

	if !strings.HasPrefix(strings.ToLower(cwd), GoPath) {
		die("Gobot must be run within your $GOPATH")
	}

	// FIXME This won't work if golint is installed elsewhere and available on the path.
	golint = filepath.Join(build.Default.GOPATH, "bin", "golint")
	if runtime.GOOS == "windows" {
		golint += ".exe"
	}
	_, lerr := os.Stat(golint)
	if lerr != nil {
		warn("Linter not found, disabling linting for the current session")
	}

	sw := status.NewWriter(termWidth)

	f := func(path string, info os.FileInfo) bool {
		directory := filepath.Dir(path)
		//start := Timestamp()

		fmt.Print("\n")
		sw.Pendingf("Building %v", ToModuleName(directory))
		stdout, stderr, err := Execute(directory, "go", "build")
		if err != nil {
			sw.MkFailure().Done()
			fmt.Fprintln(os.Stderr, stderr.String())
			return true
		}
		sw.MkSuccess().Done()

		stdout.Reset()
		stderr.Reset()
		err = nil

		if len(golint) > 0 {
			sw.Pending("Linting")
			stdout, stderr, err = Execute(directory, golint)
			if err != nil {
				sw.MkFailure().Done()
				warn("An error occured while linting, disabling linting for this session")
				golint = ""
			} else if len(stdout.String()) > 0 {
				sw.MkWarning().Done()
				fmt.Fprintln(os.Stdout, "#"+filepath.Base(directory)+"\n"+stdout.String())
			} else {
				sw.MkSuccess().Done()
			}
		}

		stdout.Reset()
		stderr.Reset()
		err = nil

		sw.Pending("Testing")
		stdout, stderr, err = Execute(directory, "go", "test")
		if err != nil {
			sw.MkFailure().Done()
			fmt.Fprintln(os.Stdout, "Some unit tests failed :\n", stdout.String())
		} else {
			sw.MkSuccess().Done()
		}

		return true
	}

	w, e := fswatcher.NewFsWatcher(cwd, pollInterval)
	if e != nil {
		die(fmt.Sprintf("Cannot monitor %s : %s", cwd, e.Error()))
	}

	w.Skip(".git")
	w.RegisterFileExtension(".go", f)
	w.Watch()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
	fmt.Println("Exiting")
}

func die(msg string) {
	fmt.Fprintln(stderr, "\x1b[1m\x1b[31m"+msg+"\x1b[0m")
	os.Exit(1)
}

func warn(msg string) {
	fmt.Fprintln(stdout, "\x1b[1m\x1b[33m"+msg+"\x1b[0m")
}
