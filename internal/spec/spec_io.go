package spec

import (
	"droplet/internal/oci"
	"droplet/internal/utils"
	"path/filepath"
)

func buildRootSpec(opts ConfigOptions) RootObject {
	return RootObject{
		Path: opts.Rootfs,
	}
}

func buildMountSpec(opts ConfigOptions) []MountObject {
	var mounts = []MountObject{}

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

func buildLinuxSpec(opts ConfigOptions) LinuxSpecObject {
	var linuxSpec = LinuxSpecObject{
		Resources: ResourceObject{
			Memory: MemoryObject{ // memory limit: 512MiB
				Limit: 536870912,
			},
			Cpu: CpuObject{ // cpu limit: 80%
				Period: 100000,
				Quota:  80000,
			},
		},
		Namespaces: []NamespaceObject{},
	}

	for _, ns := range opts.Namespace {
		linuxSpec.Namespaces = append(linuxSpec.Namespaces, NamespaceObject{
			Type: ns,
		})
	}

	return linuxSpec
}

func buildNetSpec(opts ConfigOptions) NetConfigObject {
	return NetConfigObject{
		HostInterface:   opts.Net.HostInterface,
		BridgeInterface: opts.Net.BridgeInterfaceName,
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
	}
}

func buildHookSpec(opts ConfigOptions) HookLifecycleObject {
	var hookLifeCycleObject HookLifecycleObject

	// prestart
	for _, h := range opts.Hooks.Prestart {
		hookLifeCycleObject.Prestart = append(hookLifeCycleObject.Prestart,
			HookObject{
				Path:    h.Path,
				Args:    h.Args,
				Env:     h.Env,
				Timeout: h.Timeout,
			},
		)
	}
	// crateRuntime
	for _, h := range opts.Hooks.CreateRuntime {
		hookLifeCycleObject.CreateRuntime = append(hookLifeCycleObject.CreateRuntime,
			HookObject{
				Path:    h.Path,
				Args:    h.Args,
				Env:     h.Env,
				Timeout: h.Timeout,
			},
		)
	}
	// crateContainer
	for _, h := range opts.Hooks.CreateContainer {
		hookLifeCycleObject.CreateContainer = append(hookLifeCycleObject.CreateContainer,
			HookObject{
				Path:    h.Path,
				Args:    h.Args,
				Env:     h.Env,
				Timeout: h.Timeout,
			},
		)
	}
	// startContainer
	for _, h := range opts.Hooks.StartContainer {
		hookLifeCycleObject.StartContainer = append(hookLifeCycleObject.StartContainer,
			HookObject{
				Path:    h.Path,
				Args:    h.Args,
				Env:     h.Env,
				Timeout: h.Timeout,
			},
		)
	}
	// poststart
	for _, h := range opts.Hooks.Poststart {
		hookLifeCycleObject.Poststart = append(hookLifeCycleObject.Poststart,
			HookObject{
				Path:    h.Path,
				Args:    h.Args,
				Env:     h.Env,
				Timeout: h.Timeout,
			},
		)
	}
	// stopContainer
	for _, h := range opts.Hooks.StopContainer {
		hookLifeCycleObject.StopContainer = append(hookLifeCycleObject.StopContainer,
			HookObject{
				Path:    h.Path,
				Args:    h.Args,
				Env:     h.Env,
				Timeout: h.Timeout,
			},
		)
	}
	// poststop
	for _, h := range opts.Hooks.Poststop {
		hookLifeCycleObject.Poststop = append(hookLifeCycleObject.Poststop,
			HookObject{
				Path:    h.Path,
				Args:    h.Args,
				Env:     h.Env,
				Timeout: h.Timeout,
			},
		)
	}

	return hookLifeCycleObject
}

func buildAnnotationSpec(opts ConfigOptions) AnnotationObject {
	netSpec, _ := utils.JsonToString(buildNetSpec(opts))
	imageSpec, _ := utils.JsonToString(buildImageSpec(opts))
	return AnnotationObject{
		Version: oci.AnnotationVersion,
		Net:     netSpec,
		Image:   imageSpec,
	}
}

func buildSpec(opts ConfigOptions) Spec {
	ociVersion := oci.OCIVersion

	// root path
	root := buildRootSpec(opts)

	// mounts
	mounts := buildMountSpec(opts)

	// process
	process := buildProcessSpec(opts)

	// hostname
	hostname := opts.Hostname

	// linux spec
	linuxSpec := buildLinuxSpec(opts)

	// hook spec
	hookSpec := buildHookSpec(opts)

	// annotation
	annotation := buildAnnotationSpec(opts)

	return Spec{
		OciVersion:  ociVersion,
		Root:        root,
		Mounts:      mounts,
		Process:     process,
		Hostname:    hostname,
		LinuxSpec:   linuxSpec,
		Hooks:       hookSpec,
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
