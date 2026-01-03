package container

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTestFifo(t *testing.T) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "exec.fifo")
	if err := syscall.Mkfifo(path, 0o600); err != nil {
		t.Fatalf("create test FIFO failed")
	}

	return path
}

func TestNewContainerFifoHandler_Success(t *testing.T) {
	// == arrange ==

	// == act ==
	containerFifoHandler := newContainerFifoHandler()

	// == assert ==
	assert.NotNil(t, containerFifoHandler)
}

func TestCreateFifo_Success(t *testing.T) {
	// == arrange ==
	containerFifoHandler := &containerFifoHandler{}
	path := filepath.Join(t.TempDir(), "exec.fifo")

	// == act ==
	err := containerFifoHandler.createFifo(path)

	// == assert ==
	assert.Nil(t, err)
}

func TestCreateFifo_PathNotExistsError(t *testing.T) {
	// == arrange ==
	containerFifoHandler := &containerFifoHandler{}
	path := "/not/exists/path"

	// == act ==
	err := containerFifoHandler.createFifo(path)

	// == assert ==
	assert.NotNil(t, err)
}

func TestRemoveFifo_Success(t *testing.T) {
	// == arrange ==
	path := createTestFifo(t)
	containerFifoHandler := &containerFifoHandler{}

	// == act ==
	err := containerFifoHandler.removeFifo(path)

	// == assert ==
	assert.Nil(t, err)
}

func TestRemoveFifo_FileNotExistsError(t *testing.T) {
	// == arrange ==
	containerFifoHandler := &containerFifoHandler{}
	path := "/not/exists/path"

	// == act ==
	err := containerFifoHandler.removeFifo(path)

	// == assert ==
	assert.NotNil(t, err)
}

func TestReadFifo_Success(t *testing.T) {
	// == arrange ==
	path := createTestFifo(t)
	containerFifoHandler := &containerFifoHandler{}

	// writer goroutine
	writerErrCh := make(chan error, 1)
	go func() {
		f, err := os.OpenFile(path, os.O_WRONLY, 0)
		if err != nil {
			writerErrCh <- fmt.Errorf("open writer failed: %w", err)
			return
		}
		defer f.Close()

		if _, err := f.Write([]byte{0}); err != nil {
			writerErrCh <- fmt.Errorf("write failed: %w", err)
			return
		}
		writerErrCh <- nil
	}()

	// == act ==
	err := containerFifoHandler.readFifo(path)

	// == assert ==
	assert.Nil(t, err)
}

func TestReadFifo_FileNotEsistsError(t *testing.T) {
	// == arrange ==
	path := "/not/exists/path"
	containerFifoHandler := &containerFifoHandler{}

	// == act ==
	err := containerFifoHandler.readFifo(path)

	// == assert ==
	assert.NotNil(t, err)
}

func TestReadFifo_FileOpenError(t *testing.T) {
	// == arrange ==
	path := t.TempDir()
	containerFifoHandler := &containerFifoHandler{}

	// == act ==
	err := containerFifoHandler.readFifo(path)

	// == assert ==
	assert.NotNil(t, err)
}

func TestWriteFifo_Success(t *testing.T) {
	// == arrange ==
	path := createTestFifo(t)
	containerFifoHandler := &containerFifoHandler{}

	// reader goroutine
	readerErrCh := make(chan error, 1)
	go func() {
		f, err := os.OpenFile(path, os.O_RDONLY, 0)
		if err != nil {
			readerErrCh <- fmt.Errorf("open reader failed: %w", err)
			return
		}
		defer f.Close()

		buf := make([]byte, 1)
		if _, err := f.Read(buf); err != nil {
			readerErrCh <- fmt.Errorf("read failed: %w", err)
			return
		}
		readerErrCh <- nil
	}()

	// == act ==
	err := containerFifoHandler.writeFifo(path)

	// == assert ==
	assert.Nil(t, err)
}

func TestWriteFifo_FileNotExistsError(t *testing.T) {
	// == arrange ==
	path := "/not/exists/path"
	containerFifoHandler := &containerFifoHandler{}

	// == act ==
	err := containerFifoHandler.writeFifo(path)

	// == assert ==
	assert.NotNil(t, err)
}
