package log

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Logger interface {
	LogCtxf(identifier string, level LogLevel, message string, args ...interface{}) LogRecord
	Logf(level LogLevel, message string, args ...interface{}) LogRecord

	Debugf(message string, args ...interface{}) LogRecord
	Infof(message string, args ...interface{}) LogRecord
	Warnf(message string, args ...interface{}) LogRecord
	Errorf(message string, args ...interface{}) LogRecord
	Fatalf(message string, args ...interface{})

	DebugCtxf(identifier string, message string, args ...interface{}) LogRecord
	InfoCtxf(identifier string, message string, args ...interface{}) LogRecord
	WarnCtxf(identifier string, message string, args ...interface{}) LogRecord
	ErrorCtxf(identifier string, message string, args ...interface{}) LogRecord
	FatalCtxf(identifier string, message string, args ...interface{})

	Println(args ...interface{}) LogRecord
	Printf(message string, args ...interface{}) LogRecord
	Fatal(args ...interface{})

	Handlers() []Handler
	SetHandlers(hs []Handler)
	SetLogLevel(l LogLevel)
	LogLevel() LogLevel
	Identifier() string
	SetIdentifier(i string)

	NewContextLogger() Logger
}

type logger struct {
	handlers          []Handler
	minLevel          LogLevel
	defaultIdentifier string
}

func (l *logger) Identifier() string {
	return l.defaultIdentifier
}

func (l *logger) SetIdentifier(i string) {
	l.defaultIdentifier = i
}

func (l *logger) Handlers() []Handler {
	return l.handlers
}

func (l *logger) SetHandlers(hs []Handler) {
	l.handlers = hs
}

func (l *logger) Logf(level LogLevel, message string, args ...interface{}) LogRecord {
	return l.LogCtxf(l.defaultIdentifier, level, message, args...)
}

func (l *logger) LogCtxf(identifier string, level LogLevel, message string, args ...interface{}) LogRecord {
	if level < l.minLevel {
		return nil
	}
	lr := &logRecord{
		time:      time.Now(),
		level:     level,
		message:   message,
		args:      args,
		caller:    getCaller(),
		identifer: identifier,
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

func (l *logger) DebugCtxf(identifier string, message string, args ...interface{}) LogRecord {
	return l.LogCtxf(identifier, LevelDebug, message, args...)
}

func (l *logger) InfoCtxf(identifier string, message string, args ...interface{}) LogRecord {
	return l.LogCtxf(identifier, LevelInfo, message, args...)
}

func (l *logger) WarnCtxf(identifier string, message string, args ...interface{}) LogRecord {
	return l.LogCtxf(identifier, LevelWarn, message, args...)
}

func (l *logger) ErrorCtxf(identifier string, message string, args ...interface{}) LogRecord {
	return l.LogCtxf(identifier, LevelError, message, args...)
}

func (l *logger) FatalCtxf(identifier string, message string, args ...interface{}) {
	record := l.LogCtxf(identifier, LevelFatal, message, args...)
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

func (l *logger) NewContextLogger() Logger {
	log := &logger{
		minLevel:          l.minLevel,
		handlers:          l.handlers,
		defaultIdentifier: NewIdentifier(),
	}
	return log
}

func NewIdentifier() string {
	u := uuid.Must(uuid.NewUUID())
	return u.String()
}
