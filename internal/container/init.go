package container

import (
	"droplet/internal/spec"
	"droplet/internal/utils"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"syscall"

	"github.com/syndtr/gocapability/capability"
)

// NewContainerInit returns a ContainerInit wired with the default
// implementations of its dependencies (fifoReader and processReplacer).
// This is the standard entry point for executing the container init phase.
func NewContainerInit() *ContainerInit {
	return &ContainerInit{
		fifoReader:           newContainerFifoHandler(),
		specLoader:           newFileSpecLoader(),
		containerEnvPreparer: newRootContainerEnvPrepare(),
		syscallHandler:       utils.NewSyscallHandler(),
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
		syscallHandler: utils.NewSyscallHandler(),
		seccompHandler: NewSeccompManager(),
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
	syscallHandler       utils.SyscallHandler
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
	if err := c.containerEnvPreparer.prepare(opt.ContainerId, spec); err != nil {
		return err
	}

	// 4. replace process image with the container entrypoint
	if err := c.syscallHandler.Exec(entrypoint[0], entrypoint, slices.Concat(os.Environ(), spec.Process.Env)); err != nil {
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
	prepare(containerId string, spec spec.Spec) error
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
	syscallHandler utils.KernelSyscallHandler
	seccompHandler SeccompHandler
}

// prepare sets up the runtime environment for the root container process
// according to the provided OCI spec.
//
// The workflow is:
//  1. Switch to uid=0 (root) inside the user namespace
//  2. Set the hostname to the container ID from the spec
//  3. Set up the overlay filesystem based on rootfs and image annotations
//  4. Mount the configured filesystems
//  5. Mount standard device files under the new root
//  6. Create required symbolic links under the new root
//  7. Perform pivot_root into the container root filesystem
//  8. Configure Linux capabilities for the process
//
// If any step fails, the error is returned immediately and the remaining
// steps are not executed.
func (p *rootContainerEnvPreparer) prepare(containerId string, spec spec.Spec) error {
	// 1. change uid=0(root) inside container
	if err := p.switchToUserNamespaceRoot(); err != nil {
		return err
	}
	// 2. set hostname
	if err := p.setHostnameToContainerId(spec.Hostname); err != nil {
		return err
	}
	// 3. set env
	if err := p.setEnv(spec.Process.Env); err != nil {
		return err
	}
	// 4. setup overlay
	if err := p.setupOverlay(spec.Root.Path, spec.Annotations.Image); err != nil {
		return err
	}
	// 5. mount filesystem
	if err := p.mountFilesystem(containerId, spec.Root.Path, spec.Mounts); err != nil {
		return err
	}
	// 6. mount standard device
	if err := p.mountStdDevice(spec.Root.Path); err != nil {
		return err
	}
	// 7. create symbolic link
	if err := p.createSymbolicLink(spec.Root.Path); err != nil {
		return err
	}
	// 8. pivot_root
	if err := p.pivotRoot(spec.Root.Path); err != nil {
		return err
	}
	// 9. set capability
	if err := p.setCapability(spec.Process.Capabilities); err != nil {
		return err
	}
	// 10. install seccomp (NO_NEW_PRIVS + filter)
	if err := p.seccompHandler.InstallDenyFilter(*spec.LinuxSpec.Seccomp); err != nil {
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

func (p *rootContainerEnvPreparer) setEnv(envlist []string) error {
	for _, e := range envlist {
		envParts := strings.Split(e, "=")
		k, v := envParts[0], envParts[1]
		if err := p.syscallHandler.Setenv(k, v); err != nil {
			return err
		}
	}
	return nil
}

// setupOverlay mounts the container root filesystem using overlayfs.
//
// imageAnnotation is a JSON string that is decoded into ImageConfigObject,
// which contains lower (image layers), upper, and work directories.
// The overlay filesystem is mounted at the given rootfs path.
func (p *rootContainerEnvPreparer) setupOverlay(rootfs string, imageAnnotation string) error {
	// convert string to json
	var imageConfig spec.ImageConfigObject
	if err := utils.StringToJson(imageAnnotation, &imageConfig); err != nil {
		return err
	}

	// mount parameter
	mountSource := "overlay"
	mountTarget := rootfs
	mountFstype := imageConfig.RootfsType
	mountFlags := uintptr(0)
	// mount data contains following parameter
	// - lowerdir : container image layers
	// - upperdir : directory for storing differences with lowerdir
	// - workdir  : directory for working directory
	lowerDir := strings.Join(imageConfig.ImageLayer, ":")
	upperDir := imageConfig.UpperDir
	workDir := imageConfig.WorkDir
	mountData := fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s", lowerDir, upperDir, workDir)

	// overlay
	if err := p.syscallHandler.Mount(mountSource, mountTarget, mountFstype, mountFlags, mountData); err != nil {
		return err
	}

	return nil
}

// mountFilesystem mounts all filesystems required for the container runtime
// as well as user-specified bind mounts.
//
// The mountList contains entries such as /proc, /dev, /sys, cgroup, tmpfs,
// and arbitrary host paths. For bind mounts, this method prepares the
// destination path depending on whether the source is a file or a directory.
func (p *rootContainerEnvPreparer) mountFilesystem(containerId string, rootfs string, mountList []spec.MountObject) error {
	// mount path validation for protecting path traversal
	securePath := func(rootfs, dest string) (string, error) {
		rel := strings.TrimPrefix(dest, "/")
		clean := filepath.Clean(rel)

		if clean == "." || strings.HasPrefix(clean, ".."+string(os.PathSeparator)) || clean == ".." {
			return "", fmt.Errorf("invalid destination: %q", dest)
		}

		fullPath := filepath.Join(rootfs, clean)
		root := filepath.Clean(rootfs) + string(os.PathSeparator)
		if !strings.HasPrefix(fullPath+string(os.PathSeparator), root) && fullPath != filepath.Clean(rootfs) {
			return "", fmt.Errorf("destination escapes rootfs: %q -> %q", dest, fullPath)
		}
		return fullPath, nil
	}

	// source mount validation
	// the following source is denied by default
	//   /proc, /sys, /dev, /run, /var/run, /boot, /root, /
	hasDeniedSource := func(source string) bool {
		p := filepath.Clean(source)
		if p == "/" {
			return true
		}
		deniedList := []string{"/proc", "/sys", "/dev", "/run", "/var/run", "/boot", "/root"}
		for _, d := range deniedList {
			if p == d || strings.HasPrefix(p, d+string(os.PathSeparator)) {
				return true
			}
		}
		return false
	}

	// mount file system required for operation. required fs is the following
	//   /proc, /dev, /dev/pts, /sys, /sys/fs/cgroup, /dev/mqueue, /dev/shm
	// additionally, mount user-specified host directories
	var prerequiredMounts = []spec.MountObject{
		{
			Destination: "/proc",
			Type:        "proc",
			Source:      "proc",
			Options: []string{
				"nosuid",
				"noexec",
				"nodev",
			},
		},
		{
			Destination: "/dev",
			Type:        "tmpfs",
			Source:      "tmpfs",
			Options: []string{
				"nosuid",
				"strictatime",
				"mode=755",
				"size=65536k",
			},
		},
		{
			Destination: "/dev/pts",
			Type:        "devpts",
			Source:      "devpts",
			Options: []string{
				"nosuid",
				"noexec",
				"newinstance",
				"ptmxmode=0666",
				"mode=0620",
				"gid=5",
			},
		},
		{
			Destination: "/sys",
			Type:        "sysfs",
			Source:      "sysfs",
			Options: []string{
				"nosuid",
				"noexec",
				"nodev",
				"ro",
			},
		},
		{
			Destination: "/tmp",
			Type:        "tmpfs",
			Source:      "tmpfs",
			Options: []string{
				"nosuid",
				"nodev",
				"mode=1777",
				"size=65536k",
			},
		},
		{
			Destination: "/run",
			Type:        "tmpfs",
			Source:      "tmpfs",
			Options: []string{
				"nosuid",
				"nodev",
				"mode=755",
				"size=65536k",
			},
		},
		{
			Destination: "/proc/sys",
			Type:        "tmpfs",
			Source:      "tmpfs",
			Options: []string{
				"nosuid",
				"noexec",
				"nodev",
				"mode=0555",
				"size=0",
			},
		},
		{
			Destination: "/proc/sysrq-trigger",
			Type:        "bind",
			Source:      "/dev/null",
			Options: []string{
				"rbind",
				"ro",
			},
		},
		{
			Destination: "/sys/firmware",
			Type:        "tmpfs",
			Source:      "tmpfs",
			Options: []string{
				"nosuid",
				"noexec",
				"nodev",
				"mode=0555",
				"size=0",
			},
		},
		{
			Destination: "/sys/fs/bpf",
			Type:        "tmpfs",
			Source:      "tmpfs",
			Options: []string{
				"nosuid",
				"noexec",
				"nodev",
				"ro",
				"mode=0555",
				"size=0",
			},
		},
		{
			Destination: "/sys/fs/cgroup",
			Type:        "cgroup2",
			Source:      "cgroup",
			Options: []string{
				"nosuid",
				"nodev",
				"noexec",
				"ro",
			},
		},
		{
			Destination: "/dev/mqueue",
			Type:        "mqueue",
			Source:      "mqueue",
			Options: []string{
				"nosuid",
				"noexec",
				"nodev",
			},
		},
		{
			Destination: "/dev/shm",
			Type:        "tmpfs",
			Source:      "shm",
			Options: []string{
				"nosuid",
				"noexec",
				"nodev",
				"mode=1777",
				"size=67108864",
			},
		},
		{
			Destination: "/etc/resolv.conf",
			Type:        "bind",
			Source:      fmt.Sprintf("/etc/raind/container/%s/etc/resolv.conf", containerId),
			Options: []string{
				"rbind",
				"rprivate",
			},
		},
		{
			Destination: "/etc/hostname",
			Type:        "bind",
			Source:      fmt.Sprintf("/etc/raind/container/%s/etc/hostname", containerId),
			Options: []string{
				"rbind",
				"rprivate",
			},
		},
		{
			Destination: "/etc/hosts",
			Type:        "bind",
			Source:      fmt.Sprintf("/etc/raind/container/%s/etc/hosts", containerId),
			Options: []string{
				"rbind",
				"rprivate",
			},
		},
	}

	// user mounts
	for _, user_mount := range mountList {
		// validate
		if hasDeniedSource(user_mount.Source) {
			return fmt.Errorf("invalid mount source: %s", user_mount.Source)
		}
		prerequiredMounts = append(prerequiredMounts, spec.MountObject{
			Destination: user_mount.Destination,
			Type:        user_mount.Type,
			Source:      user_mount.Source,
			Options:     user_mount.Options,
		})
	}

	for _, mountConfig := range prerequiredMounts {
		var (
			mountFlags uintptr
			mountData  string
			dataStrTmp []string
		)
		if mountConfig.Options != nil {
			for _, option := range mountConfig.Options {
				switch option {
				case "nosuid":
					mountFlags |= syscall.MS_NOSUID
				case "noexec":
					mountFlags |= syscall.MS_NOEXEC
				case "nodev":
					mountFlags |= syscall.MS_NODEV
				case "ro":
					mountFlags |= syscall.MS_RDONLY
				case "rw":
					// ignore
				case "bind":
					mountFlags |= syscall.MS_BIND
				case "strictatime":
					mountFlags |= syscall.MS_STRICTATIME
				case "noatime":
					mountFlags |= syscall.MS_NOATIME
				case "relatime":
					mountFlags |= syscall.MS_RELATIME
				case "rbind":
					mountFlags |= syscall.MS_BIND | syscall.MS_REC
				case "rprivate":
					// ignore
				default:
					dataStrTmp = append(dataStrTmp, option)
				}
			}
			mountData = strings.Join(dataStrTmp, ",")
		} else {
			mountFlags = uintptr(0)
		}

		mountPath, err := securePath(rootfs, mountConfig.Destination)
		if err != nil {
			return err
		}

		// the process differs depending on whether the source to be mounted is a directory or a file.
		// if the source is a directory, the destination directory is checked for existence and created if it does not exist.
		// if the source is a file, the parent directory is created and an empty file is created.
		// this process is only bind mount
		if mountConfig.Type == "bind" {
			// retrieve source info
			srcInfo, statErr := p.syscallHandler.Stat(mountConfig.Source)
			if statErr != nil {
				return statErr
			}

			if srcInfo.IsDir() { // source: directory
				// check if target directory is exists
				if _, err := p.syscallHandler.Stat(mountPath); p.syscallHandler.IsNotExist(err) {
					if err := p.syscallHandler.MkdirAll(mountPath, os.ModePerm); err != nil {
						return err
					}
				}
			} else { // source: file
				// create parent directory if not exists
				if err := p.syscallHandler.MkdirAll(filepath.Dir(mountPath), os.ModePerm); err != nil {
					return err
				}
				if _, err := p.syscallHandler.Stat(mountPath); p.syscallHandler.IsNotExist(err) {
					f, err := p.syscallHandler.OpenFile(mountPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
					if err != nil {
						return err
					}
					f.Close()
				}
			}
		} else {
			// check if target directory is exists
			if _, err := p.syscallHandler.Stat(mountPath); p.syscallHandler.IsNotExist(err) {
				if err := p.syscallHandler.MkdirAll(mountPath, os.ModePerm); err != nil {
					return err
				}
			}
		}

		// mount
		if err := p.syscallHandler.Mount(
			mountConfig.Source,
			mountPath,
			mountConfig.Type,
			mountFlags,
			mountData,
		); err != nil {
			return err
		}

		if err := p.syscallHandler.Mount(
			"",
			mountPath,
			"",
			syscall.MS_PRIVATE|syscall.MS_REC,
			"",
		); err != nil {
			return err
		}
	}

	return nil
}

// mountStdDevice bind-mounts standard device files into the container's /dev.
//
// The following devices are mounted from the host:
//   - /dev/random
//   - /dev/urandom
//   - /dev/null
//   - /dev/full
//   - /dev/zero
//   - /dev/tty
//
// If the destination file does not exist under rootfs, it is created first.
func (p *rootContainerEnvPreparer) mountStdDevice(rootfs string) error {
	devices := []string{
		"random",
		"urandom",
		"null",
		"zero",
		"full",
		"tty",
	}
	for _, device := range devices {
		destination := filepath.Join(rootfs, "dev", device)
		// check if the file exist
		if _, err := p.syscallHandler.Stat(destination); p.syscallHandler.IsNotExist(err) {
			// create
			if _, err := p.syscallHandler.Create(destination); err != nil {
				return err
			}
		}
		// mount
		if err := p.syscallHandler.Mount(
			"/dev/"+device,
			destination,
			"",
			syscall.MS_BIND,
			"",
		); err != nil {
			return err
		}
		// remount for setting read-only flag
		if err := p.syscallHandler.Mount(
			"",
			destination,
			"",
			syscall.MS_BIND|syscall.MS_REMOUNT|syscall.MS_RDONLY|syscall.MS_NOEXEC|syscall.MS_NOSUID,
			"",
		); err != nil {
			return err
		}
	}
	return nil
}

// createSymbolicLink creates standard device-related symlinks under /dev
// inside the container rootfs.
//
// The following symlinks are created if they do not already exist:
//   - /dev/fd     -> /proc/self/fd
//   - /dev/stdin  -> /proc/self/fd/0
//   - /dev/stdout -> /proc/self/fd/1
//   - /dev/stderr -> /proc/self/fd/2
//   - /dev/ptmx   -> /dev/pts/ptmx
func (p *rootContainerEnvPreparer) createSymbolicLink(rootfs string) error {
	deviceDir := filepath.Join(rootfs, "dev")
	symlinks := []struct {
		link   string
		target string
	}{
		{filepath.Join(deviceDir, "fd"), "/proc/self/fd"},
		{filepath.Join(deviceDir, "stdin"), "/proc/self/fd/0"},
		{filepath.Join(deviceDir, "stdout"), "/proc/self/fd/1"},
		{filepath.Join(deviceDir, "stderr"), "/proc/self/fd/2"},
		{filepath.Join(deviceDir, "ptmx"), "/dev/pts/ptmx"},
	}

	for _, symlink := range symlinks {
		if _, err := p.syscallHandler.Lstat(symlink.link); err == nil {
			continue
		}
		if err := p.syscallHandler.Symlink(symlink.target, symlink.link); err != nil {
			return err
		}
	}

	return nil
}

// pivotRoot performs a pivot_root into the given rootfs and cleans up the old root.
//
// The sequence is:
//  1. create a put_old directory under the new root
//  2. call pivot_root(new_root, put_old)
//  3. chdir to "/"
//  4. unmount the old root at /put_old with MNT_DETACH
//  5. remove the /put_old directory
func (p *rootContainerEnvPreparer) pivotRoot(rootfs string) error {
	// oldroot directory
	putoldDir := filepath.Join(rootfs, "put_old")

	// 1. create put_old directory
	if err := p.syscallHandler.Mkdir(putoldDir, 0700); err != nil {
		return err
	}
	// 2. pivot_root
	if err := p.syscallHandler.PivotRoot(rootfs, putoldDir); err != nil {
		return err
	}
	// 3. change directory to root
	if err := p.syscallHandler.Chdir("/"); err != nil {
		return err
	}
	// 4. unmount put_old
	if err := p.syscallHandler.Unmount("/put_old", syscall.MNT_DETACH); err != nil {
		return err
	}
	// 5. remove put_old
	if err := p.syscallHandler.Rmdir("/put_old"); err != nil {
		return err
	}

	return nil
}

// setCapability configures Linux capabilities for the current (init) process
// according to the provided OCI capability configuration.
//
// The workflow is:
//  1. Create a capability set for PID 0 (the calling process)
//  2. Clear all capability sets (BOUNDING, PERMITTED, INHERITABLE, EFFECTIVE, AMBIENT)
//  3. Convert capability names from the spec to capability.Cap values
//  4. Populate each capability set from the corresponding field in capConfig
//  5. Apply the updated capability sets to the process
//
// If capability initialization or application fails, an error is returned.
func (p *rootContainerEnvPreparer) setCapability(capConfig spec.CapabilityObject) error {
	// set current process(init process) capability
	c, err := capability.NewPid2(0)
	if err != nil {
		return err
	}

	// clear all cap
	c.Clear(capability.BOUNDING | capability.PERMITTED | capability.INHERITABLE | capability.EFFECTIVE | capability.AMBIENT)

	// set bounding
	if len(capConfig.Bounding) > 0 {
		c.Set(capability.BOUNDING, toCaps(capConfig.Bounding)...)
	}
	// set permitted
	if len(capConfig.Permitted) > 0 {
		c.Set(capability.PERMITTED, toCaps(capConfig.Permitted)...)
	}
	// set inheritable
	if len(capConfig.Inheritable) > 0 {
		c.Set(capability.INHERITABLE, toCaps(capConfig.Inheritable)...)
	}
	// set effective
	if len(capConfig.Effective) > 0 {
		c.Set(capability.EFFECTIVE, toCaps(capConfig.Effective)...)
	}
	// set ambient
	if len(capConfig.Ambient) > 0 {
		c.Set(capability.AMBIENT, toCaps(capConfig.Ambient)...)
	}

	// apply
	if err := c.Apply(capability.BOUNDING | capability.PERMITTED | capability.INHERITABLE | capability.EFFECTIVE | capability.AMBIENT); err != nil {
		return fmt.Errorf("apply capability failed: %w", err)
	}

	return nil
}
