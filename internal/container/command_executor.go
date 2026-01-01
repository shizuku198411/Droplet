package container

import (
	"io"
	"os/exec"
	"syscall"
)

func newCommandFactory() *execCommandFactory {
	return &execCommandFactory{}
}

// commandFactory creates commandExecutor instances.
//
// The factory abstracts process creation so that callers do not depend
// directly on exec.Command. This makes the behavior testable by replacing
// the factory with a mock implementation.
type commandFactory interface {
	Command(name string, args ...string) commandExecutor
}

// execCommandFactory is the default implementation of commandFactory.
//
// It creates commandExecutor values backed by *exec.Cmd and launches
// real OS processes.
type execCommandFactory struct{}

// Command returns a commandExecutor that executes the given command
// using exec.Cmd.
func (e *execCommandFactory) Command(name string, args ...string) commandExecutor {
	return &execCmd{cmd: exec.Command(name, args...)}
}

// commandExecutor represents a process that can be started.
//
// It provides a minimal surface over exec.Cmd so that command execution
// can be substituted or mocked in tests.
type commandExecutor interface {
	Start() error
	Wait() error
	Pid() int
	SetStdout(w io.Writer)
	SetStderr(w io.Writer)
	SetStdin(r io.Reader)
	SetSysProcAttr(attr *syscall.SysProcAttr)
}

// execCmd is the concrete commandExecutor backed by exec.Cmd.
//
// It delegates all operations to the underlying exec.Cmd instance.
type execCmd struct {
	cmd *exec.Cmd
}

// Start starts the underlying process.
//
// It mirrors (*exec.Cmd).Start.
func (e *execCmd) Start() error {
	return e.cmd.Start()
}

func (e *execCmd) Wait() error {
	return e.cmd.Wait()
}

// Pid returns the PID of the started process.
//
// If the process has not been started, -1 is returned.
func (e *execCmd) Pid() int {
	if e.cmd.Process == nil {
		return -1
	}
	return e.cmd.Process.Pid
}

// SetStdout sets the stdout writer for the underlying command.
func (e *execCmd) SetStdout(w io.Writer) {
	e.cmd.Stdout = w
}

// SetStderr sets the stderr writer for the underlying command.
func (e *execCmd) SetStderr(w io.Writer) {
	e.cmd.Stderr = w
}

func (e *execCmd) SetStdin(r io.Reader) {
	e.cmd.Stdin = r
}

func (e *execCmd) SetSysProcAttr(attr *syscall.SysProcAttr) {
	e.cmd.SysProcAttr = attr
}
