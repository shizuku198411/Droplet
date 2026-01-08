package container

import (
	"droplet/internal/hook"
	"droplet/internal/status"
	"droplet/internal/utils"
	"fmt"
)

// NewContainerDelete constructs a ContainerDelete with the default
// implementations of its dependencies (SpecLoader, StatusManager, HookController).
// This acts as the main entry point for the container deletion workflow.
func NewContainerDelete() *ContainerDelete {
	return &ContainerDelete{
		specLoader:              newFileSpecLoader(),
		fifoHandler:             newContainerFifoHandler(),
		containerStatusManager:  status.NewStatusHandler(),
		containerHookController: hook.NewHookController(),
	}
}

// ContainerDelete orchestrates the container deletion flow.
//
// It is responsible for:
//   - Validating the current container status
//   - Loading the OCI spec (for hooks)
//   - Executing poststop hooks
//   - Removing the container state file
//
// Low-level operations are delegated to its collaborators so that
// the logic can be tested and substituted.
type ContainerDelete struct {
	specLoader  specLoader
	fifoHandler interface {
		removeFifo(path string) error
	}
	containerStatusManager  status.ContainerStatusManager
	containerHookController hook.ContainerHookController
}

// Delete executes the container deletion pipeline for the given container ID.
//
// The workflow is:
//  1. Check the container status and fail if it is still running
//  2. Load the OCI spec (config.json)
//  3. Run poststop hooks
//  4. Remove the container state file (state.json)
//  5. Remove the FIFO if the container status is created
//
// If any step fails, the error is returned immediately and subsequent
// steps are not executed.
func (c *ContainerDelete) Delete(opt DeleteOption) error {
	// 1. check container status
	//    if status is running, return error
	containerStatus, err := c.containerStatusManager.GetStatusFromId(opt.ContainerId)
	if err != nil {
		return err
	}
	if containerStatus == status.RUNNING {
		return fmt.Errorf("container: %s is not stopped. current status: %s", opt.ContainerId, containerStatus)
	}

	// 2. load config.json
	spec, err := c.specLoader.loadFile(opt.ContainerId)
	if err != nil {
		return err
	}

	// 3. HOOK: poststop
	if err := c.containerHookController.RunPoststopHooks(
		opt.ContainerId,
		spec.Hooks.Poststop,
	); err != nil {
		return err
	}

	// 4. remove state.json
	if err := c.containerStatusManager.RemoveStatusFile(opt.ContainerId); err != nil {
		return err
	}

	// 5. remove exec.fifo if status is created
	if containerStatus == status.CREATED {
		if err := c.fifoHandler.removeFifo(utils.FifoPath(opt.ContainerId)); err != nil {
			return err
		}
	}

	return nil
}
