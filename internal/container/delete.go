package container

import (
	"droplet/internal/hook"
	"droplet/internal/status"
	"fmt"
)

func NewContainerDelete() *ContainerDelete {
	return &ContainerDelete{
		specLoader:              newFileSpecLoader(),
		containerStatusManager:  status.NewStatusHandler(),
		containerHookController: hook.NewHookController(),
	}
}

type ContainerDelete struct {
	specLoader              specLoader
	containerStatusManager  status.ContainerStatusManager
	containerHookController hook.ContainerHookController
}

func (c *ContainerDelete) Delete(opt DeleteOption) error {
	// check container status
	// if status is running, return error
	containerStatus, err := c.containerStatusManager.GetStatusFromId(opt.ContainerId)
	if err != nil {
		return err
	}
	if containerStatus == status.RUNNING {
		return fmt.Errorf("container: %s is running.", opt.ContainerId)
	}

	// load config.json
	spec, err := c.specLoader.loadFile(opt.ContainerId)
	if err != nil {
		return err
	}

	// HOOK: poststop
	if err := c.containerHookController.RunPoststopHooks(
		opt.ContainerId,
		spec.Hooks.Poststop,
	); err != nil {
		return err
	}

	// remove state.json
	if err := c.containerStatusManager.RemoveStatusFile(opt.ContainerId); err != nil {
		return err
	}

	return nil
}
