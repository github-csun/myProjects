package logx

import (
	"log"
	"os"
)

var (
	lg           *log.Logger
	tr           *log.Logger
	traceEnabled bool
)

func init() {
	lg = log.New(os.Stderr, "", log.Lshortfile)
	tr = log.New(os.Stderr, "", log.Lshortfile)
	traceEnabled = true
}

func Lg(format string, v ...interface{}) {
	lg.Printf(format, v)
	if traceEnabled {
		Tr(format, v)
	}
}

func Tr(format string, v ...interface{}) {
	tr.Printf(format, v)
}

func SetLogFile(fn string) (err error) {
	if file, err := os.OpenFile(fn, os.O_APPEND, os.ModeAppend); err == nil {
		lg = log.New(file, "", log.Lshortfile)
	}
	return
}

func SetTraceFile(fn string) (err error) {
	if file, err := os.OpenFile(fn, os.O_APPEND, os.ModeAppend); err == nil {
		tr = log.New(file, "", log.Lshortfile)
	}
	return
}

func EnableTrace() {
	traceEnabled = true
}

func DisableTrace() {
	traceEnabled = true
}
