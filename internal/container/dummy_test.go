package container

import (
	"io"
	"syscall"

	"droplet/internal/spec"
)

// dummy spec loader
type dummyFileSpecLoader struct {
	containerId string
	spec        spec.Spec
	err         error
}

func (d *dummyFileSpecLoader) loadFile(containerId string) (spec.Spec, error) {
	d.containerId = containerId
	return d.spec, d.err
}

// dummy fifo handler
type dummyFifoHandler struct {
	path      string
	createErr error
	removeErr error
	writeErr  error
	readErr   error
	callFlag  bool
}

func (d *dummyFifoHandler) createFifo(path string) error {
	d.callFlag = true
	d.path = path
	return d.createErr
}

func (d *dummyFifoHandler) removeFifo(path string) error {
	d.callFlag = true
	d.path = path
	return d.removeErr
}

func (d *dummyFifoHandler) readFifo(path string) error {
	d.callFlag = true
	d.path = path
	return d.readErr
}

func (d *dummyFifoHandler) writeFifo(path string) error {
	d.callFlag = true
	d.path = path
	return d.writeErr
}

// dummy container init executor
type dummyContainerInitExecutor struct {
	Spec spec.Spec
	Fifo string
	Pid  int
	Err  error
}

func (d *dummyContainerInitExecutor) executeInit(spec spec.Spec, fifo string) (int, error) {
	d.Spec = spec
	d.Fifo = fifo
	return d.Pid, d.Err
}

type dummyCmd struct {
	name   string
	args   []string
	stdout io.Writer
	stderr io.Writer
	stdin  io.Reader
	pid    int
	err    error
	attr   *syscall.SysProcAttr

	startFlag bool
}

func (d *dummyCmd) Start() error {
	d.startFlag = true
	return d.err
}

func (d *dummyCmd) Wait() error {
	d.startFlag = true
	return d.err
}

func (d *dummyCmd) Pid() int {
	return d.pid
}

func (d *dummyCmd) SetStdout(w io.Writer) {
	d.stdout = w
}

func (d *dummyCmd) SetStderr(w io.Writer) {
	d.stderr = w
}

func (d *dummyCmd) SetStdin(r io.Reader) {
	d.stdin = r
}

func (d *dummyCmd) SetSysProcAttr(attr *syscall.SysProcAttr) {
	d.attr = attr
}

type dummyCommandFactory struct {
	commandName string
	commandArgs []string
	cmd         *dummyCmd
}

func (d *dummyCommandFactory) Command(name string, args ...string) commandExecutor {
	d.commandName = name
	d.commandArgs = args
	return d.cmd
}

func dummySpec() spec.Spec {
	return spec.Spec{
		Process: spec.ProcessObject{
			Args: []string{
				"/bin/sh",
			},
		},
	}
}

// dummy process replacer
type dummySyscallProcessReplacer struct {
	argv0    string
	argv     []string
	envv     []string
	err      error
	callFlag bool
}

func (d *dummySyscallProcessReplacer) Exec(argv0 string, argv []string, envv []string) error {
	d.callFlag = true
	d.argv0 = argv0
	d.argv = argv
	d.envv = envv
	return d.err
}
