package container

import (
	"droplet/internal/hook"
	"droplet/internal/status"
	"droplet/internal/utils"
	"fmt"
)

// NewContainerKill constructs a ContainerKill with the default
// implementations of its dependencies (SyscallHandler, StatusManager).
// This serves as the main entry point for the `kill` workflow, which
// delivers a signal to a running container’s init process.
func NewContainerKill() *ContainerKill {
	return &ContainerKill{
		specLoader:              newFileSpecLoader(),
		syscallHandler:          utils.NewSyscallHandler(),
		containerStatusManager:  status.NewStatusHandler(),
		containerHookController: hook.NewHookController(),
	}
}

// ContainerKill orchestrates the container termination flow.
//
// It is responsible for:
//   - Verifying that the container is currently RUNNING
//   - Resolving the container’s init process PID from state.json
//   - Sending the requested signal to that process
//   - Updating the container status to STOPPED
//
// Low-level system interactions are delegated to collaborators to
// keep the workflow testable and replaceable.
type ContainerKill struct {
	specLoader              specLoader
	syscallHandler          utils.KernelSyscallHandler
	containerStatusManager  status.ContainerStatusManager
	containerHookController hook.ContainerHookController
}

// Kill sends a signal to the container’s init process and updates its state.
//
// The workflow is:
//  1. Check that the container is RUNNING
//  2. Retrieve the init PID from state.json
//  3. Send the configured signal to that PID
//  4. Update the status file to STOPPED and clear the PID
//
// If any step fails, the method stops and returns the error.
func (c *ContainerKill) Kill(opt KillOption) error {
	// 1. check container status
	//    if status is not running, return error
	containerStatus, containerStatusErr := c.containerStatusManager.GetStatusFromId(opt.ContainerId)
	if containerStatusErr != nil {
		return containerStatusErr
	}
	if containerStatus != status.RUNNING {
		return fmt.Errorf("container: %s not running.", opt.ContainerId)
	}

	// 2. retrieve pid and shimpid from state.json
	containerPid, containerPidErr := c.containerStatusManager.GetPidFromId(opt.ContainerId)
	if containerPidErr != nil {
		return containerPidErr
	}
	shimPid, shimPidErr := c.containerStatusManager.GetShimPidFromId(opt.ContainerId)
	if shimPidErr != nil {
		return shimPidErr
	}

	// 3. send signal to pid
	if err := c.syscallHandler.Kill(containerPid, signalMap[opt.Signal]); err != nil {
		return err
	}
	// if shim pid > 0, the container created with interactive mode
	// clean up files for shim
	if shimPid > 0 {
		if err := c.cleanupShim(opt.ContainerId); err != nil {
			return err
		}
	}

	// 4. update status file
	//      status = stopped
	//      pid = 0
	//		shimPid = 0
	if err := c.containerStatusManager.UpdateStatus(
		opt.ContainerId,
		status.STOPPED,
		0,
		0,
	); err != nil {
		return err
	}

	// 5. load config.json
	spec, err := c.specLoader.loadFile(opt.ContainerId)
	if err != nil {
		return err
	}

	// 6. HOOK: stopContainer
	if err := c.containerHookController.RunStopContainerHooks(
		opt.ContainerId,
		spec.Hooks.StopContainer,
	); err != nil {
		return err
	}

	return nil
}

func (c *ContainerKill) cleanupShim(containerId string) error {
	// remove tty.sock
	if err := c.syscallHandler.Remove(utils.SockPath(containerId)); err != nil {
		return err
	}
	// remove init.pid
	if err := c.syscallHandler.Remove(utils.InitPidFilePath(containerId)); err != nil {
		return err
	}
	return nil
}
