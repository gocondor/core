package core

import (
	"log"
	"os"
)

type Logger struct {
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
	debugLogger   *log.Logger
}

var l *Logger

func NewLogger() *Logger {
	opts := log.Lshortfile
	infoLogger := log.New(os.Stdout, "info: ", opts)
	warningLogger := log.New(os.Stdout, "warning: ", opts)
	errorLogger := log.New(os.Stdout, "error: ", opts)
	debugLogger := log.New(os.Stdout, "debug: ", opts)
	l = &Logger{
		infoLogger,
		warningLogger,
		errorLogger,
		debugLogger,
	}

	return l
}

func ResolveLogger() *Logger {
	return l
}

func (l *Logger) Info(msg string) {
	l.infoLogger.Println(msg)
}

func (l *Logger) Debug(msg string) {
	l.debugLogger.Println(msg)
}

func (l *Logger) Warning(msg string) {
	l.warningLogger.Println(msg)
}

func (l *Logger) Error(msg string) {
	l.errorLogger.Println(msg)
}
