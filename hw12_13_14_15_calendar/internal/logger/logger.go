package logger

import (
	"errors"
	"io"
	"log"
)

var ErrUnknownLoggerLevel = errors.New("unknown logger level")

const (
	DEBUG = iota
	INFO
	ERROR
)

type Logger struct {
	severityLevel int
	writer        *log.Logger
}

func New(level string, writer io.Writer) (*Logger, error) {
	var severity int

	switch level {
	case "INFO":
		severity = INFO
	case "ERROR":
		severity = ERROR
	case "DEBUG":
		severity = DEBUG
	default:
		return nil, ErrUnknownLoggerLevel
	}

	logger := log.Logger{}
	logger.SetFlags(log.Ldate | log.Ltime)
	logger.SetOutput(writer)

	return &Logger{
		severityLevel: severity,
		writer:        &logger,
	}, nil
}

func (l Logger) Debug(msg string, params ...any) {
	if l.severityLevel > DEBUG {
		return
	}
	l.writer.Printf("DEBUG "+msg+" ", params...)
}

func (l Logger) Info(msg string) {
	if l.severityLevel > INFO {
		return
	}
	l.writer.Printf("INFO " + msg)
}

func (l Logger) Error(msg string) {
	if l.severityLevel > ERROR {
		return
	}
	l.writer.Printf("ERROR " + msg)
}
