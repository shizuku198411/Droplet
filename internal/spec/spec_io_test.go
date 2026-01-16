package spec

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

// build ConfigOptions helper
func buildConfigOptions(t *testing.T) ConfigOptions {
	t.Helper()

	return ConfigOptions{
		Rootfs: "rootfs",
		Mounts: []MountOption{
			{
				Destination: "/dst",
				Type:        "",
				Source:      "/src",
				Options:     []string{"bind"},
			},
		},
		Process: ProcessOption{
			Cwd:  "/",
			Env:  []string{"KEY=VALUE"},
			Args: []string{"/bin/sh"},
		},
		Namespace: []string{"mount"},
		Hostname:  "mycontainer",
		Net: NetOption{
			HostInterface:       "eth0",
			BridgeInterfaceName: "br0",
			InterfaceName:       "eth0",
			Address:             "10.166.0.1/24",
			Gateway:             "10.166.0.254",
			Dns:                 []string{"8.8.8.8"},
		},
		Image: ImageOption{
			ImageLayer: []string{"/image/path"},
			UpperDir:   "/upper/path",
			WorkDir:    "/work/path",
		},
	}
}

func TestBuildRootSpec_Success(t *testing.T) {
	// == arrange ==
	opts := buildConfigOptions(t)

	// == act ==
	got := buildRootSpec(opts)

	// == assert ==
	assert.Equal(t, RootObject{Path: "rootfs"}, got)
}

func TestBuildMountSpec_Success(t *testing.T) {
	// == arrange ==
	opts := buildConfigOptions(t)

	// == act ==
	got := buildMountSpec(opts)

	// == assert ==
	expect := []MountObject{
		{
			Destination: "/dst",
			Type:        "",
			Source:      "/src",
			Options: []string{
				"bind",
			},
		},
	}
	assert.Equal(t, expect, got)
}

func TestBuildProcessSpec_Success(t *testing.T) {
	// == arrange ==
	opts := buildConfigOptions(t)

	// == act ==
	got := buildProcessSpec(opts)

	// == assert ==
	expect := ProcessObject{
		Cwd:  "/",
		Env:  []string{"KEY=VALUE"},
		Args: []string{"/bin/sh"},
		Capabilities: CapabilityObject{
			Bounding: []string{
				"CAP_CHOWN",
				"CAP_DAC_OVERRIDE",
				"CAP_FSETID",
				"CAP_FOWNER",
				"CAP_MKNOD",
				"CAP_NET_RAW",
				"CAP_SETGID",
				"CAP_SETUID",
				"CAP_SETFCAP",
				"CAP_SETPCAP",
				"CAP_NET_BIND_SERVICE",
				"CAP_SYS_CHROOT",
				"CAP_KILL",
				"CAP_AUDIT_WRITE",
			},
			Effective: []string{
				"CAP_CHOWN",
				"CAP_DAC_OVERRIDE",
				"CAP_FSETID",
				"CAP_FOWNER",
				"CAP_MKNOD",
				"CAP_NET_RAW",
				"CAP_SETGID",
				"CAP_SETUID",
				"CAP_SETFCAP",
				"CAP_SETPCAP",
				"CAP_NET_BIND_SERVICE",
				"CAP_SYS_CHROOT",
				"CAP_KILL",
				"CAP_AUDIT_WRITE",
			},
			Permitted: []string{
				"CAP_CHOWN",
				"CAP_DAC_OVERRIDE",
				"CAP_FSETID",
				"CAP_FOWNER",
				"CAP_MKNOD",
				"CAP_NET_RAW",
				"CAP_SETGID",
				"CAP_SETUID",
				"CAP_SETFCAP",
				"CAP_SETPCAP",
				"CAP_NET_BIND_SERVICE",
				"CAP_SYS_CHROOT",
				"CAP_KILL",
				"CAP_AUDIT_WRITE",
			},
		},
	}
	assert.Equal(t, expect, got)
}

func TestBuildLinuxSpec_Success(t *testing.T) {
	// == arrange ==
	opts := buildConfigOptions(t)

	// == act ==
	got := buildLinuxSpec(opts)

	// == assert ==
	ep := uint32(1)
	ociArch := func() string {
		arch := runtime.GOARCH
		switch arch {
		case "amd64":
			return "SCMP_ARCH_X86_64"
		case "arm64":
			return "SCMP_ARCH_AARCH64"
		case "riscv64":
			return "SCMP_ARCH_RISCV64"
		default:
			return ""
		}
	}

	expect := LinuxSpecObject{
		Resources: ResourceObject{
			Memory: MemoryObject{
				Limit: 536870912,
			},
			Cpu: CpuObject{
				Period: 100000,
				Quota:  80000,
			},
		},
		Seccomp: &SeccompObject{
			DefaultAction:   "SCMP_ACT_ALLOW",
			DefaultErrnoRet: &ep,
			Architectures: []string{
				ociArch(),
			},
			Syscalls: []SeccompSyscallObject{
				{
					Names: []string{
						"bpf",
						"perf_event_open",
						"kexec_load",
						"open_by_handle_at",
						"ptrace",
						"process_vm_readv",
						"process_vm_writev",
						"userfaultfd",
						"reboot",
						"swapon",
						"swapoff",
						"open_by_handle_at",
						"name_to_handle_at",
						"init_module",
						"finit_module",
						"delete_module",
						"kcmp",
						"mount",
						"unshare",
						"setns",
					},
					Action:   "SCMP_ACT_ERRNO",
					ErrnoRet: &ep,
				},
			},
		},
		Namespaces: []NamespaceObject{
			{
				Type: "mount",
			},
		},
	}
	assert.Equal(t, expect, got)
}

func TestBuildNetSpec_Success(t *testing.T) {
	// == arrange ==
	opts := buildConfigOptions(t)

	// == act ==
	got := buildNetSpec(opts)

	// == assert ==
	expect := NetConfigObject{
		HostInterface:   "eth0",
		BridgeInterface: "br0",
		Interface: InterfaceObject{
			Name: "eth0",
			IPv4: IPv4Object{
				Address: "10.166.0.1/24",
				Gateway: "10.166.0.254",
			},
			Dns: DnsObject{
				Servers: []string{"8.8.8.8"},
			},
		},
	}
	assert.Equal(t, expect, got)
}

func TestBuildImageSpec_Success(t *testing.T) {
	// == arrange ==
	opts := buildConfigOptions(t)

	// == act ==
	got := buildImageSpec(opts)

	// == assert ==
	expect := ImageConfigObject{
		RootfsType: "overlay",
		ImageLayer: []string{"/image/path"},
		UpperDir:   "/upper/path",
		WorkDir:    "/work/path",
	}
	assert.Equal(t, expect, got)
}

func TestBuildAnnotationSpec_Success(t *testing.T) {
	// == arrange ==
	opts := buildConfigOptions(t)

	// == act ==
	got := buildAnnotationSpec(opts)

	// == assert ==
	expect := AnnotationObject{
		Version: "0.1.0",
		Net:     "{\"hostInterface\":\"eth0\",\"bridgeInterface\":\"br0\",\"interface\":{\"name\":\"eth0\",\"ipv4\":{\"address\":\"10.166.0.1/24\",\"gateway\":\"10.166.0.254\"},\"dns\":{\"servers\":[\"8.8.8.8\"]}}}",
		Image:   "{\"rootfsType\":\"overlay\",\"imageLayer\":[\"/image/path\"],\"upperDir\":\"/upper/path\",\"workDir\":\"/work/path\"}",
	}
	assert.Equal(t, expect, got)
}

func TestBuildSpec_Success(t *testing.T) {
	// == arrange ==
	opts := buildConfigOptions(t)

	// == act ==
	got := buildSpec(opts)

	// == assert ==
	ep := uint32(1)
	ociArch := func() string {
		arch := runtime.GOARCH
		switch arch {
		case "amd64":
			return "SCMP_ARCH_X86_64"
		case "arm64":
			return "SCMP_ARCH_AARCH64"
		case "riscv64":
			return "SCMP_ARCH_RISCV64"
		default:
			return ""
		}
	}

	expect := Spec{
		OciVersion: "1.3.0",
		Root: RootObject{
			Path: "rootfs",
		},
		Mounts: []MountObject{
			{
				Destination: "/dst",
				Type:        "",
				Source:      "/src",
				Options: []string{
					"bind",
				},
			},
		},
		Process: ProcessObject{
			Cwd:  "/",
			Env:  []string{"KEY=VALUE"},
			Args: []string{"/bin/sh"},
			Capabilities: CapabilityObject{
				Bounding: []string{
					"CAP_CHOWN",
					"CAP_DAC_OVERRIDE",
					"CAP_FSETID",
					"CAP_FOWNER",
					"CAP_MKNOD",
					"CAP_NET_RAW",
					"CAP_SETGID",
					"CAP_SETUID",
					"CAP_SETFCAP",
					"CAP_SETPCAP",
					"CAP_NET_BIND_SERVICE",
					"CAP_SYS_CHROOT",
					"CAP_KILL",
					"CAP_AUDIT_WRITE",
				},
				Effective: []string{
					"CAP_CHOWN",
					"CAP_DAC_OVERRIDE",
					"CAP_FSETID",
					"CAP_FOWNER",
					"CAP_MKNOD",
					"CAP_NET_RAW",
					"CAP_SETGID",
					"CAP_SETUID",
					"CAP_SETFCAP",
					"CAP_SETPCAP",
					"CAP_NET_BIND_SERVICE",
					"CAP_SYS_CHROOT",
					"CAP_KILL",
					"CAP_AUDIT_WRITE",
				},
				Permitted: []string{
					"CAP_CHOWN",
					"CAP_DAC_OVERRIDE",
					"CAP_FSETID",
					"CAP_FOWNER",
					"CAP_MKNOD",
					"CAP_NET_RAW",
					"CAP_SETGID",
					"CAP_SETUID",
					"CAP_SETFCAP",
					"CAP_SETPCAP",
					"CAP_NET_BIND_SERVICE",
					"CAP_SYS_CHROOT",
					"CAP_KILL",
					"CAP_AUDIT_WRITE",
				},
			},
		},
		Hostname: "mycontainer",
		LinuxSpec: LinuxSpecObject{
			Resources: ResourceObject{
				Memory: MemoryObject{
					Limit: 536870912,
				},
				Cpu: CpuObject{
					Period: 100000,
					Quota:  80000,
				},
			},
			Seccomp: &SeccompObject{
				DefaultAction:   "SCMP_ACT_ALLOW",
				DefaultErrnoRet: &ep,
				Architectures: []string{
					ociArch(),
				},
				Syscalls: []SeccompSyscallObject{
					{
						Names: []string{
							"bpf",
							"perf_event_open",
							"kexec_load",
							"open_by_handle_at",
							"ptrace",
							"process_vm_readv",
							"process_vm_writev",
							"userfaultfd",
							"reboot",
							"swapon",
							"swapoff",
							"open_by_handle_at",
							"name_to_handle_at",
							"init_module",
							"finit_module",
							"delete_module",
							"kcmp",
							"mount",
							"unshare",
							"setns",
						},
						Action:   "SCMP_ACT_ERRNO",
						ErrnoRet: &ep,
					},
				},
			},
			Namespaces: []NamespaceObject{
				{
					Type: "mount",
				},
			},
		},
		Annotations: AnnotationObject{
			Version: "0.1.0",
			Net:     "{\"hostInterface\":\"eth0\",\"bridgeInterface\":\"br0\",\"interface\":{\"name\":\"eth0\",\"ipv4\":{\"address\":\"10.166.0.1/24\",\"gateway\":\"10.166.0.254\"},\"dns\":{\"servers\":[\"8.8.8.8\"]}}}",
			Image:   "{\"rootfsType\":\"overlay\",\"imageLayer\":[\"/image/path\"],\"upperDir\":\"/upper/path\",\"workDir\":\"/work/path\"}",
		},
	}
	assert.Equal(t, expect, got)
}

func TestCreateConfigFile_Success(t *testing.T) {
	// == arrange ==
	opts := buildConfigOptions(t)
	path := filepath.Join(t.TempDir(), "config.json")

	// == act ==
	err := CreateConfigFile(path, opts)

	// == assert ==
	assert.Nil(t, err)
}

func TestCreateConfigFile_PathNotExistsError(t *testing.T) {
	// == arrange ==
	opts := buildConfigOptions(t)
	path := "/not/exists/path"

	// == act ==
	err := CreateConfigFile(path, opts)

	// == assert ==
	assert.NotNil(t, err)
}

func TestLoadConfigFile_Success(t *testing.T) {
	// == arrange ==
	opts := buildConfigOptions(t)
	path := filepath.Join(t.TempDir(), "config.json")
	configCreateErr := CreateConfigFile(path, opts)
	if configCreateErr != nil {
		t.Fatalf("create config.json failed")
	}

	// == act ==
	got, err := LoadConfigFile(path)

	// == assert ==
	assert.Nil(t, err)

	ep := uint32(1)
	ociArch := func() string {
		arch := runtime.GOARCH
		switch arch {
		case "amd64":
			return "SCMP_ARCH_X86_64"
		case "arm64":
			return "SCMP_ARCH_AARCH64"
		case "riscv64":
			return "SCMP_ARCH_RISCV64"
		default:
			return ""
		}
	}

	expect := Spec{
		OciVersion: "1.3.0",
		Root: RootObject{
			Path: "rootfs",
		},
		Mounts: []MountObject{
			{
				Destination: "/dst",
				Type:        "",
				Source:      "/src",
				Options: []string{
					"bind",
				},
			},
		},
		Process: ProcessObject{
			Cwd:  "/",
			Env:  []string{"KEY=VALUE"},
			Args: []string{"/bin/sh"},
			Capabilities: CapabilityObject{
				Bounding: []string{
					"CAP_CHOWN",
					"CAP_DAC_OVERRIDE",
					"CAP_FSETID",
					"CAP_FOWNER",
					"CAP_MKNOD",
					"CAP_NET_RAW",
					"CAP_SETGID",
					"CAP_SETUID",
					"CAP_SETFCAP",
					"CAP_SETPCAP",
					"CAP_NET_BIND_SERVICE",
					"CAP_SYS_CHROOT",
					"CAP_KILL",
					"CAP_AUDIT_WRITE",
				},
				Effective: []string{
					"CAP_CHOWN",
					"CAP_DAC_OVERRIDE",
					"CAP_FSETID",
					"CAP_FOWNER",
					"CAP_MKNOD",
					"CAP_NET_RAW",
					"CAP_SETGID",
					"CAP_SETUID",
					"CAP_SETFCAP",
					"CAP_SETPCAP",
					"CAP_NET_BIND_SERVICE",
					"CAP_SYS_CHROOT",
					"CAP_KILL",
					"CAP_AUDIT_WRITE",
				},
				Permitted: []string{
					"CAP_CHOWN",
					"CAP_DAC_OVERRIDE",
					"CAP_FSETID",
					"CAP_FOWNER",
					"CAP_MKNOD",
					"CAP_NET_RAW",
					"CAP_SETGID",
					"CAP_SETUID",
					"CAP_SETFCAP",
					"CAP_SETPCAP",
					"CAP_NET_BIND_SERVICE",
					"CAP_SYS_CHROOT",
					"CAP_KILL",
					"CAP_AUDIT_WRITE",
				},
			},
		},
		Hostname: "mycontainer",
		LinuxSpec: LinuxSpecObject{
			Resources: ResourceObject{
				Memory: MemoryObject{
					Limit: 536870912,
				},
				Cpu: CpuObject{
					Period: 100000,
					Quota:  80000,
				},
			},
			Seccomp: &SeccompObject{
				DefaultAction:   "SCMP_ACT_ALLOW",
				DefaultErrnoRet: &ep,
				Architectures: []string{
					ociArch(),
				},
				Syscalls: []SeccompSyscallObject{
					{
						Names: []string{
							"bpf",
							"perf_event_open",
							"kexec_load",
							"open_by_handle_at",
							"ptrace",
							"process_vm_readv",
							"process_vm_writev",
							"userfaultfd",
							"reboot",
							"swapon",
							"swapoff",
							"open_by_handle_at",
							"name_to_handle_at",
							"init_module",
							"finit_module",
							"delete_module",
							"kcmp",
							"mount",
							"unshare",
							"setns",
						},
						Action:   "SCMP_ACT_ERRNO",
						ErrnoRet: &ep,
					},
				},
			},
			Namespaces: []NamespaceObject{
				{
					Type: "mount",
				},
			},
		},
		Annotations: AnnotationObject{
			Version: "0.1.0",
			Net:     "{\"hostInterface\":\"eth0\",\"bridgeInterface\":\"br0\",\"interface\":{\"name\":\"eth0\",\"ipv4\":{\"address\":\"10.166.0.1/24\",\"gateway\":\"10.166.0.254\"},\"dns\":{\"servers\":[\"8.8.8.8\"]}}}",
			Image:   "{\"rootfsType\":\"overlay\",\"imageLayer\":[\"/image/path\"],\"upperDir\":\"/upper/path\",\"workDir\":\"/work/path\"}",
		},
	}
	assert.Equal(t, expect, got)
}

func TestLoadConfigFile_FileNotExistsError(t *testing.T) {
	// == arrange ==
	path := filepath.Join(t.TempDir(), "config.json")

	// == act ==
	_, err := LoadConfigFile(path)

	// == assert ==
	assert.NotNil(t, err)
}
