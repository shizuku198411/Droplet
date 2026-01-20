package container

import (
	"droplet/internal/hook"
	"droplet/internal/logs"
	"droplet/internal/spec"
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
func (c *ContainerKill) Kill(opt KillOption) (err error) {
	var (
		spec  spec.Spec
		event = "kill"
		stage string
		pid   int
	)

	// audit log
	defer func() {
		result := "success"
		if err != nil {
			result = "fail"
		}
		_ = logs.RecordAuditLog(logs.AuditRecord{
			ContainerId: opt.ContainerId,
			Event:       event,
			Stage:       stage,
			Spec:        &spec,
			Pid:         pid,
			Result:      result,
			Error:       err,
		})
	}()

	// 1. check container status
	//    if status is not running, return error
	stage = "get_status"
	containerStatus, err := c.containerStatusManager.GetStatusFromId(opt.ContainerId)
	if err != nil {
		return err
	}

	stage = "check_status"
	if containerStatus != status.RUNNING {
		return fmt.Errorf("container: %s not running.", opt.ContainerId)
	}

	// 2. retrieve pid and shimpid from state.json
	stage = "get_pid"
	containerPid, err := c.containerStatusManager.GetPidFromId(opt.ContainerId)
	if err != nil {
		return err
	}
	stage = "get_shim_pid"
	shimPid, err := c.containerStatusManager.GetShimPidFromId(opt.ContainerId)
	if err != nil {
		return err
	}

	// 3. send signal to pid
	stage = "send_signal"
	err = c.syscallHandler.Kill(containerPid, signalMap[opt.Signal])
	if err != nil {
		return err
	}

	// if shim pid > 0, the container created with interactive mode
	// clean up files for shim
	stage = "cleanup_shim"
	if shimPid > 0 {
		err = c.cleanupShim(opt.ContainerId)
		if err != nil {
			return err
		}
	}

	// 4. update status file
	//      status = stopped
	//      pid = 0
	//		shimPid = 0
	stage = "update_state"
	err = c.containerStatusManager.UpdateStatus(
		opt.ContainerId,
		status.STOPPED,
		0,
		0,
	)
	if err != nil {
		return err
	}

	// 5. load config.json
	stage = "load_spec"
	spec, err = c.specLoader.loadFile(opt.ContainerId)
	if err != nil {
		return err
	}

	// 6. HOOK: stopContainer
	stage = "hook_stopContainer"
	err = c.containerHookController.RunStopContainerHooks(
		opt.ContainerId,
		spec.Hooks.StopContainer,
	)
	if err != nil {
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
