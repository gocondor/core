package core

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/google/uuid"
)

func TestNewEventsManager(t *testing.T) {
	m := NewEventsManager()
	if fmt.Sprintf("%T", m) != "*core.EventsManager" {
		t.Errorf("failed testing new events manager")
	}
}

func TestResolveEventsManager(t *testing.T) {
	NewEventsManager()
	m := ResolveEventsManager()
	if fmt.Sprintf("%T", m) != "*core.EventsManager" {
		t.Errorf("failed testing new events manager")
	}
}

func TestEvents(t *testing.T) {
	pwd, _ := os.Getwd()
	const eventName1 string = "test-event-name1"
	const eventName2 string = "test-event-name2"
	var tmpDir string
	if runtime.GOOS == "linux" {
		tmpDir = t.TempDir()
	} else {
		tmpDir = filepath.Join(pwd, "/testingdata/tmp")
	}
	tmpFile1 := filepath.Join(tmpDir, uuid.NewString())
	tmpFile2 := filepath.Join(tmpDir, uuid.NewString())
	tmpFile3 := filepath.Join(tmpDir, uuid.NewString())
	m := NewEventsManager()
	m.Register(eventName1, func(event *Event, requestContext *Context) {
		os.Create(tmpFile1)
		f, err := os.Create(tmpFile1)
		if err != nil {
			t.Errorf("error testing register event: %v", err.Error())
		}
		f.WriteString(event.Name)
		f.Close()
	})
	m.Register(eventName1, func(event *Event, requestContext *Context) {
		os.Create(tmpFile3)
		f, err := os.Create(tmpFile3)
		if err != nil {
			t.Errorf("error testing register event: %v", err.Error())
		}
		f.WriteString(event.Name)
		f.Close()
	})
	m.Fire(&Event{Name: eventName1})
	m.executeEventsJobs()

	ff, err := os.Open(tmpFile1)
	if err != nil {
		t.Errorf("error testing register event : %v", err.Error())
	}

	d, err := io.ReadAll(ff)
	if string(d) != eventName1 {
		t.Error("faild testing events")
	}
	ff.Close()
	os.Remove(tmpFile1)

	ff, err = os.Open(tmpFile3)
	if err != nil {
		t.Errorf("error testing register event : %v", err.Error())
	}

	d, err = io.ReadAll(ff)
	if string(d) != eventName1 {
		t.Error("faild testing events")
	}
	ff.Close()
	os.Remove(tmpFile3)

	m.Register(eventName2, func(event *Event, requestContext *Context) {
		f, err := os.Create(tmpFile2)
		if err != nil {
			t.Errorf("error testing register event: %v", err.Error())
		}
		f.WriteString(event.Name)
		f.Close()
	})
	m.Fire(&Event{Name: eventName2})
	m.executeEventsJobs()

	ff, err = os.Open(tmpFile2)
	if err != nil {
		t.Errorf("error testing register event : %v", err.Error())
	}

	d, err = io.ReadAll(ff)
	if string(d) != eventName2 {
		t.Error("faild testing events")
	}
	ff.Close()
}
