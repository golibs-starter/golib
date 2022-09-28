package log

import (
	"bytes"
	"testing"
)

// TestingWriter is a WriteSyncer that writes to the given testing.TB.
type TestingWriter struct {
	tb testing.TB
}

func NewTestingWriter(tb testing.TB) TestingWriter {
	return TestingWriter{tb: tb}
}

func (w TestingWriter) Write(p []byte) (n int, err error) {
	n = len(p)

	// Strip trailing newline because t.Log always adds one.
	p = bytes.TrimRight(p, "\n")

	// Note: tb.Log is safe for concurrent use.
	w.tb.Logf("%s", p)
	return n, nil
}

func (w TestingWriter) Sync() error {
	return nil
}
