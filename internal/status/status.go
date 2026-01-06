package status

type StatusObject struct {
	OciVersion string `json:"ociVersion"`
	Id         string `json:"id"`
	Status     string `json:"status"`
	Pid        int    `json:"pid"`
	Bundle     string `json:"bundle"`
}

// container status
//
//	creating = 0
//	created  = 1
//	running  = 2
//	stopped  = 3
type ContainerStatus int

const (
	CREATING ContainerStatus = iota
	CREATED
	RUNNING
	STOPPED
)

func (s ContainerStatus) String() string {
	switch s {
	case CREATING:
		return "creating"
	case CREATED:
		return "created"
	case RUNNING:
		return "running"
	case STOPPED:
		return "stopped"
	default:
		return "unknown"
	}
}
