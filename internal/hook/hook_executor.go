package hook

import (
	"bytes"
	"droplet/internal/spec"
	"droplet/internal/status"
	"droplet/internal/utils"
	"fmt"
	"os"
	"strconv"
)

type ContainerHookController interface {
	RunCreateRuntimeHooks(containerId string, hookList []spec.HookObject) error
	RunCreateContainerHooks(containerId string, hookList []spec.HookObject) error
	RunStartContainerHooks(containerId string, hookList []spec.HookObject) error
	RunPoststartHooks(containerId string, hookList []spec.HookObject) error
	RunPoststopHooks(containerId string, hookList []spec.HookObject) error
}

func NewHookController() *HookController {
	return &HookController{
		commandFactory:         utils.NewCommandFactory(),
		containerStatusManager: status.NewStatusHandler(),
	}
}

type HookController struct {
	commandFactory         utils.CommandFactory
	containerStatusManager status.ContainerStatusManager
}

func (c *HookController) RunCreateRuntimeHooks(containerId string, hookList []spec.HookObject) error {
	if hookList == nil || len(hookList) == 0 {
		return nil
	}
	return c.runHookList(containerId, "createRuntime", hookList)
}

func (c *HookController) RunCreateContainerHooks(containerId string, hookList []spec.HookObject) error {
	if hookList == nil || len(hookList) == 0 {
		return nil
	}
	return c.runHookListWithNsenter(containerId, "createContainer", hookList)
}

func (c *HookController) RunStartContainerHooks(containerId string, hookList []spec.HookObject) error {
	if hookList == nil || len(hookList) == 0 {
		return nil
	}
	return c.runHookListWithNsenter(containerId, "startContainer", hookList)
}

func (c *HookController) RunPoststartHooks(containerId string, hookList []spec.HookObject) error {
	if hookList == nil || len(hookList) == 0 {
		return nil
	}
	return c.runHookList(containerId, "poststart", hookList)
}

func (c *HookController) RunPoststopHooks(containerId string, hookList []spec.HookObject) error {
	if hookList == nil || len(hookList) == 0 {
		return nil
	}
	return c.runHookList(containerId, "poststop", hookList)
}

func (c *HookController) runHookList(containerId string, phase string, hookList []spec.HookObject) error {
	// read state.json
	stateJson, err := c.containerStatusManager.ReadStatusFile(containerId)
	if err != nil {
		return err
	}

	for i, hook := range hookList {
		if hook.Path == "" {
			return fmt.Errorf("hook %s[%d]: empty path", phase, i)
		}

		// set args
		args := hook.Args
		if len(args) == 0 {
			args = []string{}
		}

		// prepare hook environment
		cmd := c.commandFactory.Command(hook.Path, args...)
		cmd.SetEnv(append(os.Environ(), hook.Env...))
		cmd.SetStdin(bytes.NewReader([]byte(stateJson)))
		cmd.SetStdout(os.Stdout)
		cmd.SetStderr(os.Stderr)

		// execute hook
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("hook %s[%d] failed: %w", phase, i, err)
		}
	}
	return nil
}

func (c *HookController) runHookListWithNsenter(containerId string, phase string, hookList []spec.HookObject) error {
	// read state.json
	stateJson, err := c.containerStatusManager.ReadStatusFile(containerId)
	if err != nil {
		return err
	}

	// get pid
	initPid, err := c.containerStatusManager.GetPidFromId(containerId)
	if err != nil {
		return err
	}

	for i, hook := range hookList {
		if hook.Path == "" {
			return fmt.Errorf("hook %s[%d]: empty path", phase, i)
		}

		// set args
		args := hook.Args
		if len(args) == 0 {
			args = []string{}
		}

		// prepare hook environment with nsenter
		nsenterArgs := []string{
			"nsenter",
			"-t", strconv.Itoa(initPid),
			"-m", "-u", "-i", "-n", "-p",
			"--",
			hook.Path,
		}
		nsenterArgs = append(nsenterArgs, args...)

		cmd := c.commandFactory.Command("/usr/bin/nsenter", nsenterArgs...)
		cmd.SetEnv(append(os.Environ(), hook.Env...))
		cmd.SetStdin(bytes.NewReader([]byte(stateJson)))
		cmd.SetStdout(os.Stdout)
		cmd.SetStderr(os.Stderr)

		// execute hook
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("hook %s[%d] failed: %w", phase, i, err)
		}
	}
	return nil
}
