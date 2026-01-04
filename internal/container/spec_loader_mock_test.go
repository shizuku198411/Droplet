package container

import (
	"droplet/internal/spec"
	"testing"
)

type mockFileSpecLoader struct {
	// loadFile()
	loadFileCallFlag    bool
	loadFileContainerId string
	loadFileSpec        spec.Spec
	loadFileErr         error
}

func (m *mockFileSpecLoader) loadFile(containerId string) (spec.Spec, error) {
	m.loadFileCallFlag = true
	m.loadFileContainerId = containerId
	return m.loadFileSpec, m.loadFileErr
}

func buildMockSpec(t *testing.T) spec.Spec {
	t.Helper()

	return spec.Spec{
		OciVersion: "1.2.0",
		Root: spec.RootObject{
			Path: "/etc/raind/container/mycontainer/merge",
		},
		Mounts: []spec.MountObject{
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
				Destination: "/etc/resolv.conf",
				Type:        "bind",
				Source:      "/etc/raind/container/mycontainer/etc/resolv.conf",
				Options: []string{
					"rbind",
					"rprivate",
				},
			},
			{
				Destination: "/etc/hostname",
				Type:        "bind",
				Source:      "/etc/raind/container/mycontainer/etc/hostname",
				Options: []string{
					"rbind",
					"rprivate",
				},
			},
			{
				Destination: "/etc/hosts",
				Type:        "bind",
				Source:      "/etc/raind/container/mycontainer/etc/hosts",
				Options: []string{
					"rbind",
					"rprivate",
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
		Process: spec.ProcessObject{
			Cwd:  "/",
			Env:  []string{"KEY=VALUE"},
			Args: []string{"/bin/sh"},
			Capabilities: spec.CapabilityObject{
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
		LinuxSpec: spec.LinuxSpecObject{
			Resources: spec.ResourceObject{
				Memory: spec.MemoryObject{
					Limit: 536870912,
				},
				Cpu: spec.CpuObject{
					Period: 100000,
					Quota:  80000,
				},
			},
			Namespaces: []spec.NamespaceObject{
				{
					Type: "mount",
				},
			},
		},
		Annotations: spec.AnnotationObject{
			Version: "0.1.0",
			Net:     "{\"defaultInterface\":\"eth0\",\"interface\":{\"name\":\"eth0\",\"ipv4\":{\"address\":\"10.166.0.1/24\",\"gateway\":\"10.166.0.254\"},\"dns\":{\"servers\":[\"8.8.8.8\"]}}}",
			Image:   "{\"rootfsType\":\"overlay\",\"imageLayer\":[\"/image/path\"],\"upperDir\":\"/upper/path\",\"workDir\":\"/work/path\",\"mergeDir\":\"/merge/path\"}",
		},
	}
}
