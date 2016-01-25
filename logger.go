package robin

import (
	"log"
	"os"
)

const (
	LogLevelFatal = iota
	LogLevelError
	LogLevelInfo
	LogLevelDebug
)

type AppLogger interface {
	Fatal(f string)
	Error(f string)
	Info(f string)
	Debug(f string)
}

type applogger struct {
	level int
	log   *log.Logger
}

func newAppLogger(level int) *applogger {
	return &applogger{
		level: level,
		log:   log.New(os.Stderr, "", 0),
	}
}

func (l *applogger) Fatal(f string) {
	l.output(LogLevelFatal, "FATAL", f)
}

func (l *applogger) Error(f string) {
	l.output(LogLevelError, "ERROR", f)
}

func (l *applogger) Info(f string) {
	l.output(LogLevelInfo, "INFO", f)
}

func (l *applogger) Debug(f string) {
	l.output(LogLevelDebug, "DEBUG", f)
}

func (l *applogger) output(level int, prefix string, f string) {
	if level <= l.level {
		l.log.Printf("# [%s] %s", prefix, f)
	}

	if level == LogLevelFatal {
		os.Exit(1)
	}
}
