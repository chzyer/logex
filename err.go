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

func NewError(info string, format ...interface{}) *TrackError {
	err := errors.New(fmt.Sprintf(info, format...))
	return &TrackError{
		error: err,
	}
}

func Is(e1, e2 error) bool {
	if e, ok := e1.(*TrackError); ok {
		e1 = e.error
	}
	if e, ok := e2.(*TrackError); ok {
		e2 = e.error
	}
	return e1 == e2
}

type TrackError struct {
	error
	format []interface{}
	stack  []string
}

func (t *TrackError) Error() string {
	if t == nil {
		return "<NILTrackErr>"
	}
	if t.format == nil {
		if t.error == nil {
			panic(t.stack)
		}
		return t.error.Error()
	}
	return fmt.Sprintf(t.error.Error(), t.format...)
}

func (t *TrackError) Follow(err error) *TrackError {
	if t == nil {
		return nil
	}
	if te, ok := err.(*TrackError); ok {
		if len(te.stack) > 0 {
			te.stack[len(te.stack)-1] += ":" + err.Error()
		}
		t.stack = append(te.stack, t.stack...)
	}
	return t
}

func (t *TrackError) Format(obj ...interface{}) *TrackError {
	if t == nil {
		return nil
	}
	t.format = obj
	return t
}

func (t *TrackError) StackError() string {
	if t == nil {
		return t.Error()
	}
	if len(t.stack) == 0 {
		return t.Error()
	}
	return fmt.Sprintf("[%s] %s", strings.Join(t.stack, ";"), t.Error())
}

func Trackf(err error, obj ...interface{}) *TrackError {
	e := TrackEx(1, err).Format(obj...)
	return e
}

// set runtime info to error
func Track(err error, info ...interface{}) *TrackError {
	return TrackEx(1, err, info...)
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

func TrackEx(depth int, err error, info ...interface{}) *TrackError {
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
	if te, ok := err.(*TrackError); ok {
		te.stack = append(te.stack, stack)
		return te
	}
	return &TrackError{err, nil, []string{stack}}
}

func NewTrackError(info ...interface{}) *TrackError {
	return TrackEx(1, errors.New(fmt.Sprint(info...)))
}
