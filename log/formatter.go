package log

import (
	"io"
	"time"
)

type LogFormatter interface {
	Format(writer io.Writer, record LogRecord)
}

type DefaultFormatter struct {
}

func (df *DefaultFormatter) Format(writer io.Writer, lr LogRecord) {
	io.WriteString(writer, "[")
	io.WriteString(writer, lr.Time().Format(time.RFC3339))
	io.WriteString(writer, "] ")
	if lr.Identifier() != "" {
		io.WriteString(writer, lr.Identifier())
		io.WriteString(writer, " ")
	}
	io.WriteString(writer, lr.Level().String())
	io.WriteString(writer, ": <")
	io.WriteString(writer, lr.Caller())
	io.WriteString(writer, ">: ")
	io.WriteString(writer, lr.Format())
}

type ColorFormatter struct {
}

func (cf *ColorFormatter) Format(writer io.Writer, lr LogRecord) {
	io.WriteString(writer, "[")
	switch lr.Level() {
	case LevelDebug:
		io.WriteString(writer, "\033[96m")
	case LevelInfo:
		io.WriteString(writer, "\033[94m ")
	case LevelWarn:
		io.WriteString(writer, "\033[33m ")
	case LevelError:
		io.WriteString(writer, "\033[91m")
	case LevelFatal:
		io.WriteString(writer, "\033[7;91m")
	default:
		io.WriteString(writer, "\033[37m")
	}
	io.WriteString(writer, lr.Level().String())
	io.WriteString(writer, "\033[m] ")
	io.WriteString(writer, lr.Format())
}
