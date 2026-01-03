package container

import (
	"droplet/internal/spec"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildNamespaceConfig_Success(t *testing.T) {
	// == arrange ==
	spec := spec.Spec{
		LinuxSpec: spec.LinuxSpecObject{
			Namespaces: []spec.NamespaceObject{
				{
					Type: "mount",
				},
				{
					Type: "network",
				},
				{
					Type: "uts",
				},
				{
					Type: "pid",
				},
				{
					Type: "ipc",
				},
				{
					Type: "user",
				},
				{
					Type: "cgroup",
				},
			},
		},
	}

	// == act ==
	nsConfig := buildNamespaceConfig(spec)

	// == assert ==
	assert.True(t, nsConfig.mount)
	assert.True(t, nsConfig.network)
	assert.True(t, nsConfig.uts)
	assert.True(t, nsConfig.pid)
	assert.True(t, nsConfig.ipc)
	assert.True(t, nsConfig.user)
	assert.True(t, nsConfig.cgroup)
}

func TestBuildCloneFlags_Success(t *testing.T) {
	// == arrange ==
	nsConfig := namespaceConfig{
		mount:   true,
		network: true,
		uts:     true,
		pid:     true,
		ipc:     true,
		user:    true,
		cgroup:  true,
	}

	// == act ==
	got := buildCloneFlags(nsConfig)

	// == assert ==
	var expect uintptr
	expect |= (syscall.CLONE_NEWNS |
		syscall.CLONE_NEWNET |
		syscall.CLONE_NEWUTS |
		syscall.CLONE_NEWPID |
		syscall.CLONE_NEWIPC |
		syscall.CLONE_NEWUSER |
		syscall.CLONE_NEWCGROUP)
	assert.Equal(t, expect, got)
}

func TestBuildRootUserNamespaceIDMap_NsUserSuccess(t *testing.T) {
	// == arrange ==
	nsConfig := namespaceConfig{
		mount:   true,
		network: true,
		uts:     true,
		pid:     true,
		ipc:     true,
		user:    true,
		cgroup:  true,
	}

	// == arrange ==
	uidMap, gidMap := buildRootUserNamespaceIDMap(nsConfig)

	// == assert ==
	expectUidMap := []syscall.SysProcIDMap{
		{
			ContainerID: 0,
			HostID:      0,
			Size:        65535,
		},
	}
	expectGidMap := []syscall.SysProcIDMap{
		{
			ContainerID: 0,
			HostID:      0,
			Size:        65535,
		},
	}
	assert.Equal(t, expectUidMap, uidMap)
	assert.Equal(t, expectGidMap, gidMap)
}

func TestBuildRootUserNamespaceIDMap_NonNsUserSuccess(t *testing.T) {
	// == arrange ==
	nsConfig := namespaceConfig{
		mount:   true,
		network: true,
		uts:     true,
		pid:     true,
		ipc:     true,
		user:    false,
		cgroup:  true,
	}

	// == arrange ==
	uidMap, gidMap := buildRootUserNamespaceIDMap(nsConfig)

	// == assert ==
	assert.Nil(t, uidMap)
	assert.Nil(t, gidMap)
}

func TestBuildProcAttrForRootContainer_Success(t *testing.T) {
	// == arrange ==
	nsConfig := namespaceConfig{
		mount:   true,
		network: true,
		uts:     true,
		pid:     true,
		ipc:     true,
		user:    true,
		cgroup:  true,
	}

	// == act ==
	got := buildProcAttrForRootContainer(nsConfig)

	// == assert ==
	expect := procAttr{
		cloneFlags: (syscall.CLONE_NEWNS |
			syscall.CLONE_NEWNET |
			syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWUSER |
			syscall.CLONE_NEWCGROUP),
		uidMap: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      0,
				Size:        65535,
			},
		},
		gidMap: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      0,
				Size:        65535,
			},
		},
		setGroupsFlag: true,
	}
	assert.Equal(t, expect, got)
}

func TestBuildSysProcAttr_Success(t *testing.T) {
	// == arrange ==
	procAttr := procAttr{
		cloneFlags: (syscall.CLONE_NEWNS |
			syscall.CLONE_NEWNET |
			syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWUSER |
			syscall.CLONE_NEWCGROUP),
		uidMap: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      0,
				Size:        65535,
			},
		},
		gidMap: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      0,
				Size:        65535,
			},
		},
		setGroupsFlag: true,
	}

	// == act ==
	got := buildSysProcAttr(procAttr)

	// == assert ==
	expect := &syscall.SysProcAttr{
		Cloneflags: (syscall.CLONE_NEWNS |
			syscall.CLONE_NEWNET |
			syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWUSER |
			syscall.CLONE_NEWCGROUP),
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      0,
				Size:        65535,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      0,
				Size:        65535,
			},
		},
		GidMappingsEnableSetgroups: true,
	}
	assert.Equal(t, expect, got)
}
