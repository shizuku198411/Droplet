package container

import (
	"os"
	"syscall"
)

func InitContainer(opt InitOption) error {
	fifo := opt.Fifo
	entrypoint := opt.Entrypoint

	// read fifo for waiting start signal
	if err := readFifo(fifo); err != nil {
		return err
	}

	if err := syscall.Exec(entrypoint[0], entrypoint, os.Environ()); err != nil {
		return err
	}

	return nil
}
