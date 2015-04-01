package logex

import (
	"os"
	"strings"
	"testing"
)

func b() error {
	_, err := os.Open("dflkjasldfkas")
	return Trace(err)
}

func a() error {
	return Trace(b())
}

func TestError(t *testing.T) {
	te := Trace(a())
	errInfo := te.StackError()
	if strings.Contains(errInfo, "logex%2ev1.b:11") &&
		strings.Contains(errInfo, "logex%2ev1.a:15") &&
		strings.Contains(errInfo, "logex%2ev1.TestError:19") {
	} else {
		t.Error("fail", te.StackError())
	}
}
