package logex

import (
	"bytes"
	"errors"
	"fmt"
	"path"
	"runtime"
	"strconv"
	"strings"
)

func NewError(info string, format ...interface{}) *TraceError {
	err := fmt.Errorf(info, format...)
	return &TraceError{
		error: err,
	}
}

func Is(e1, e2 error) bool {
	if e, ok := e1.(*TraceError); ok {
		e1 = e.error
	}
	if e, ok := e2.(*TraceError); ok {
		e2 = e.error
	}
	return e1 == e2
}

type TraceError struct {
	error
	format []interface{}
	stack  []string
}

func (t *TraceError) Error() string {
	if t == nil {
		return "<NILTraceErr>"
	}
	if t.format == nil {
		if t.error == nil {
			panic(t.stack)
		}
		return t.error.Error()
	}
	return fmt.Sprintf(t.error.Error(), t.format...)
}

func (t *TraceError) Follow(err error) *TraceError {
	if t == nil {
		return nil
	}
	if te, ok := err.(*TraceError); ok {
		if len(te.stack) > 0 {
			te.stack[len(te.stack)-1] += ":" + err.Error()
		}
		t.stack = append(te.stack, t.stack...)
	}
	return t
}

func (t *TraceError) Format(obj ...interface{}) *TraceError {
	if t == nil {
		return nil
	}
	t.format = obj
	return t
}

func (t *TraceError) StackError() string {
	if t == nil {
		return t.Error()
	}
	if len(t.stack) == 0 {
		return t.Error()
	}
	return fmt.Sprintf("[%s] %s", strings.Join(t.stack, ";"), t.Error())
}

func Tracef(err error, obj ...interface{}) *TraceError {
	e := TraceEx(1, err).Format(obj...)
	return e
}

// set runtime info to error
func Trace(err error, info ...interface{}) *TraceError {
	return TraceEx(1, err, info...)
}

func TraceIfError(err error, info ...interface{}) error {
	if err != nil {
		return Trace(err, info...)
	}
	return nil
}

func joinInterface(info []interface{}, ch string) string {
	ret := bytes.NewBuffer(make([]byte, 0, 512))
	for idx, o := range info {
		if idx > 0 {
			ret.WriteString(ch)
		}
		ret.WriteString(fmt.Sprint(o))
	}
	return ret.String()
}

func TraceEx(depth int, err error, info ...interface{}) *TraceError {
	if err == nil {
		return nil
	}
	pc, _, line, _ := runtime.Caller(1 + depth)
	name := runtime.FuncForPC(pc).Name()
	name = path.Base(name)
	stack := name + ":" + strconv.Itoa(line)
	if len(info) > 0 {
		stack += "(" + joinInterface(info, ",") + ")"
	}
	if te, ok := err.(*TraceError); ok {
		te.stack = append(te.stack, stack)
		return te
	}
	return &TraceError{err, nil, []string{stack}}
}

func NewTraceError(info ...interface{}) *TraceError {
	return TraceEx(1, errors.New(fmt.Sprint(info...)))
}
