package spec

type RootObject struct {
	Path string `json:"path"`
}

type MountObject struct {
	Destination string   `json:"destination"`
	Type        string   `json:"type"`
	Source      string   `json:"source"`
	Options     []string `json:"options"`
}

type CapabilityObject struct {
	Bounding    []string `json:"bounding"`
	Permitted   []string `json:"permitted"`
	Inheritable []string `json:"inheritable"`
	Effective   []string `json:"effective"`
	Ambient     []string `json:"ambient"`
}

type ProcessObject struct {
	Cwd          string           `json:"cwd"`
	Env          []string         `json:"env"`
	Args         []string         `json:"args"`
	Capabilities CapabilityObject `json:"capabilities"`
}

type MemoryObject struct {
	Limit int `json:"limit"`
}

type CpuObject struct {
	Period int `json:"period"`
	Quota  int `json:"quota"`
}

type ResourceObject struct {
	Memory MemoryObject `json:"memory"`
	Cpu    CpuObject    `json:"cpu"`
}

type LinuxSpecObject struct {
	Resources ResourceObject `json:"resources"`
}

type AnnotationObject struct {
	Version string `json:"io.raind.runtime.annotation.version"`
	Net     string `json:"io.raind.net.config"`
	Image   string `json:"io.raind.image.config"`
}

type Spec struct {
	OciVersion  string           `json:"ociVersion"`
	Root        RootObject       `json:"root"`
	Mounts      []MountObject    `json:"mounts"`
	Process     ProcessObject    `json:"process"`
	Hostname    string           `json:"hostname"`
	LinuxSpec   LinuxSpecObject  `json:"linux"`
	Annotations AnnotationObject `json:"annotations"`
}

// Annotation: io.raind.net.config
type IPv4Object struct {
	Address string `json:"address"`
	Gateway string `json:"gateway"`
}

type DnsObject struct {
	Servers []string `json:"servers"`
}

type InterfaceObject struct {
	Name string     `json:"name"`
	IPv4 IPv4Object `json:"ipv4"`
	Dns  DnsObject  `json:"dns"`
}

type NetConfigObject struct {
	DefaultInterface string          `json:"defaultInterface"`
	Interface        InterfaceObject `json:"interface"`
}

// Annotation: io.raind.image.config
type ImageConfigObject struct {
	RootfsType string `json:"rootfsType"`
	ImageLayer string `json:"imageLayer"`
	UpperDir   string `json:"upperDir"`
	WorkDir    string `json:"workDir"`
	MergeDir   string `json:"mergeDir"`
}
