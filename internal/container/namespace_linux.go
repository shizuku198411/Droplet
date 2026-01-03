package container

import (
	"syscall"

	"droplet/internal/spec"
)

// procAttr represents the low-level process attributes that will be applied
// when starting the container init process.
//
// At present it contains only the cloneFlags derived from the selected
// namespaces, but the struct exists to allow future extension (for example,
// UID/GID mappings, capability settings, or seccomp configuration) without
// changing the function signatures that depend on it.
type procAttr struct {
	cloneFlags    uintptr
	uidMap        []syscall.SysProcIDMap
	gidMap        []syscall.SysProcIDMap
	setGroupsFlag bool
}

// buildSysProcAttr converts the given procAttr into a syscall.SysProcAttr,
// which can be assigned to exec.Cmd.SysProcAttr when launching the init
// process.
//
// The returned SysProcAttr currently sets only the Cloneflags field, but
// additional process attributes (such as UID/GID mappings for user
// namespaces) may be added here in the future as the runtime evolves.
func buildSysProcAttr(procAttr procAttr) *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
		Cloneflags:                 procAttr.cloneFlags,
		UidMappings:                procAttr.uidMap,
		GidMappings:                procAttr.gidMap,
		GidMappingsEnableSetgroups: procAttr.setGroupsFlag,
	}
}

// buildProcAttrForRootContainer builds a procAttr for a root-executed
// container, including user namespaces if requested in nsConfig.
//
// For now this configures a simple identity mapping (0 -> 0, size 65535)
// when the user namespace is enabled. Other namespaces are expressed as
// clone flags only. This function can be extended to support additional
// mappings (e.g. non-root users, rootless mode, or spec-based mappings)
// without changing the SysProcAttr construction logic.
func buildProcAttrForRootContainer(nsConfig namespaceConfig) procAttr {
	cloneFlags := buildCloneFlags(nsConfig)
	uidMap, gidMap := buildRootUserNamespaceIDMap(nsConfig)
	setGroupsFlag := true

	return procAttr{
		cloneFlags:    cloneFlags,
		uidMap:        uidMap,
		gidMap:        gidMap,
		setGroupsFlag: setGroupsFlag,
	}
}

// namespaceConfig represents the set of Linux namespaces that should be
// created for the container's init process.
//
// Each field corresponds to an OCI runtime-spec namespace type.
// A value of true indicates that the namespace should be created
// (i.e., the associated CLONE_NEW* flag will be applied).
type namespaceConfig struct {
	mount   bool
	network bool
	uts     bool
	pid     bool
	ipc     bool
	user    bool
	cgroup  bool
}

// buildNamespaceConfig constructs a namespaceConfig from the namespaces
// defined in the OCI runtime-spec.
//
// The function inspects spec.LinuxSpec.Namespaces and marks each namespace
// as enabled in the returned namespaceConfig. If a namespace type is not
// present in the spec, the corresponding field remains false.
//
// This function does not perform any system calls; it simply derives the
// configuration that will later be used to construct SysProcAttr.
func buildNamespaceConfig(spec spec.Spec) namespaceConfig {
	var nsConfig namespaceConfig
	for _, ns := range spec.LinuxSpec.Namespaces {
		switch ns.Type {
		case "mount":
			nsConfig.mount = true
		case "network":
			nsConfig.network = true
		case "uts":
			nsConfig.uts = true
		case "pid":
			nsConfig.pid = true
		case "ipc":
			nsConfig.ipc = true
		case "user":
			nsConfig.user = true
		case "cgroup":
			nsConfig.cgroup = true
		}
	}
	return nsConfig
}

// buildCloneFlags constructs the Linux namespace clone flags from the given
// namespaceConfig and returns the bitwise OR of the corresponding CLONE_NEW*
//
// Each enabled namespace in nsConfig results in the associated clone flag
// being added to the returned value. The resulting flags value is intended
// to be used as syscall.SysProcAttr.Cloneflags when spawning the container
// init process.
//
// This function does not perform any system calls; it only derives the flag
// mask based on the requested namespaces.
func buildCloneFlags(nsConfig namespaceConfig) uintptr {
	var flags uintptr

	if nsConfig.mount {
		flags |= syscall.CLONE_NEWNS
	}
	if nsConfig.network {
		flags |= syscall.CLONE_NEWNET
	}
	if nsConfig.uts {
		flags |= syscall.CLONE_NEWUTS
	}
	if nsConfig.pid {
		flags |= syscall.CLONE_NEWPID
	}
	if nsConfig.ipc {
		flags |= syscall.CLONE_NEWIPC
	}
	if nsConfig.user {
		flags |= syscall.CLONE_NEWUSER
	}
	if nsConfig.cgroup {
		flags |= syscall.CLONE_NEWCGROUP
	}

	return flags
}

// buildRootUserNamespaceIDMaps returns UID/GID ID maps suitable for a
// root-executed container when the user namespace is enabled.
//
// When nsConfig.user is true, this function creates an identity mapping
// from container UID/GID 0..(size-1) to host UID/GID 0..(size-1).
// When the user namespace is disabled, it returns nil maps.
//
// This function is the main extension point for future mapping policies:
// for example, supporting rootless containers, using /etc/subuid/subgid,
// or honoring OCI spec.Process.User fields.
func buildRootUserNamespaceIDMap(nsConfig namespaceConfig) (uidMap, gidMap []syscall.SysProcIDMap) {
	if !nsConfig.user {
		return nil, nil
	}

	const idMapSize = 65535

	uidMap = []syscall.SysProcIDMap{
		{
			ContainerID: 0,
			HostID:      0,
			Size:        idMapSize,
		},
	}
	gidMap = []syscall.SysProcIDMap{
		{
			ContainerID: 0,
			HostID:      0,
			Size:        idMapSize,
		},
	}

	return uidMap, gidMap
}
