package container

import (
	"fmt"
	"os"

	"droplet/internal/hook"
	"droplet/internal/spec"
	"droplet/internal/status"
	"droplet/internal/utils"
)

// NewContainerCreator constructs a ContainerCreator with the default
// implementations of its dependencies (SpecLoader, FifoCreator,
// ProcessExecutor, network/cgroup preparers, status manager, hook controller).
// This acts as the main entry point for the container creation workflow.
func NewContainerCreator() *ContainerCreator {
	return &ContainerCreator{
		specLoader:               newFileSpecLoader(),
		fifoCreator:              newContainerFifoHandler(),
		processExecutor:          newContainerInitExecutor(),
		containerNetworkPreparer: newContainerNetworkController(),
		containerCgroupPreparer:  newContainerCgroupController(),
		containerStatusManager:   status.NewStatusHandler(),
		containerHookController:  hook.NewHookController(),
	}
}

// newContainerInitExecutor constructs a containerInitExecutor with the default
// command factory. It is used as the default implementation for spawning
// the container init process workflow.
func newContainerInitExecutor() *containerInitExecutor {
	return &containerInitExecutor{
		commandFactory: &utils.ExecCommandFactory{},
	}
}

// ContainerCreator orchestrates the container creation flow for a single
// container instance.
//
// The flow currently consists of:
//
//  1. Loading the OCI spec (config.json)
//  2. Creating the initial state.json (status=creating, pid=0)
//  3. Running createRuntime hooks
//  4. Creating the FIFO used for init synchronization
//  5. Launching the init process via the init subcommand
//  6. Configuring cgroups for the init process
//  7. Configuring network for the init process
//  8. Updating state.json (status=created, pid=init pid)
//  9. Running createContainer hooks
//
// Each step is delegated to an interface to allow testing and substitution.
type ContainerCreator struct {
	specLoader               specLoader
	fifoCreator              fifoCreator
	processExecutor          processExecutor
	containerNetworkPreparer containerNetworkPreparer
	containerCgroupPreparer  containerCgroupPreparer
	containerStatusManager   status.ContainerStatusManager
	containerHookController  hook.ContainerHookController
}

// Create executes the container creation pipeline for the given container ID.
//
// It coordinates the high-level workflow by:
//   - Loading the spec
//   - Initializing container state
//   - Running lifecycle hooks
//   - Spawning the init process
//   - Applying cgroup and network configuration
//   - Updating final status
//
// This method performs no low-level work itself and relies entirely on
// its collaborators. If any step fails, the error is returned immediately.
func (c *ContainerCreator) Create(opt CreateOption) error {
	// 1. load config.json
	spec, err := c.specLoader.loadFile(opt.ContainerId)
	if err != nil {
		return err
	}

	// 2. create state.json
	//      status = creating
	//      pid = 0
	if err := c.containerStatusManager.CreateStatusFile(
		opt.ContainerId,
		0,
		status.CREATING,
		spec.Root.Path,
		utils.ContainerDir(opt.ContainerId),
		spec.Annotations,
	); err != nil {
		return err
	}

	// 3. HOOK: createRuntime
	if err := c.containerHookController.RunCreateRuntimeHooks(
		opt.ContainerId,
		spec.Hooks.CreateRuntime,
	); err != nil {
		return err
	}

	// 4. create fifo
	fifo := utils.FifoPath(opt.ContainerId)
	if err := c.fifoCreator.createFifo(fifo); err != nil {
		return err
	}

	// 5. execute init subcommand
	initPid, err := c.processExecutor.executeInit(opt.ContainerId, spec, fifo)
	if err != nil {
		return err
	}

	// output when init process has been created
	// if --print-pid is setted, print message with pid
	// otherwise print message with Container ID
	if opt.PrintPidFlag {
		fmt.Printf("create container success. pid: %d\n", initPid)
	} else {
		fmt.Printf("create container success. ID: %s\n", opt.ContainerId)
	}

	// 6. cgroup setup
	if err := c.containerCgroupPreparer.prepare(opt.ContainerId, spec, initPid); err != nil {
		return err
	}

	// 7. network setup
	if err := c.containerNetworkPreparer.prepare(opt.ContainerId, initPid, spec.Annotations); err != nil {
		return err
	}

	// 8. update state.json
	//      status = created
	//      pid    = init pid
	if err := c.containerStatusManager.UpdateStatus(
		opt.ContainerId,
		status.CREATED,
		initPid,
	); err != nil {
		return err
	}

	// 9. HOOK: createContainer
	if err := c.containerHookController.RunCreateContainerHooks(
		opt.ContainerId,
		spec.Hooks.CreateContainer,
	); err != nil {
		return err
	}

	return nil
}

// processExecutor defines the behavior for spawning the container init process.
//
// It is an interface so that the behavior can be mocked in tests and
// replaced by alternative implementations if needed.
type processExecutor interface {
	executeInit(containerId string, spec spec.Spec, fifo string) (int, error)
}

// containerInitExecutor is the default implementation of processExecutor.
//
// It invokes this binary with the `init` subcommand and the FIFO path,
// passing the spec's process args as the container entrypoint.
type containerInitExecutor struct {
	commandFactory utils.CommandFactory
}

// executeInit starts the init process and returns its PID.
//
// The init process is started as a child of the current runtime binary
// with the appropriate namespace and process attributes applied. The FIFO
// path is passed as an argument so that the init process can synchronize
// with the runtime before proceeding.
func (c *containerInitExecutor) executeInit(containerId string, spec spec.Spec, fifo string) (int, error) {
	// retrieve entrypoint from spec
	entrypoint := spec.Process.Args

	// prepare init subcommand
	initArgs := append([]string{"init", containerId, fifo}, entrypoint...)
	cmd := c.commandFactory.Command(os.Args[0], initArgs...)
	// TODO: set stdout/stderr to log files

	// apply SysProcAttr
	nsConfig := buildNamespaceConfig(spec)
	procAttr := buildProcAttrForRootContainer(nsConfig)
	sysProcAttr := buildSysProcAttr(procAttr)
	cmd.SetSysProcAttr(sysProcAttr)

	// execute init subcommand
	if err := cmd.Start(); err != nil {
		return -1, err
	}

	return cmd.Pid(), nil
}
