package spec

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// == test for func:buildRootSpec ==
func TestBuildRootSpec_1(t *testing.T) {
	input := ConfigOptions{
		Rootfs: "/path/to/root",
	}

	result := buildRootSpec(input)

	expect := RootObject{
		Path: "/path/to/root",
	}

	assert.Equal(t, expect, result)
}

// =================================

// == test for func:buildMountSpec ==
func TestBuildMountSpec_1(t *testing.T) {
	input := ConfigOptions{
		Mounts: []MountOption{
			{
				Destination: "/dst",
				Type:        "",
				Source:      "/src",
				Options: []string{
					"bind",
				},
			},
		},
	}

	result := buildMountSpec(input)

	expect := []MountObject{
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
			Destination: "/sys/fs/cgroup",
			Type:        "cgroup",
			Source:      "cgroup",
			Options: []string{
				"ro",
				"nosuid",
				"noexec",
				"nodev",
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
			Destination: "/dst",
			Type:        "",
			Source:      "/src",
			Options: []string{
				"bind",
			},
		},
	}

	assert.Equal(t, expect, result)
}

// ==================================

// == test for func:buildProcessSpec ==
func TestBuildProcessSpec_1(t *testing.T) {
	input := ConfigOptions{
		Process: ProcessOption{
			Cwd: "/",
			Env: []string{
				"KEY=VALUE",
				"TEST=CODE",
			},
			Args: []string{
				"echo",
				"hello world",
			},
		},
	}

	result := buildProcessSpec(input)

	expect := ProcessObject{
		Cwd: "/",
		Env: []string{
			"KEY=VALUE",
			"TEST=CODE",
		},
		Args: []string{
			"echo",
			"hello world",
		},
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

	assert.Equal(t, expect, result)
}

// ====================================

// == test for func:buildLinuxSpec ==
func TestBuildLinuxSpec_1(t *testing.T) {
	input := ConfigOptions{
		Namespace: []string{
			"mount",
			"network",
		},
	}

	result := buildLinuxSpec(input)

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
		Namespaces: []NamespaceObject{
			{
				Type: "mount",
			},
			{
				Type: "network",
			},
		},
	}

	assert.Equal(t, expect, result)
}

// ==================================

// == test for func:buildNetSpec ==
func TestBuildNetSpec_1(t *testing.T) {
	input := ConfigOptions{
		Net: NetOption{
			InterfaceName: "eth0",
			Address:       "192.168.0.1/24",
			Gateway:       "192.168.0.254",
			Dns: []string{
				"8.8.8.8",
				"8.8.4.4",
			},
		},
	}

	result := buildNetSpec(input)

	expect := NetConfigObject{
		DefaultInterface: "eth0",
		Interface: InterfaceObject{
			Name: "eth0",
			IPv4: IPv4Object{
				Address: "192.168.0.1/24",
				Gateway: "192.168.0.254",
			},
			Dns: DnsObject{
				Servers: []string{
					"8.8.8.8",
					"8.8.4.4",
				},
			},
		},
	}

	assert.Equal(t, expect, result)
}

// ================================

// == test for func:buildImageSpec ==
func TestBuildImageSpec_1(t *testing.T) {
	input := ConfigOptions{
		Image: ImageOption{
			ImageLayer: "/image/path",
			UpperDir:   "/upper/path",
			WorkDir:    "/work/path",
			MergeDir:   "/merge/path",
		},
	}

	result := buildImageSpec(input)

	expect := ImageConfigObject{
		RootfsType: "overlay",
		ImageLayer: "/image/path",
		UpperDir:   "/upper/path",
		WorkDir:    "/work/path",
		MergeDir:   "/merge/path",
	}

	assert.Equal(t, expect, result)
}

// ==================================

// == test for func:buildAnnotationSpec ==
func TestBuildAnnotationSpec_1(t *testing.T) {
	input := ConfigOptions{
		Net: NetOption{
			InterfaceName: "eth0",
			Address:       "192.168.0.1/24",
			Gateway:       "192.168.0.254",
			Dns: []string{
				"8.8.8.8",
				"8.8.4.4",
			},
		},
		Image: ImageOption{
			ImageLayer: "/image/path",
			UpperDir:   "/upper/path",
			WorkDir:    "/work/path",
			MergeDir:   "/merge/path",
		},
	}

	result := buildAnnotationSpec(input)

	expect := AnnotationObject{
		Version: "0.1.0",
		Net:     "{\"defaultInterface\":\"eth0\",\"interface\":{\"name\":\"eth0\",\"ipv4\":{\"address\":\"192.168.0.1/24\",\"gateway\":\"192.168.0.254\"},\"dns\":{\"servers\":[\"8.8.8.8\",\"8.8.4.4\"]}}}",
		Image:   "{\"rootfsType\":\"overlay\",\"imageLayer\":\"/image/path\",\"upperDir\":\"/upper/path\",\"workDir\":\"/work/path\",\"mergeDir\":\"/merge/path\"}",
	}

	assert.Equal(t, expect, result)
}

// =======================================

// == test for func:buildSpec ==
func TestBuildSpec_1(t *testing.T) {
	input := ConfigOptions{
		Rootfs: "/path/to/root",
		Mounts: []MountOption{
			{
				Destination: "/dst",
				Type:        "",
				Source:      "/src",
				Options: []string{
					"bind",
				},
			},
		},
		Process: ProcessOption{
			Cwd: "/",
			Env: []string{
				"KEY=VALUE",
				"TEST=CODE",
			},
			Args: []string{
				"echo",
				"hello world",
			},
		},
		Namespace: []string{
			"mount",
			"network",
		},
		Hostname: "mycontainer",
		Net: NetOption{
			InterfaceName: "eth0",
			Address:       "192.168.0.1/24",
			Gateway:       "192.168.0.254",
			Dns: []string{
				"8.8.8.8",
				"8.8.4.4",
			},
		},
		Image: ImageOption{
			ImageLayer: "/image/path",
			UpperDir:   "/upper/path",
			WorkDir:    "/work/path",
			MergeDir:   "/merge/path",
		},
	}

	result := buildSpec(input)

	expect := Spec{
		OciVersion: "1.2.0",
		Root: RootObject{
			Path: "/path/to/root",
		},
		Mounts: []MountObject{
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
				Destination: "/sys/fs/cgroup",
				Type:        "cgroup",
				Source:      "cgroup",
				Options: []string{
					"ro",
					"nosuid",
					"noexec",
					"nodev",
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
				Destination: "/dst",
				Type:        "",
				Source:      "/src",
				Options: []string{
					"bind",
				},
			},
		},
		Process: ProcessObject{
			Cwd: "/",
			Env: []string{
				"KEY=VALUE",
				"TEST=CODE",
			},
			Args: []string{
				"echo",
				"hello world",
			},
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
			Namespaces: []NamespaceObject{
				{
					Type: "mount",
				},
				{
					Type: "network",
				},
			},
		},
		Annotations: AnnotationObject{
			Version: "0.1.0",
			Net:     "{\"defaultInterface\":\"eth0\",\"interface\":{\"name\":\"eth0\",\"ipv4\":{\"address\":\"192.168.0.1/24\",\"gateway\":\"192.168.0.254\"},\"dns\":{\"servers\":[\"8.8.8.8\",\"8.8.4.4\"]}}}",
			Image:   "{\"rootfsType\":\"overlay\",\"imageLayer\":\"/image/path\",\"upperDir\":\"/upper/path\",\"workDir\":\"/work/path\",\"mergeDir\":\"/merge/path\"}",
		},
	}

	assert.Equal(t, expect, result)
}

// =============================

// == test for func:CreateConfigFile ==
func TestCreateConfigFile_1(t *testing.T) {
	t.Parallel()

	// create temporary directory
	dir := t.TempDir()

	input := ConfigOptions{
		Rootfs: "rootfs",
	}
	path := filepath.Join(dir, "config.json")

	if err := CreateConfigFile(path, input); err != nil {
		t.Fatalf("CreateConfigFile failed: %v", err)
	}

	// check if the file:config.json exists
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}

	// validate output file contents
	var result Spec
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("invalid json written: %v", err)
	}

	expect := buildSpec(input)

	//if diff := cmp.Diff(expect, result); diff != "" {
	//	t.Errorf("config.json mismatch (-want +got):\n%s", diff)
	//}

	assert.Equal(t, expect, result)
}

// ====================================

// == test for func:LoadConfigFile ==
func TestLoadConfigFile_1(t *testing.T) {
	t.Parallel()

	// create temporary directory
	dir := t.TempDir()

	// create configuration file
	path := filepath.Join(dir, "config.json")
	content := `{
					"ociVersion": "1.2.0",
					"root": {
						"path": "/"
					},
					"mounts": [
						{
							"destination": "/proc",
							"type": "proc",
							"source": "proc",
							"options": [
								"nosuid",
								"noexec",
								"nodev"
							]
						},
						{
							"destination": "/dev",
							"type": "tmpfs",
							"source": "tmpfs",
							"options": [
								"nosuid",
								"strictatime",
								"mode=755",
								"size=65536k"
							]
						},
						{
							"destination": "/dev/pts",
							"type": "devpts",
							"source": "devpts",
							"options": [
								"nosuid",
								"noexec",
								"newinstance",
								"ptmxmode=0666",
								"mode=0620",
								"gid=5"
							]
						},
						{
							"destination": "/sys",
							"type": "sysfs",
							"source": "sysfs",
							"options": [
								"nosuid",
								"noexec",
								"nodev",
								"ro"
							]
						},
						{
							"destination": "/sys/fs/cgroup",
							"type": "cgroup",
							"source": "cgroup",
							"options": [
								"ro",
								"nosuid",
								"noexec",
								"nodev"
							]
						},
						{
							"destination": "/dev/mqueue",
							"type": "mqueue",
							"source": "mqueue",
							"options": [
								"nosuid",
								"noexec",
								"nodev"
							]
						},
						{
							"destination": "/dev/shm",
							"type": "tmpfs",
							"source": "shm",
							"options": [
								"nosuid",
								"noexec",
								"nodev",
								"mode=1777",
								"size=67108864"
							]
						},
						{
							"destination": "/dst",
							"type": "",
							"source": "/src",
							"options": [
								"bind"
							]
						},
						{
							"destination": "/dst2",
							"type": "",
							"source": "/src2",
							"options": [
								"bind"
							]
						}
					],
					"process": {
						"cwd": "/",
						"env": [
							"PATH=/usr/bin:/bin",
							"KEY=VALUE"
						],
						"args": [
							"/bin/bash"
						],
						"capabilities": {
							"bounding": [
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
								"CAP_AUDIT_WRITE"
							],
							"permitted": [
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
								"CAP_AUDIT_WRITE"
							],
							"inheritable": null,
							"effective": [
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
								"CAP_AUDIT_WRITE"
							],
							"ambient": null
						}
					},
					"hostname": "mycontainer",
					"linux": {
						"resources": {
							"memory": {
								"limit": 536870912
							},
							"cpu": {
								"period": 100000,
								"quota": 80000
							}
						},
						"namespaces": [
							{
								"type": "mount"
							},
							{
								"type": "network"
							},
							{
								"type": "uts"
							},
							{
								"type": "pid"
							},
							{
								"type": "ipc"
							},
							{
								"type": "cgroup"
							}
						]
					},
					"annotations": {
						"io.raind.runtime.annotation.version": "0.1.0",
						"io.raind.net.config": "{\"defaultInterface\":\"eth0\",\"interface\":{\"name\":\"eth0\",\"ipv4\":{\"address\":\"10.166.0.1/24\",\"gateway\":\"10.166.0.254\"},\"dns\":{\"servers\":[\"8.8.8.8\",\"8.8.4.4\"]}}}",
						"io.raind.image.config": "{\"rootfsType\":\"overlay\",\"imageLayer\":\"/image/path\",\"upperDir\":\"/upper/path\",\"workDir\":\"/work/path\",\"mergeDir\":\"/merge/path\"}"
					}
				}
				`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	result, err := LoadConfigFile(path)
	if err != nil {
		t.Fatalf("LoadConfigFile failed: %v", err)
	}

	expect := Spec{
		OciVersion: "1.2.0",
		Root: RootObject{
			Path: "/",
		},
		Mounts: []MountObject{
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
				Destination: "/sys/fs/cgroup",
				Type:        "cgroup",
				Source:      "cgroup",
				Options: []string{
					"ro",
					"nosuid",
					"noexec",
					"nodev",
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
				Destination: "/dst",
				Type:        "",
				Source:      "/src",
				Options: []string{
					"bind",
				},
			},
			{
				Destination: "/dst2",
				Type:        "",
				Source:      "/src2",
				Options: []string{
					"bind",
				},
			},
		},
		Process: ProcessObject{
			Cwd: "/",
			Env: []string{
				"PATH=/usr/bin:/bin",
				"KEY=VALUE",
			},
			Args: []string{
				"/bin/bash",
			},
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
			Namespaces: []NamespaceObject{
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
					Type: "cgroup",
				},
			},
		},
		Annotations: AnnotationObject{
			Version: "0.1.0",
			Net:     "{\"defaultInterface\":\"eth0\",\"interface\":{\"name\":\"eth0\",\"ipv4\":{\"address\":\"10.166.0.1/24\",\"gateway\":\"10.166.0.254\"},\"dns\":{\"servers\":[\"8.8.8.8\",\"8.8.4.4\"]}}}",
			Image:   "{\"rootfsType\":\"overlay\",\"imageLayer\":\"/image/path\",\"upperDir\":\"/upper/path\",\"workDir\":\"/work/path\",\"mergeDir\":\"/merge/path\"}",
		},
	}

	assert.Equal(t, expect, result)
}

// ==================================
