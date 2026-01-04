package container

import (
	"io"
	"os"
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
	Mount(source string, target string, fstype string, flags uintptr, data string) error
	Unmount(target string, flags int) error
	PivotRoot(newroot string, putold string) error
	Chdir(path string) error
	Mkdir(path string, mode uint32) error
	MkdirAll(path string, perm os.FileMode) error
	Rmdir(path string) error
	Stat(name string) (os.FileInfo, error)
	Create(name string) (*os.File, error)
	IsNotExist(err error) bool
	Symlink(oldname string, newname string) error
	Lstat(name string) (os.FileInfo, error)
	OpenFile(name string, flag int, perm os.FileMode) (*os.File, error)
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

// Mount performs a mount(2) system call.
//
// It is a thin wrapper around syscall.Mount and is provided so that
// filesystem operations can be abstracted and mocked during testing.
func (k *kernelSyscall) Mount(source string, target string, fstype string, flags uintptr, data string) error {
	return syscall.Mount(source, target, fstype, flags, data)
}

// Unmount performs an umount(2) / umount2(2) system call.
//
// target specifies the mountpoint to unmount and flags corresponds to
// MNT_* constants such as MNT_DETACH.
func (k *kernelSyscall) Unmount(target string, flags int) error {
	return syscall.Unmount(target, flags)
}

// PivotRoot performs a pivot_root(2) system call.
//
// newroot becomes the new root filesystem ("/") and putold becomes the
// location where the previous root filesystem is mounted temporarily.
// This is typically used during container startup after preparing a
// private mount namespace.
func (k *kernelSyscall) PivotRoot(newroot string, putold string) error {
	return syscall.PivotRoot(newroot, putold)
}

// Chdir changes the current working directory of the calling process.
//
// It wraps syscall.Chdir to allow testing and syscall abstraction.
func (k *kernelSyscall) Chdir(path string) error {
	return syscall.Chdir(path)
}

// Mkdir creates a single directory using the mkdir(2) syscall.
//
// Unlike MkdirAll, this does not create parent directories.
// The mode parameter corresponds to POSIX permission bits (e.g., 0700).
func (k *kernelSyscall) Mkdir(path string, mode uint32) error {
	return syscall.Mkdir(path, mode)
}

// MkdirAll creates a directory and all missing parent directories.
//
// This wraps os.MkdirAll rather than a raw syscall because recursive
// directory creation is a library-level operation rather than a single
// kernel call.
func (k *kernelSyscall) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

// Rmdir removes an empty directory using the rmdir(2) syscall.
//
// If the directory is not empty, the call fails with an error.
func (k *kernelSyscall) Rmdir(path string) error {
	return syscall.Rmdir(path)
}

// Stat returns file or directory metadata using the stat(2) syscall.
//
// The returned FileInfo describes the file. If the file does not exist,
// an error that can be checked with IsNotExist is returned.
func (k *kernelSyscall) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

// Create creates a new file or truncates an existing one.
//
// This is a thin wrapper around os.Create, which internally issues
// an open(2) call with O_CREAT|O_TRUNC.
func (k *kernelSyscall) Create(name string) (*os.File, error) {
	return os.Create(name)
}

// IsNotExist reports whether an error indicates that a file or directory
// does not exist.
//
// This wraps os.IsNotExist so that callers do not depend directly on
// os package helpers.
func (k *kernelSyscall) IsNotExist(err error) bool {
	return os.IsNotExist(err)
}

// Symlink creates a symbolic link using the symlink(2) syscall.
//
// oldname is the target of the link and newname is the created symlink path.
func (k *kernelSyscall) Symlink(oldname string, newname string) error {
	return os.Symlink(oldname, newname)
}

// Lstat retrieves metadata for the named file without following symlinks.
//
// It is equivalent to the lstat(2) syscall and differs from Stat,
// which follows symbolic links.
func (k *kernelSyscall) Lstat(name string) (os.FileInfo, error) {
	return os.Lstat(name)
}

// OpenFile opens a file with the given flags and permissions.
//
// This corresponds to the open(2) syscall and supports flags such as
// O_RDONLY, O_WRONLY, O_APPEND, O_CREAT, and O_EXCL.
func (k *kernelSyscall) OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(name, flag, perm)
}
