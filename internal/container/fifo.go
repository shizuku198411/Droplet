package container

import (
	"os"
	"syscall"
)

func createFifo(path string) error {
	if err := syscall.Mkfifo(path, 0o600); err != nil {
		return err
	}

	return nil
}

func removeFifo(path string) error {
	if err := os.Remove(path); err != nil {
		return err
	}

	return nil
}

func readFifo(path string) error {
	f, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer f.Close()

	buf := make([]byte, 1)
	if _, err := f.Read(buf); err != nil {
		return err
	}

	return nil
}

func writeFifo(path string) error {
	f, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Write([]byte{1}); err != nil {
		return err
	}

	return nil
}
