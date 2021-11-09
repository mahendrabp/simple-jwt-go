package utils

import log "github.com/sirupsen/logrus"

type Logger interface {
	Init(debug bool, level string)

	Info(args ...interface{})
	InfoFormat(format string, args ...interface{})

	Warning(args ...interface{})
	WarningFormat(format string, args ...interface{})

	Error(args ...interface{})
	ErrorFormat(format string, args ...interface{})

	Fatal(args ...interface{})
	FatalFormat(format string, args ...interface{})

	Panic(args ...interface{})
	PanicFormat(format string, args ...interface{})
}

type logger struct {
	log *log.Logger
}

var levelMap = map[string]log.Level{
	"trace": log.TraceLevel,
	"info":  log.InfoLevel,
	"warn":  log.WarnLevel,
	"error": log.ErrorLevel,
	"fatal": log.FatalLevel,
	"panic": log.PanicLevel,
}

func NewLogger() Logger {
	return &logger{log.New()}
}

func (l *logger) Init(debug bool, level string) {
	if !debug {
		l.log.SetFormatter(&log.JSONFormatter{})
	}

	logLevel := l.getLoggerLevel(level)
	l.log.SetLevel(logLevel)
}

func (l *logger) getLoggerLevel(level string) log.Level {
	lvl, ok := levelMap[level]
	if !ok {
		return log.DebugLevel
	}

	return lvl
}

func (l *logger) DebugFormat(format string, args ...interface{}) {
	l.log.Debugf(format, args...)
}

func (l *logger) Info(args ...interface{}) {
	l.log.Info(args...)
}

func (l *logger) InfoFormat(format string, args ...interface{}) {
	l.log.Infof(format, args...)
}

func (l *logger) Warning(args ...interface{}) {
	l.log.Warn(args...)
}

func (l *logger) WarningFormat(format string, args ...interface{}) {
	l.log.Warnf(format, args...)
}

func (l *logger) Error(args ...interface{}) {
	l.log.Error(args...)
}

func (l *logger) ErrorFormat(format string, args ...interface{}) {
	l.log.Errorf(format, args...)
}

func (l *logger) Fatal(args ...interface{}) {
	l.log.Fatal(args...)
}

func (l *logger) FatalFormat(format string, args ...interface{}) {
	l.log.Fatalf(format, args...)
}

func (l *logger) Panic(args ...interface{}) {
	l.log.Panic(args...)
}

func (l *logger) PanicFormat(format string, args ...interface{}) {
	l.log.Panicf(format, args...)
}
