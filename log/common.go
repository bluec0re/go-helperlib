package log

import "context"

var (
	defaultLogger = &logger{
		handlers: []Handler{
			NewStderrHandler(),
		},
		defaultCtx: NewContext(),
	}
)

const (
	LevelDebug  = LogLevel(iota)
	LevelInfo   = LogLevel(iota)
	LevelWarn   = LogLevel(iota)
	LevelError  = LogLevel(iota)
	LevelFatal  = LogLevel(iota)
	LevelNotSet = -1
)

type LogLevel int

func (ll LogLevel) String() string {
	switch ll {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

func Handlers() []Handler {
	return defaultLogger.Handlers()
}

func SetHandlers(h []Handler) {
	defaultLogger.SetHandlers(h)
}

func Logf(level LogLevel, message string, args ...interface{}) LogRecord {
	return defaultLogger.Logf(level, message, args...)
}

func Debugf(message string, args ...interface{}) LogRecord {
	return defaultLogger.Debugf(message, args...)
}

func Infof(message string, args ...interface{}) LogRecord {
	return defaultLogger.Infof(message, args...)
}

func Warnf(message string, args ...interface{}) LogRecord {
	return defaultLogger.Warnf(message, args...)
}

func Errorf(message string, args ...interface{}) LogRecord {
	return defaultLogger.Errorf(message, args...)
}

func Fatalf(message string, args ...interface{}) {
	defaultLogger.Fatalf(message, args...)
}

func Println(args ...interface{}) LogRecord {
	return defaultLogger.Println(args...)
}

func Printf(message string, args ...interface{}) LogRecord {
	return defaultLogger.Printf(message, args...)
}

func Fatal(args ...interface{}) {
	defaultLogger.Fatal(args...)
}

func AddFileHandler(filename string) (Handler, error) {
	h, err := NewFileHandler(filename)
	if err != nil {
		return nil, err
	}
	handlers := Handlers()
	handlers = append(handlers, h)
	SetHandlers(handlers)
	return h, nil
}

func NewContextLogger() Logger {
	return defaultLogger.NewContextLogger()
}

func NewContextLoggerWithParent(ctx context.Context) Logger {
	return defaultLogger.NewContextLoggerWithParent(ctx)
}

func CloseLogger() error {
	return defaultLogger.Close()
}
