package container

import (
	"fmt"
	"os"

	"droplet/internal/spec"
)

// NewContainerCreator constructs a ContainerCreator with the default
// implementations of its dependencies (SpecLoader, FifoCreator, ProcessExecutor).
// This acts as the main entry point for the container creation workflow.
func NewContainerCreator() *ContainerCreator {
	return &ContainerCreator{
		specLoader:      newFileSpecLoader(),
		fifoCreator:     newContainerFifoHandler(),
		processExecutor: newContainerInitExecutor(),
	}
}

// newContainerInitExecutor constructs a containerInitExecutor with the default implementations.
// This acts as the main entry point for spawning the container init process workflow.
func newContainerInitExecutor() *containerInitExecutor {
	return &containerInitExecutor{
		commandFactory: &execCommandFactory{},
	}
}

// ContainerCreator orchestrates the container creation flow.
//
// The flow currently consists of:
//
//  1. Loading the OCI spec (config.json)
//  2. Creating the FIFO used for init synchronization
//  3. Launching the init process via the init subcommand
//
// Each step is delegated to an interface to allow testing and substitution.
type ContainerCreator struct {
	specLoader      specLoader
	fifoCreator     fifoCreator
	processExecutor processExecutor
}

// Create executes the container creation pipeline for the given container ID.
// This method performs no low-level work itself — it coordinates collaborators.
func (c *ContainerCreator) Create(opt CreateOption) error {
	// load config.json
	spec, err := c.specLoader.loadFile(opt.ContainerId)
	if err != nil {
		return err
	}

	// create fifo
	fifo := fifoPath(opt.ContainerId)
	if err := c.fifoCreator.createFifo(fifo); err != nil {
		return err
	}

	// execute init subcommand
	initPid, err := c.processExecutor.executeInit(opt.ContainerId, spec, fifo)
	if err != nil {
		return err
	}

	fmt.Printf("init process has been created. pid: %d\n", initPid)

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
	commandFactory commandFactory
}

// executeInit starts the init process and returns its PID.
//
// The init process is started as a child of the current runtime binary.
// The FIFO path is passed as an argument so that the init process can
// synchronize with the runtime.
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
