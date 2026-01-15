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
		OciVersion: "1.3.0",
		Root: spec.RootObject{
			Path: "/etc/raind/container/mycontainer/merge",
		},
		Mounts: []spec.MountObject{
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
