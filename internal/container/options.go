package container

// create options
type CreateOption struct {
	ContainerId string
}

// init options
type InitOption struct {
	ContainerId string
	Fifo        string
	Entrypoint  []string
}

// start options
type StartOption struct {
	ContainerId string
}

// run options
type RunOption struct {
	ContainerId string
	Interactive bool
}
