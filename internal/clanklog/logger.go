package clanklog

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

var logger *loggerDefault

func Info(args ...interface{}) {
	logger.Info(args...)
}
func Infoln(args ...interface{}) {
	logger.Infoln(args...)
}
func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}
func Warning(args ...interface{}) {
	logger.Warning(args...)
}
func Warningln(args ...interface{}) {
	logger.Warningln(args...)
}
func Warningf(format string, args ...interface{}) {
	logger.Warningf(format, args...)
}
func Error(args ...interface{}) {
	logger.Error(args...)
}
func Errorln(args ...interface{}) {
	logger.Errorln(args...)
}
func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}
func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}
func Fatalln(args ...interface{}) {
	logger.Fatalln(args...)
}
func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}
func V(l int) bool {
	return logger.V(l)
}

type logType int

const (
	infoLog logType = iota
	warningLog
	errorLog
	fatalLog
)

// severityName contains the string representation of each severity.
var severityName = []string{
	infoLog:    "INFO",
	warningLog: "WARNING",
	errorLog:   "ERROR",
	fatalLog:   "FATAL",
}

type loggerDefault struct {
	m          []*log.Logger
	v          int
	jsonFormat bool
	callDepth  int
}

type loggerV2Config struct {
	verbose    int
	jsonFormat bool
}

func NewLogger() {
	logger = newLogger()
}

func newLogger() *loggerDefault {
	errorW := ioutil.Discard
	warningW := ioutil.Discard
	infoW := ioutil.Discard

	logLevel := os.Getenv("CLANK_LOG_LEVEL")
	switch logLevel {
	case "", "ERROR", "error": // If env is unset, set level to ERROR.
		errorW = os.Stderr
	case "WARNING", "warning":
		warningW = os.Stderr
	case "INFO", "info":
		infoW = os.Stderr
	}

	var v int
	vLevel := os.Getenv("CLANK_LOG_VERBOSITY_LEVEL")
	if vl, err := strconv.Atoi(vLevel); err == nil {
		v = vl
	}
	jsonFormat := strings.EqualFold(os.Getenv("CLANK_LOG_FORMATTER"), "json")
	return newLoggerWithConfig(infoW, warningW, errorW, loggerV2Config{
		verbose:    v,
		jsonFormat: jsonFormat,
	})
}

func newLoggerWithConfig(infoW, warningW, errorW io.Writer, c loggerV2Config) *loggerDefault {
	var m []*log.Logger
	flag := log.LstdFlags
	if c.jsonFormat {
		flag = 0
	}
	m = append(m, log.New(infoW, "", flag))                           // info
	m = append(m, log.New(io.MultiWriter(infoW, warningW), "", flag)) // warning
	ew := io.MultiWriter(infoW, warningW, errorW)                     // ew will be used for error and fatal.
	m = append(m, log.New(ew, "", flag))
	m = append(m, log.New(ew, "", flag))
	return &loggerDefault{m: m, v: c.verbose, jsonFormat: c.jsonFormat, callDepth: 2}
}

func (l *loggerDefault) output(t logType, s string) {
	sevStr := severityName[t]
	if !l.jsonFormat {
		l.m[t].Output(l.callDepth, fmt.Sprintf("%v: %v", sevStr, s))
		return
	}
	// TODO: we can also include the logging component, but that needs more
	// (API) changes.
	b, _ := json.Marshal(map[string]string{
		"Log":     sevStr,
		"message": s,
	})
	l.m[t].Output(l.callDepth, string(b))
}

func (l *loggerDefault) Info(args ...interface{}) {
	l.output(infoLog, fmt.Sprint(args...))
}

func (l *loggerDefault) Infoln(args ...interface{}) {
	l.output(infoLog, fmt.Sprintln(args...))
}

func (l *loggerDefault) Infof(format string, args ...interface{}) {
	l.output(infoLog, fmt.Sprintf(format, args...))
}

func (l *loggerDefault) Warning(args ...interface{}) {
	l.output(warningLog, fmt.Sprint(args...))
}

func (l *loggerDefault) Warningln(args ...interface{}) {
	l.output(warningLog, fmt.Sprintln(args...))
}

func (l *loggerDefault) Warningf(format string, args ...interface{}) {
	l.output(warningLog, fmt.Sprintf(format, args...))
}

func (l *loggerDefault) Error(args ...interface{}) {
	l.output(errorLog, fmt.Sprint(args...))
}

func (l *loggerDefault) Errorln(args ...interface{}) {
	l.output(errorLog, fmt.Sprintln(args...))
}

func (l *loggerDefault) Errorf(format string, args ...interface{}) {
	l.output(errorLog, fmt.Sprintf(format, args...))
}

func (l *loggerDefault) Fatal(args ...interface{}) {
	l.output(fatalLog, fmt.Sprint(args...))
	os.Exit(1)
}

func (l *loggerDefault) Fatalln(args ...interface{}) {
	l.output(fatalLog, fmt.Sprintln(args...))
	os.Exit(1)
}

func (l *loggerDefault) Fatalf(format string, args ...interface{}) {
	l.output(fatalLog, fmt.Sprintf(format, args...))
	os.Exit(1)
}

func (l *loggerDefault) V(v int) bool {
	return v <= l.v
}
