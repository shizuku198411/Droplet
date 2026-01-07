package container

import (
	"droplet/internal/hook"
	"droplet/internal/status"
	"droplet/internal/utils"
)

// NewContainerStart returns a ContainerStart wired with the default
// FIFO handler implementation. This is the standard entry point for
// executing the container start phase.
func NewContainerStart() *ContainerStart {
	return &ContainerStart{
		specLoader:              newFileSpecLoader(),
		fifoHandler:             newContainerFifoHandler(),
		containerStatusManager:  status.NewStatusHandler(),
		containerHookController: hook.NewHookController(),
	}
}

// ContainerStart coordinates the logic for starting a container
// from the runtime side.
//
// The start phase signals the already-created init process by
// writing to the FIFO and then removes the FIFO after the signal
// is delivered.
type ContainerStart struct {
	specLoader  specLoader
	fifoHandler interface {
		writeFifo(path string) error
		removeFifo(path string) error
	}
	containerStatusManager  status.ContainerStatusManager
	containerHookController hook.ContainerHookController
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
	// load config.json
	spec, err := c.specLoader.loadFile(opt.ContainerId)
	if err != nil {
		return err
	}

	// HOOK: startContainer
	if err := c.containerHookController.RunStartContainerHooks(
		opt.ContainerId,
		spec.Hooks.StartContainer,
	); err != nil {
		return err
	}

	// write fifo
	fifo := utils.FifoPath(opt.ContainerId)
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
		opt.ContainerId,
		status.RUNNING,
		-1, // no update
	); err != nil {
		return err
	}

	// HOOK: poststart
	if err := c.containerHookController.RunPoststartHooks(
		opt.ContainerId,
		spec.Hooks.Poststart,
	); err != nil {
		return err
	}

	return nil
}
