package logex

import (
	"os"
	"strings"
	"testing"
)

func b() error {
	_, err := os.Open("dflkjasldfkas")
	return Track(err)
}

func a() error {
	return Track(b())
}

func TestError(t *testing.T) {
	te := Track(a())
	errInfo := te.StackError()
	if strings.Contains(errInfo, "logex.b:11") &&
		strings.Contains(errInfo, "logex.a:15") &&
		strings.Contains(errInfo, "logex.TestError:19") {
	} else {
		t.Error("fail", te.StackError())
	}
}
