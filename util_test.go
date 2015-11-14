package main

import (
	"os"
	"strings"
	"testing"
)

func TestToModulePath(t *testing.T) {
	cwd, _ := os.Getwd()

	if strings.Compare(ToModuleName(cwd), "github.com/fredjeck/gobot") != 0 {
		t.Fail()
	}
}
