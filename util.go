package main

import (
	"bytes"
	"go/build"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Lowercased GOPATH used for trimming.
var GoPath = strings.ToLower(build.Default.GOPATH)

// Path to the SRC directory.
var SrcRoot = filepath.Join(GoPath, "src")

// Timestamp creates a unix timestamp for the current date/time
func Timestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// Execute runs a command from the given working directory passing the provided arguments and returns
// the StdOut, StdErr and an Error as third parameter if the command execution returned an error
// The returned error is nil if the command runs, has no problems copying stdin, stdout, and stderr, and exits with a zero exit status.
// If the command fails to run or doesn't complete successfully, the error is of type *ExitError.
// Other error types may be returned for I/O problems.
func Execute(workingDir string, command string, args ...string) (bytes.Buffer, bytes.Buffer, error) {
	cmd := exec.Command(command, args...)
	cmd.Dir = workingDir

	var err, out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &err
	return out, err, cmd.Run()
}

// ToModuleName converts a path under the current GOROOT to a module name.
// If path isn't located under GOROOT, returns path, untouched.
func ToModuleName(path string) string {
	p := strings.ToLower(path)

	if strings.HasPrefix(p, GoPath) {
		return strings.TrimPrefix(filepath.ToSlash(strings.TrimPrefix(p, SrcRoot)), "/")
	}

	return path
}
