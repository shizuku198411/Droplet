package command

import (
	"bytes"
	"os"
	"testing"
)

func readStdout(t *testing.T, fnc func()) string {
	t.Helper()

	// store original stdout
	orgStdOut := os.Stdout
	// set original stdout to os.Stdout
	defer func() { os.Stdout = orgStdOut }()
	// create pipe
	reader, writer, _ := os.Pipe()
	// set stdout to write
	os.Stdout = writer

	// execute func
	fnc()

	// close writer
	writer.Close()
	// read buffer
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(reader); err != nil {
		t.Fatalf("failed to read buffer: %v", err)
	}
	return buf.String()
}
