package container

import (
	"droplet/internal/status"
)

// NewContainerStart returns a ContainerStart wired with the default
// FIFO handler implementation. This is the standard entry point for
// executing the container start phase.
func NewContainerStart() *ContainerStart {
	return &ContainerStart{
		fifoHandler:            newContainerFifoHandler(),
		containerStatusManager: status.NewStatusHandler(),
	}
}

// ContainerStart coordinates the logic for starting a container
// from the runtime side.
//
// The start phase signals the already-created init process by
// writing to the FIFO and then removes the FIFO after the signal
// is delivered.
type ContainerStart struct {
	fifoHandler interface {
		writeFifo(path string) error
		removeFifo(path string) error
	}
	containerStatusManager status.ContainerStatusManager
}

// Execute performs the container start sequence for the given container.
//
// The sequence is:
//
//  1. Open and write to the FIFO to notify the init process that it may start
//  2. Remove the FIFO after the notification is complete
//
// An error is returned if either the write or removal operation fails.
func (c *ContainerStart) Execute(opt StartOption) error {
	fifo := fifoPath(opt.ContainerId)

	// write fifo
	if err := c.fifoHandler.writeFifo(fifo); err != nil {
		return err
	}

	// remove fifo
	if err := c.fifoHandler.removeFifo(fifo); err != nil {
		return err
	}

	// update status file
	//   status = running
	if err := c.containerStatusManager.UpdateStatus(
		containerDir(opt.ContainerId),
		opt.ContainerId,
		status.RUNNING,
		-1, // no update
	); err != nil {
		return err
	}

	return nil
}
