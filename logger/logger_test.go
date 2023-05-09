package logger

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestInfo(t *testing.T) {
	path := filepath.Join(t.TempDir(), uuid.NewString())
	l := NewLogger(&LogFileDriver{
		FilePath: path,
	})
	l.Info("DFT2V56H")
	lf, err := os.Open(path)
	if err != nil {
		t.Error("failed testing info")
	}
	d, err := io.ReadAll(lf)
	if err != nil {
		t.Error("error testing info")
	}
	if !strings.Contains(string(d), "DFT2V56H") {
		t.Error("error testing info")
	}
	t.Cleanup(func() {
		CloseLogsFile()
	})
}

func TestWarning(t *testing.T) {
	path := filepath.Join(t.TempDir(), uuid.NewString())
	l := NewLogger(&LogFileDriver{
		FilePath: path,
	})

	l.Warning("DFT2V56H")
	lf, err := os.Open(path)
	if err != nil {
		t.Error("failed testing warning")
	}
	d, err := io.ReadAll(lf)
	if err != nil {
		t.Error("failed testing warning")
	}
	if !strings.Contains(string(d), "DFT2V56H") {
		t.Error("failed testing warning")
	}
	t.Cleanup(func() {
		CloseLogsFile()
	})
}

func TestDebug(t *testing.T) {
	path := filepath.Join(t.TempDir(), uuid.NewString())
	l := NewLogger(&LogFileDriver{
		FilePath: path,
	})
	l.Debug("DFT2V56H")
	lf, err := os.Open(path)
	if err != nil {
		t.Error("failed testing debug")
	}
	d, err := io.ReadAll(lf)
	if err != nil {
		t.Error("error testing debug")
	}
	if !strings.Contains(string(d), "DFT2V56H") {
		t.Error("error testing debug")
	}
	t.Cleanup(func() {
		CloseLogsFile()
	})
}

func TestError(t *testing.T) {
	path := filepath.Join(t.TempDir(), uuid.NewString())
	l := NewLogger(&LogFileDriver{
		FilePath: path,
	})
	l.Error("DFT2V56H")
	lf, err := os.Open(path)
	if err != nil {
		t.Error("failed testing error")
	}
	d, err := io.ReadAll(lf)
	if err != nil {
		t.Error("failed testing error")
	}
	if !strings.Contains(string(d), "DFT2V56H") {
		t.Error("failed testing error")
	}
	t.Cleanup(func() {
		CloseLogsFile()
	})
}
