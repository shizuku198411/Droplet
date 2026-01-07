package container

import (
	"droplet/internal/status"
	"fmt"
)

func NewContainerKill() *ContainerKill {
	return &ContainerKill{
		syscallHandler:         newSyscallHandler(),
		containerStatusManager: status.NewStatusHandler(),
	}
}

type ContainerKill struct {
	syscallHandler         containerEnvPrepareSyscallHandler
	containerStatusManager status.ContainerStatusManager
}

func (c *ContainerKill) Kill(opt KillOption) error {
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

	// send signal to pid
	if err := c.syscallHandler.Kill(containerPid, signalMap[opt.Signal]); err != nil {
		return err
	}

	// update status file
	//   status = stopped
	//   pid = 0
	if err := c.containerStatusManager.UpdateStatus(
		opt.ContainerId,
		status.STOPPED,
		0,
	); err != nil {
		return err
	}

	return nil
}
