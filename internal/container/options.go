package container

// create options
type CreateOption struct {
	ContainerId string
}

// init options
type InitOption struct {
	Fifo       string
	Entrypoint []string
}

// start options
type StartOption struct {
	ContainerId string
}
