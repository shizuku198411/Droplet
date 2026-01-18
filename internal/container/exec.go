package container

import (
	"droplet/internal/status"
	"droplet/internal/utils"
	"fmt"
	"os"
	"slices"
	"strconv"
)

// NewContainerExec constructs a ContainerExec with the default
// implementations of its dependencies (CommandFactory, StatusManager).
// This acts as the main entry point for the `exec` workflow, which
// runs an additional process inside an existing container.
func NewContainerExec() *ContainerExec {
	return &ContainerExec{
		commandFactory:         utils.NewCommandFactory(),
		containerStatusManager: status.NewStatusHandler(),
		syscallHandler:         utils.NewSyscallHandler(),
	}
}

// ContainerExec orchestrates execution of a new process inside an
// already-running container.
//
// It is responsible for:
//   - Verifying the container is in the RUNNING state
//   - Resolving the container’s init process PID
//   - Entering the container namespaces via nsenter
//   - Executing the requested command (optionally in interactive mode)
//
// Responsibility for low-level execution details is delegated to
// its collaborators to keep the workflow testable.
type ContainerExec struct {
	commandFactory         utils.CommandFactory
	containerStatusManager status.ContainerStatusManager
	syscallHandler         utils.KernelSyscallHandler
}

// Exec runs the given entrypoint inside the target container.
//
// The workflow is:
//  1. Verify that the container is RUNNING
//  2. Look up the container’s PID from state.json
//  3. Construct an nsenter invocation targeting that PID and namespaces
//  4. Start the command
//  5. If interactive mode is enabled, attach stdio and wait for completion
//
// If any step fails, execution stops and the error is returned.
func (c *ContainerExec) Exec(opt ExecOption) error {
	// 1. check container status
	//    if status is not running, return error
	containerStatus, containerStatusErr := c.containerStatusManager.GetStatusFromId(opt.ContainerId)
	if containerStatusErr != nil {
		return containerStatusErr
	}
	if containerStatus != status.RUNNING {
		return fmt.Errorf("container: %s not running.", opt.ContainerId)
	}

	// 2. retrieve pid from state.json
	containerPid, containerPidErr := c.containerStatusManager.GetPidFromId(opt.ContainerId)
	if containerPidErr != nil {
		return containerPidErr
	}

	// 3. prepare entrypoint with nsenter
	if opt.Tty {
		if err := c.executeShim(containerPid, opt); err != nil {
			return err
		}
	} else {
		if err := c.executeNsenter(containerPid, opt); err != nil {
			return err
		}
	}

	return nil
}

func (c *ContainerExec) executeNsenter(containerPid int, opt ExecOption) error {
	nsenterCommand := []string{"nsenter", "-t", strconv.Itoa(containerPid), "--all"}
	commandStr := slices.Concat(nsenterCommand, opt.Entrypoint)
	cmd := c.commandFactory.Command(commandStr[0], commandStr[1:]...)
	// set stdout/stderr to log files
	logPath := utils.ExecLogPath(opt.ContainerId)
	f, err := c.syscallHandler.OpenFile(logPath, os.O_CREATE|os.O_WRONLY, 0640)
	if err != nil {
		return err
	}
	cmd.SetStdout(f)
	cmd.SetStderr(f)

	// execute entrypoint
	if err := cmd.Start(); err != nil {
		return err
	}

	return nil
}

func (c *ContainerExec) executeShim(containerPid int, opt ExecOption) error {
	entrypoint := opt.Entrypoint
	shimArgs := append([]string{"exec-shim", opt.ContainerId, strconv.Itoa(containerPid)}, entrypoint...)
	cmd := c.commandFactory.Command(os.Args[0], shimArgs...)

	// execute exec-shim subcommand
	if err := cmd.Start(); err != nil {
		return err
	}

	return nil
}
