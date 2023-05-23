package logger

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestNewLogger(t *testing.T) {
	fp := filepath.Join(t.TempDir(), uuid.NewString())
	f, err := os.Create(fp)
	if err != nil {
		t.Errorf("failed test new logger")
	}
	f.Close()
	l := NewLogger(&LogNullDriver{})
	l.Info("testing")
	fdrv := &LogFileDriver{
		fp,
	}
	trgt := fdrv.GetTarget()
	ts, ok := trgt.(string)
	if !ok {
		t.Errorf("failed test new logger")
	}
	if ts != fp {
		t.Errorf("failed test new logger")
	}
	l = NewLogger(fdrv)
	l.Error("test-err")
	f, err = os.Open(fp)
	if err != nil {
		t.Errorf("failed test new logger")
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if !strings.Contains(string(b), "test-err") {
		t.Errorf("failed test new logger")
	}
}

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
