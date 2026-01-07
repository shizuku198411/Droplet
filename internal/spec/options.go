package spec

type MountOption struct {
	Destination string
	Type        string
	Source      string
	Options     []string
}

type ProcessOption struct {
	Cwd  string
	Env  []string
	Args []string
}

type NetOption struct {
	HostInterface       string
	BridgeInterfaceName string
	InterfaceName       string
	Address             string
	Gateway             string
	Dns                 []string
}

type ImageOption struct {
	ImageLayer []string
	UpperDir   string
	WorkDir    string
}

type HookOption struct {
	Path    string
	Args    []string
	Env     []string
	Timeout *int
}

type HookLifecycleOption struct {
	Prestart        []HookOption
	CreateRuntime   []HookOption
	CreateContainer []HookOption
	StartContainer  []HookOption
	Poststart       []HookOption
	Poststop        []HookOption
}

type ConfigOptions struct {
	Rootfs    string
	Mounts    []MountOption
	Process   ProcessOption
	Namespace []string
	Hostname  string
	Net       NetOption
	Image     ImageOption
	Hooks     HookLifecycleOption
}
