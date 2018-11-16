package log

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/google/uuid"
)

type LogRecord interface {
	Format() string
	Time() time.Time
	Level() LogLevel
	Message() string
	Args() []interface{}
	Caller() string
	UpdateCaller()
	Identifier() string
}

type logRecord struct {
	time    time.Time
	level   LogLevel
	message string
	args    []interface{}
	caller  string
	ctx     context.Context
}

func (lr *logRecord) Format() string {
	return fmt.Sprintf(lr.message, lr.args...)
}

func (lr *logRecord) Message() string {
	return lr.message
}

func (lr *logRecord) Time() time.Time {
	return lr.time
}

func (lr *logRecord) Args() []interface{} {
	return lr.args
}

func (lr *logRecord) Level() LogLevel {
	return lr.level
}

func (lr *logRecord) Caller() string {
	return lr.caller
}

func (lr *logRecord) Identifier() string {
	val := lr.ctx.Value(logCtxKey)
	if val == nil {
		return ""
	}
	return val.(uuid.UUID).String()
}

func (lr *logRecord) Context() context.Context {
	return lr.ctx
}

func (lr *logRecord) UpdateCaller() {
	lr.caller = getCaller()
}

func getCaller() string {
	fpcs := make([]uintptr, 10)
	n := runtime.Callers(3, fpcs)
	if n == 0 {
		return ""
	}
	fpcs = fpcs[:n]
	frames := runtime.CallersFrames(fpcs)
	var frame runtime.Frame
	var more bool
	for {
		frame, more = frames.Next()
		// ignore log framework calls
		if !strings.Contains(frame.File, "go-helperlib/log") {
			break
		}
		if !more {
			break
		}
	}
	f := frame.File
	l := frame.Line
	f = f[strings.LastIndex(f, "/")+1:]
	name := frame.Function
	name = name[strings.LastIndex(name, ".")+1:]
	return fmt.Sprintf("%s (%s:%d)", name, f, l)
}
