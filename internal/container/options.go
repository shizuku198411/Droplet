package container

// create options
type CreateOption struct {
	ContainerId  string
	PrintPidFlag bool
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
	ContainerId  string
	Interactive  bool
	PrintPidFlag bool
}
