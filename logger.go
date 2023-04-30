// Copyright 2021 Harran Ali <harran.m@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"log"
	"os"
)

// logs file
var logsFile *os.File

type Logger struct {
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
	debugLogger   *log.Logger
}

var l *Logger

func NewLogger(logsFilePath string) *Logger {
	logsFile, err := os.OpenFile(logsFilePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}

	infoLogger := log.New(logsFile, "info: ", log.LstdFlags)
	warningLogger := log.New(logsFile, "warning: ", log.LstdFlags)
	errorLogger := log.New(logsFile, "error: ", log.LstdFlags)
	debugLogger := log.New(logsFile, "debug: ", log.LstdFlags)
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
