package log

import (
	"context"
	"encoding/json"
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
	Vars() map[string]interface{}
	json.Marshaler
}

type logRecord struct {
	time    time.Time
	level   LogLevel
	message string
	args    []interface{}
	caller  string
	ctx     context.Context
}

func (lr logRecord) MarshalJSON() ([]byte, error) {
	record := struct {
		Time       time.Time
		Level      LogLevel
		LevelStr   string
		Message    string
		Arguments  []interface{}
		Caller     string
		Identifier string
		Vars       map[string]interface{}
	}{
		lr.time,
		lr.level,
		lr.level.String(),
		lr.message,
		lr.args,
		lr.caller,
		lr.Identifier(),
		lr.Vars(),
	}

	return json.Marshal(record)
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

func (lr *logRecord) Vars() map[string]interface{} {
	val := lr.ctx.Value(logVarsCtxKey)
	if val == nil {
		return make(map[string]interface{})
	}
	return val.(map[string]interface{})
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
