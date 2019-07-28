package log

import (
	"context"
	"io"
	"strings"
	"time"

	"github.com/google/uuid"
)

type contextKey string

func (c contextKey) String() string {
	return "LOG_" + string(c)
}

const logCtxKey = contextKey("IDENTIFIER")
const logVarsCtxKey = contextKey("VARIABLES")

type Logger interface {
	io.Closer
	LogCtxf(ctx context.Context, level LogLevel, message string, args ...interface{}) LogRecord
	Logf(level LogLevel, message string, args ...interface{}) LogRecord

	Debugf(message string, args ...interface{}) LogRecord
	Infof(message string, args ...interface{}) LogRecord
	Warnf(message string, args ...interface{}) LogRecord
	Errorf(message string, args ...interface{}) LogRecord
	Fatalf(message string, args ...interface{})

	DebugCtxf(ctx context.Context, message string, args ...interface{}) LogRecord
	InfoCtxf(ctx context.Context, message string, args ...interface{}) LogRecord
	WarnCtxf(ctx context.Context, message string, args ...interface{}) LogRecord
	ErrorCtxf(ctx context.Context, message string, args ...interface{}) LogRecord
	FatalCtxf(ctx context.Context, message string, args ...interface{})

	Println(args ...interface{}) LogRecord
	Printf(message string, args ...interface{}) LogRecord
	Fatal(args ...interface{})

	Handlers() []Handler
	SetHandlers(hs []Handler)
	SetLogLevel(l LogLevel)
	LogLevel() LogLevel
	Context() context.Context
	SetContext(ctx context.Context)
	SetVariable(key string, value interface{})

	NewContextLogger() Logger
}

type logger struct {
	handlers   []Handler
	minLevel   LogLevel
	defaultCtx context.Context
}

func (l *logger) Context() context.Context {
	return l.defaultCtx
}

func (l *logger) SetContext(ctx context.Context) {
	l.defaultCtx = ctx
}

func (l *logger) Handlers() []Handler {
	return l.handlers
}

func (l *logger) SetHandlers(hs []Handler) {
	l.handlers = hs
}

func (l *logger) Logf(level LogLevel, message string, args ...interface{}) LogRecord {
	return l.LogCtxf(l.defaultCtx, level, message, args...)
}

func (l *logger) LogCtxf(ctx context.Context, level LogLevel, message string, args ...interface{}) LogRecord {
	if level < l.minLevel {
		return nil
	}
	lr := &logRecord{
		time:    time.Now(),
		level:   level,
		message: message,
		args:    args,
		caller:  getCaller(),
		ctx:     ctx,
	}
	for _, handler := range l.handlers {
		handler.Handle(lr)
	}
	return lr
}

func (l *logger) LogLevel() LogLevel {
	return l.minLevel
}

func (l *logger) SetLogLevel(lvl LogLevel) {
	l.minLevel = lvl
}

func (l *logger) Debugf(message string, args ...interface{}) LogRecord {
	return l.Logf(LevelDebug, message, args...)
}

func (l *logger) Infof(message string, args ...interface{}) LogRecord {
	return l.Logf(LevelInfo, message, args...)
}

func (l *logger) Warnf(message string, args ...interface{}) LogRecord {
	return l.Logf(LevelWarn, message, args...)
}

func (l *logger) Errorf(message string, args ...interface{}) LogRecord {
	return l.Logf(LevelError, message, args...)
}

func (l *logger) Fatalf(message string, args ...interface{}) {
	record := l.Logf(LevelFatal, message, args...)
	panic(record.Format())
}

func (l *logger) DebugCtxf(ctx context.Context, message string, args ...interface{}) LogRecord {
	return l.LogCtxf(ctx, LevelDebug, message, args...)
}

func (l *logger) InfoCtxf(ctx context.Context, message string, args ...interface{}) LogRecord {
	return l.LogCtxf(ctx, LevelInfo, message, args...)
}

func (l *logger) WarnCtxf(ctx context.Context, message string, args ...interface{}) LogRecord {
	return l.LogCtxf(ctx, LevelWarn, message, args...)
}

func (l *logger) ErrorCtxf(ctx context.Context, message string, args ...interface{}) LogRecord {
	return l.LogCtxf(ctx, LevelError, message, args...)
}

func (l *logger) FatalCtxf(ctx context.Context, message string, args ...interface{}) {
	record := l.LogCtxf(ctx, LevelFatal, message, args...)
	panic(record.Format())
}

func (l *logger) Printf(message string, args ...interface{}) LogRecord {
	return l.Logf(l.LogLevel(), message, args...)
}

func (l *logger) Println(args ...interface{}) LogRecord {
	message := strings.Repeat("%v ", len(args))
	message = message[0 : len(message)-1]
	return l.Printf(message, args...)
}
func (l *logger) Fatal(args ...interface{}) {
	message := strings.Repeat("%v ", len(args))
	message = message[0 : len(message)-1]
	l.Fatalf(message, args...)
}

func (l *logger) SetVariable(key string, value interface{}) {
	val := l.defaultCtx.Value(logVarsCtxKey)
	if val == nil {
		val = make(map[string]interface{})
		l.defaultCtx = context.WithValue(l.defaultCtx, logVarsCtxKey, val)
	}
	val.(map[string]interface{})[key] = value
}

func (l *logger) Close() error {
	for _, handler := range l.handlers {
		if err := handler.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (l *logger) NewContextLoggerWithParent(ctx context.Context) Logger {
	log := &logger{
		minLevel:   l.minLevel,
		handlers:   l.handlers,
		defaultCtx: NewContextWithParent(ctx),
	}
	return log
}

func (l *logger) NewContextLogger() Logger {
	return l.NewContextLoggerWithParent(l.defaultCtx)
}

func NewContext() context.Context {
	return NewContextWithParent(context.Background())
}

func NewContextWithParent(ctx context.Context) context.Context {
	u, _ := uuid.NewUUID()
	return context.WithValue(ctx, logCtxKey, u)
}
