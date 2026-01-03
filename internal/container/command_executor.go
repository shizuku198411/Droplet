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

// SetStdin sets the standard input stream for the underlying command.
func (e *execCmd) SetStdin(r io.Reader) {
	e.cmd.Stdin = r
}

// SetSysProcAttr assigns the provided SysProcAttr to the underlying exec.Cmd.
func (e *execCmd) SetSysProcAttr(attr *syscall.SysProcAttr) {
	e.cmd.SysProcAttr = attr
}

// syscallHandler abstracts the operation of replacing the current
// process image with another program.
//
// It is defined as an interface to allow syscall.Exec to be mocked
// in tests and substituted by alternative implementations if needed.
type syscallHandler interface {
	Exec(argv0 string, argv []string, envv []string) error
}

// containerEnvPrepareSyscallHandler abstracts the set of syscalls used during
// container environment preparation inside the init process.
//
// This interface allows environment-setup logic (such as switching to the
// user-namespace root or configuring the UTS namespace hostname) to be tested
// without invoking real kernel syscalls. Production code typically provides a
// syscall-backed implementation, while unit tests may supply a mock or stub
// implementation to validate control flow and error handling.
type containerEnvPrepareSyscallHandler interface {
	Setresgid(rgid int, egid int, sgid int) error
	Setresuid(ruid int, euid int, suid int) error
	Sethostname(p []byte) error
}

// newSyscallHandler returns a kernelSyscall that delegates to
// syscall.Exec to replace the current process image.
func newSyscallHandler() *kernelSyscall {
	return &kernelSyscall{}
}

// kerneklSyscall is the default implementation of processReplacer.
//
// It invokes syscall.Exec directly, causing the current process to be
// replaced by the specified executable if successful.
type kernelSyscall struct{}

// Exec calls syscall.Exec with the provided arguments.
//
// On success, this call does not return. Any returned error indicates
// that the process could not be replaced.
func (k *kernelSyscall) Exec(argv0 string, argv []string, envv []string) error {
	return syscall.Exec(argv0, argv, envv)
}

// Setresgid changes the real, effective, and saved group IDs of the current
// process by invoking the kernel's setresgid(2) syscall.
//
// In the context of container initialization, this is typically used to switch
// the init process to GID 0 within the active user namespace before executing
// privileged setup operations.
func (k *kernelSyscall) Setresgid(rgid int, egid int, sgid int) error {
	return syscall.Setresgid(rgid, egid, sgid)
}

// Setresuid changes the real, effective, and saved user IDs of the current
// process by invoking the setresuid(2) syscall.
//
// During container environment preparation, this is commonly called immediately
// after Setresgid to complete the transition to UID 0 inside the user
// namespace so that privileged operations (e.g., mounts, hostname changes)
// may be performed safely.
func (k *kernelSyscall) Setresuid(ruid int, euid int, suid int) error {
	return syscall.Setresuid(ruid, euid, suid)
}

// Sethostname sets the hostname of the current UTS namespace by invoking the
// sethostname(2) syscall.
//
// Container init processes call this to assign a container-specific hostname
// derived from the OCI runtime specification or container identifier. Errors
// are returned if the operation is not permitted within the current namespace.
func (k *kernelSyscall) Sethostname(p []byte) error {
	return syscall.Sethostname(p)
}
