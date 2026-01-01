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
	InterfaceName string
	Address       string
	Gateway       string
	Dns           []string
}

type ImageOption struct {
	ImageLayer string
	UpperDir   string
	WorkDir    string
	MergeDir   string
}

type ConfigOptions struct {
	Rootfs    string
	Mounts    []MountOption
	Process   ProcessOption
	Namespace []string
	Hostname  string
	Net       NetOption
	Image     ImageOption
}
