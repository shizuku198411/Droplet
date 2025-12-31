package spec

import (
	"droplet/internal/utils"
	"path/filepath"
)

func buildRootSpec(opts ConfigOptions) RootObject {
	return RootObject{
		Path: opts.Rootfs,
	}
}

func buildMountSpec(opts ConfigOptions) []MountObject {
	// mount
	// the following device must be mounted:
	//   /proc, /dev, /dev/pts, /sys, /sys/fs/cgroup, /dev/mqueue, /dev/shm
	var mounts = []MountObject{
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
	}

	// user mounts
	for _, user_mount := range opts.Mounts {
		mounts = append(mounts, MountObject{
			Destination: user_mount.Destination,
			Type:        user_mount.Type,
			Source:      user_mount.Source,
			Options:     user_mount.Options,
		})
	}

	return mounts
}

func buildProcessSpec(opts ConfigOptions) ProcessObject {
	return ProcessObject{
		Cwd:  opts.Process.Cwd,
		Env:  opts.Process.Env,
		Args: opts.Process.Args,
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
}

func buildLinuxSpec() LinuxSpecObject {
	return LinuxSpecObject{
		Resources: ResourceObject{
			Memory: MemoryObject{ // memory limit: 512MiB
				Limit: 536870912,
			},
			Cpu: CpuObject{ // cpu limit: 80%
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
	}
}

func buildNetSpec(opts ConfigOptions) NetConfigObject {
	return NetConfigObject{
		DefaultInterface: opts.Net.InterfaceName,
		Interface: InterfaceObject{
			Name: opts.Net.InterfaceName,
			IPv4: IPv4Object{
				Address: opts.Net.Address,
				Gateway: opts.Net.Gateway,
			},
			Dns: DnsObject{
				Servers: opts.Net.Dns,
			},
		},
	}
}

func buildImageSpec(opts ConfigOptions) ImageConfigObject {
	return ImageConfigObject{
		RootfsType: "overlay",
		ImageLayer: opts.Image.ImageLayer,
		UpperDir:   opts.Image.UpperDir,
		WorkDir:    opts.Image.WorkDir,
		MergeDir:   opts.Image.MergeDir,
	}
}

func buildAnnotationSpec(opts ConfigOptions) AnnotationObject {
	netSpec, _ := utils.JsonToString(buildNetSpec(opts))
	imageSpec, _ := utils.JsonToString(buildImageSpec(opts))
	return AnnotationObject{
		Version: "0.1.0", // annotation version: 0.1.0
		Net:     netSpec,
		Image:   imageSpec,
	}
}

func buildSpec(opts ConfigOptions) Spec {
	// OCI Version: 1.2.0
	ociVersion := "1.2.0"

	// root path
	root := buildRootSpec(opts)

	// mounts
	mounts := buildMountSpec(opts)

	// process
	process := buildProcessSpec(opts)

	// hostname
	hostname := opts.Hostname

	// linux spec
	linuxSpec := buildLinuxSpec()

	// annotation
	annotation := buildAnnotationSpec(opts)

	return Spec{
		OciVersion:  ociVersion,
		Root:        root,
		Mounts:      mounts,
		Process:     process,
		Hostname:    hostname,
		LinuxSpec:   linuxSpec,
		Annotations: annotation,
	}
}

func CreateConfigFile(path string, opts ConfigOptions) error {
	// build spec
	spec := buildSpec(opts)

	// write spec to file
	configPath := filepath.Join(path)
	if err := utils.WriteJsonToFile(configPath, spec); err != nil {
		return err
	}
	return nil
}

func LoadConfigFile(path string) (Spec, error) {
	var spec Spec

	if err := utils.ReadJsonFile(path, &spec); err != nil {
		return Spec{}, err
	}

	return spec, nil
}
