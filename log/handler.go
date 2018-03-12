package log

import (
	"io"
	"os"
	"sync"
)

type Handler interface {
	Handle(lr LogRecord)
	SetLogLevel(l LogLevel)
	LogLevel() LogLevel
	Formatter() LogFormatter
	SetFormatter(lf LogFormatter)
}

type BaseHandler struct {
	minLevel LogLevel
	fmt      LogFormatter
}

func (h *BaseHandler) LogLevel() LogLevel {
	return h.minLevel
}

func (h *BaseHandler) SetLogLevel(lvl LogLevel) {
	h.minLevel = lvl
}

func (h *BaseHandler) Formatter() LogFormatter {
	return h.fmt
}

func (h *BaseHandler) SetFormatter(lf LogFormatter) {
	h.fmt = lf
}

type StreamHandler struct {
	BaseHandler
	stream io.Writer
	lock   sync.Mutex
}

func (sh *StreamHandler) Handle(lr LogRecord) {
	if lr.Level() < sh.minLevel {
		return
	}
	sh.lock.Lock()
	defer sh.lock.Unlock()

	sh.fmt.Format(sh.stream, lr)
	io.WriteString(sh.stream, "\n")
}

func NewStdoutHandler() *StreamHandler {
	return &StreamHandler{
		BaseHandler: BaseHandler{
			minLevel: LevelNotSet,
			fmt:      &ColorFormatter{},
		},
		stream: os.Stdout,
	}
}

func NewStderrHandler() *StreamHandler {
	return &StreamHandler{
		BaseHandler: BaseHandler{
			minLevel: LevelNotSet,
			fmt:      &ColorFormatter{},
		},
		stream: os.Stderr,
	}
}

func NewFileHandler(filename string) (*StreamHandler, error) {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return nil, err
	}
	return &StreamHandler{
		BaseHandler: BaseHandler{
			minLevel: LevelNotSet,
			fmt:      &DefaultFormatter{},
		},
		stream: f,
	}, nil
}
