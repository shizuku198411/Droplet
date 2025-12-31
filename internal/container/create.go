package container

import (
	"droplet/internal/spec"
	"fmt"
	"os"
	"os/exec"
)

func CreateContainer(opt CreateOption) error {
	// load config.json
	spec, err := spec.LoadConfigFile(configFilePath(opt.ContainerId))
	if err != nil {
		return err
	}

	// create fifo
	fifo := fifoPath(opt.ContainerId)
	if err := createFifo(fifo); err != nil {
		return err
	}

	// execute init subcommand
	initPid, err := ExecuteInit(spec, fifo)
	if err != nil {
		return err
	}

	fmt.Printf("init process has been created. pid: %d\n", initPid)

	return nil
}

func ExecuteInit(spec spec.Spec, fifo string) (int, error) {
	// retrieve entrypoint from spec
	entrypoint := spec.Process.Args

	// prepare init subcommand
	initArgs := append([]string{"init", fifo}, entrypoint...)
	cmd := exec.Command(os.Args[0], initArgs...)
	// set stdout/stderr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// execute init subcommand
	if err := cmd.Start(); err != nil {
		return -1, err
	}

	return cmd.Process.Pid, nil
}
