package container

import (
	"droplet/internal/spec"
	"os"
)

// NewContainerInit returns a ContainerInit wired with the default
// implementations of its dependencies (fifoReader and processReplacer).
// This is the standard entry point for executing the container init phase.
func NewContainerInit() *ContainerInit {
	return &ContainerInit{
		fifoReader:           newContainerFifoHandler(),
		specLoader:           newFileSpecLoader(),
		containerEnvPreparer: newRootContainerEnvPrepare(),
		syscallHandler:       newSyscallHandler(),
	}
}

// newRootContainerEnvPreparer returns the default environment preparer
// implementation for containers started as root on the host.
//
// This preparer performs setup steps that assume the runtime is executing
// with full privileges (e.g., user-namespace root switching, hostname
// configuration). A separate implementation can be provided for rootless
// execution environments.
func newRootContainerEnvPrepare() *rootContainerEnvPreparer {
	return &rootContainerEnvPreparer{
		syscallHandler: newSyscallHandler(),
	}
}

// ContainerInit represents the runtime logic executed inside the
// container's init process.
//
// The init process waits for a start signal via FIFO and then
// replaces itself with the container entrypoint using execve-style
// semantics (syscall.Exec).
type ContainerInit struct {
	fifoReader           fifoReader
	specLoader           specLoader
	containerEnvPreparer containerEnvPreparer
	syscallHandler       syscallHandler
}

// Execute performs the init sequence for the container.
//
// The sequence is:
//
//  1. Wait for a start signal by reading from the FIFO path
//  2. Replace the current process image with the container entrypoint
//
// On success, this function does not return because the process image
// is replaced. Errors are returned only if the FIFO read fails or
// syscall.Exec cannot be invoked.
func (c *ContainerInit) Execute(opt InitOption) error {
	fifo := opt.Fifo
	entrypoint := opt.Entrypoint

	// 1. read fifo for waiting start signal
	if err := c.fifoReader.readFifo(fifo); err != nil {
		return err
	}

	// 2. load config.json
	spec, err := c.specLoader.loadFile(opt.ContainerId)
	if err != nil {
		return err
	}

	// 3. prepare container environment
	if err := c.containerEnvPreparer.prepare(spec); err != nil {
		return err
	}

	// 4. replace process image with the container entrypoint
	if err := c.syscallHandler.Exec(entrypoint[0], entrypoint, os.Environ()); err != nil {
		return err
	}

	return nil
}

// containerEnvPreparer defines the behavior for preparing the container
// environment inside the init process.
//
// Implementations of this interface are responsible for performing
// container-local setup steps such as user namespace UID/GID switching,
// hostname configuration, filesystem setup, and other initialization logic
// that must occur before the container entrypoint is executed.
type containerEnvPreparer interface {
	prepare(spec spec.Spec) error
}

// rootContainerEnvPreparer is the default envPreparer implementation used
// for privileged (root-executed) containers.
//
// It performs environment initialization tasks inside the init process,
// such as switching to UID/GID 0 within the user namespace and configuring
// the UTS namespace hostname. Additional setup steps (mounts, pivot_root,
// capability adjustments, etc.) may be added to this implementation as
// container initialization evolves.
type rootContainerEnvPreparer struct {
	syscallHandler containerEnvPrepareSyscallHandler
}

// prepare performs the container-environment initialization sequence based
// on the provided OCI runtime specification.
//
// The current sequence consists of:
//
//  1. Switching to UID/GID 0 within the user namespace
//  2. Setting the container hostname in the UTS namespace
//
// Additional lifecycle steps (filesystem mounts, pivot_root, /proc setup,
// etc.) can be appended to this method as required.
func (p *rootContainerEnvPreparer) prepare(spec spec.Spec) error {
	// 1. change uid=0(root) inside container
	if err := p.switchToUserNamespaceRoot(); err != nil {
		return err
	}

	// 2. set hostname
	if err := p.setHostnameToContainerId(spec.Hostname); err != nil {
		return err
	}

	return nil
}

// switchToUserNamespaceRoot switches the current process credentials to
// UID and GID 0 within the active user namespace.
//
// This ensures that subsequent privileged operations (such as mount,
// pivot_root, or hostname changes) execute with the required namespace-
// scoped capabilities, even when the process was not initially running as
// namespace-root.
func (p *rootContainerEnvPreparer) switchToUserNamespaceRoot() error {
	// switch root group (gid=0)
	if err := p.syscallHandler.Setresgid(0, 0, 0); err != nil {
		return err
	}
	// switch root user (uid=0)
	if err := p.syscallHandler.Setresuid(0, 0, 0); err != nil {
		return err
	}
	return nil
}

// setHostnameToContainerId configures the hostname for the process inside
// the UTS namespace.
//
// The hostname value is typically derived from the container ID or the
// OCI spec. An error is returned if the syscall fails or the namespace
// does not permit hostname updates.
func (p *rootContainerEnvPreparer) setHostnameToContainerId(hostname string) error {
	if err := p.syscallHandler.Sethostname([]byte(hostname)); err != nil {
		return err
	}
	return nil
}
