package container

import (
	"syscall"

	"droplet/internal/spec"
)

type namespaceConfig struct {
	mount   bool
	network bool
	uts     bool
	pid     bool
	ipc     bool
	user    bool
	cgroup  bool
}

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

func buildNamespaceAttr(nsConfig namespaceConfig) *syscall.SysProcAttr {
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

	return &syscall.SysProcAttr{
		Cloneflags: flags,
	}
}
