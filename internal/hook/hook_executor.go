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

// ContainerHookController defines the interface for executing lifecycle hooks
// at various container phases. Each method takes a container ID and a list
// of hook definitions to execute in order.
type ContainerHookController interface {
	RunCreateRuntimeHooks(containerId string, hookList []spec.HookObject) error
	RunCreateContainerHooks(containerId string, hookList []spec.HookObject) error
	RunStartContainerHooks(containerId string, hookList []spec.HookObject) error
	RunPoststartHooks(containerId string, hookList []spec.HookObject) error
	RunPoststopHooks(containerId string, hookList []spec.HookObject) error
}

// NewHookController constructs a HookController with the default
// implementations of its dependencies (CommandFactory, StatusManager).
// This is the default entry point for managing and executing container hooks.
func NewHookController() *HookController {
	return &HookController{
		commandFactory:         utils.NewCommandFactory(),
		containerStatusManager: status.NewStatusHandler(),
	}
}

// HookController is the default implementation of ContainerHookController.
//
// It is responsible for:
//   - Reading container state from state.json
//   - Preparing the environment and stdin for each hook process
//   - Optionally entering the container namespaces via nsenter
//   - Executing the hook commands in sequence
type HookController struct {
	commandFactory         utils.CommandFactory
	containerStatusManager status.ContainerStatusManager
}

// RunCreateRuntimeHooks executes the createRuntime hook list in the host
// namespaces. If the list is nil or empty, it is a no-op.
//
// This corresponds to the OCI createRuntime lifecycle phase.
func (c *HookController) RunCreateRuntimeHooks(containerId string, hookList []spec.HookObject) error {
	if hookList == nil || len(hookList) == 0 {
		return nil
	}
	return c.runHookList(containerId, "createRuntime", hookList)
}

// RunCreateContainerHooks executes the createContainer hook list in the
// container's namespaces using nsenter. If the list is nil or empty,
// it is a no-op.
//
// This corresponds to the OCI createContainer lifecycle phase.
func (c *HookController) RunCreateContainerHooks(containerId string, hookList []spec.HookObject) error {
	if hookList == nil || len(hookList) == 0 {
		return nil
	}
	return c.runHookListWithNsenter(containerId, "createContainer", hookList)
}

// RunStartContainerHooks executes the startContainer hook list in the
// container's namespaces using nsenter. If the list is nil or empty,
// it is a no-op.
//
// This corresponds to the OCI startContainer lifecycle phase.
func (c *HookController) RunStartContainerHooks(containerId string, hookList []spec.HookObject) error {
	if hookList == nil || len(hookList) == 0 {
		return nil
	}
	return c.runHookListWithNsenter(containerId, "startContainer", hookList)
}

// RunPoststartHooks executes the poststart hook list in the host
// namespaces. If the list is nil or empty, it is a no-op.
//
// This corresponds to the OCI poststart lifecycle phase.
func (c *HookController) RunPoststartHooks(containerId string, hookList []spec.HookObject) error {
	if hookList == nil || len(hookList) == 0 {
		return nil
	}
	return c.runHookList(containerId, "poststart", hookList)
}

// RunPoststopHooks executes the poststop hook list in the host
// namespaces. If the list is nil or empty, it is a no-op.
//
// This corresponds to the OCI poststop lifecycle phase.
func (c *HookController) RunPoststopHooks(containerId string, hookList []spec.HookObject) error {
	if hookList == nil || len(hookList) == 0 {
		return nil
	}
	return c.runHookList(containerId, "poststop", hookList)
}

// runHookList executes a list of hooks in the host namespaces.
//
// For each hook:
//   - Validates that the hook path is non-empty
//   - Reads state.json and passes it as stdin to the hook
//   - Inherits the current environment and appends hook-specific variables
//   - Directs stdout and stderr to the runtime's stdio
//
// If any hook fails, execution stops and an error is returned which includes
// the phase name and index in the hook list.
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

// runHookListWithNsenter executes a list of hooks inside the container's
// namespaces using nsenter.
//
// For each hook:
//   - Validates that the hook path is non-empty
//   - Resolves the container init PID from state.json
//   - Builds an nsenter command targeting the container namespaces
//   - Passes state.json as stdin and appends hook environment variables
//   - Directs stdout and stderr to the runtime's stdio
//
// If any hook fails, execution stops and an error is returned which includes
// the phase name and index in the hook list.
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
