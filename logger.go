// Copyright 2021 Harran Ali <harran.m@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"io"
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

type loggerDriver interface {
	GetTarget() interface{}
}

type LogFileDriver struct {
	FilePath string
}

type LogNullDriver struct{}

func (n *LogNullDriver) GetTarget() interface{} {
	return nil
}

func (f *LogFileDriver) GetTarget() interface{} {
	return f.FilePath
}

func NewLogger(driver loggerDriver) *Logger {
	if driver.GetTarget() == nil {
		l = &Logger{
			infoLogger:    log.New(io.Discard, "info: ", log.LstdFlags),
			warningLogger: log.New(io.Discard, "warning: ", log.LstdFlags),
			errorLogger:   log.New(io.Discard, "error: ", log.LstdFlags),
			debugLogger:   log.New(io.Discard, "debug: ", log.LstdFlags),
		}
		return l
	}
	path, ok := driver.GetTarget().(string)
	if !ok {
		panic("something wrong with the file path")
	}
	logsFile, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	l = &Logger{
		infoLogger:    log.New(logsFile, "info: ", log.LstdFlags),
		warningLogger: log.New(logsFile, "warning: ", log.LstdFlags),
		errorLogger:   log.New(logsFile, "error: ", log.LstdFlags),
		debugLogger:   log.New(logsFile, "debug: ", log.LstdFlags),
	}
	return l
}

func ResolveLogger() *Logger {
	return l
}

func (l *Logger) Info(msg interface{}) {
	l.infoLogger.Println(msg)
}

func (l *Logger) Debug(msg interface{}) {
	l.debugLogger.Println(msg)
}

func (l *Logger) Warning(msg interface{}) {
	l.warningLogger.Println(msg)
}

func (l *Logger) Error(msg interface{}) {
	l.errorLogger.Println(msg)
}
