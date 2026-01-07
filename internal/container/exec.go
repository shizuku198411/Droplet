package container

import (
	"droplet/internal/status"
	"droplet/internal/utils"
	"fmt"
	"os"
	"slices"
	"strconv"
)

func NewContainerExec() *ContainerExec {
	return &ContainerExec{
		commandFactory:         utils.NewCommandFactory(),
		containerStatusManager: status.NewStatusHandler(),
	}
}

type ContainerExec struct {
	commandFactory         utils.CommandFactory
	containerStatusManager status.ContainerStatusManager
}

func (c *ContainerExec) Exec(opt ExecOption) error {
	// check container status
	// if status is not running, return error
	containerStatus, containerStatusErr := c.containerStatusManager.GetStatusFromId(opt.ContainerId)
	if containerStatusErr != nil {
		return containerStatusErr
	}
	if containerStatus != status.RUNNING {
		return fmt.Errorf("container: %s not running.", opt.ContainerId)
	}

	// retrieve pid
	containerPid, containerPidErr := c.containerStatusManager.GetPidFromId(opt.ContainerId)
	if containerPidErr != nil {
		return containerPidErr
	}

	// prepare entrypoint with nsenter
	nsenterCommand := []string{"nsenter", "-t", strconv.Itoa(containerPid), "--all"}
	commandStr := slices.Concat(nsenterCommand, opt.Entrypoint)
	cmd := c.commandFactory.Command(commandStr[0], commandStr[1:]...)
	if opt.Interactive {
		cmd.SetStdout(os.Stdout)
		cmd.SetStderr(os.Stderr)
		cmd.SetStdin(os.Stdin)
	}

	// execute entrypoint
	if err := cmd.Start(); err != nil {
		return err
	}

	// wait entrypoint if exec in interactive
	if opt.Interactive {
		if err := cmd.Wait(); err != nil {
			return err
		}
	}

	return nil
}
