package container

import (
	"fmt"
	"testing"

	"droplet/internal/testutils"

	"github.com/stretchr/testify/assert"
)

func TestCreate_Success(t *testing.T) {
	dummyContainerCreator := &ContainerCreator{
		specLoader:      &dummyFileSpecLoader{spec: dummySpec()},
		fifoCreator:     &dummyFifoHandler{},
		processExecutor: &dummyContainerInitExecutor{Pid: 11111},
	}

	result := testutils.CaptureStdout(t, func() {
		err := dummyContainerCreator.Create(CreateOption{ContainerId: "123456"})
		if err != nil {
			t.Fatal(err)
		}
	})

	expect := "init process has been created. pid: 11111\n"

	assert.Equal(t, expect, result)
}

func TestCreate_LoadConfigError(t *testing.T) {
	dummyContainerInitExecutor := &ContainerCreator{
		specLoader: &dummyFileSpecLoader{err: fmt.Errorf("failed to load config.json")},
	}

	result := dummyContainerInitExecutor.Create(CreateOption{ContainerId: "123456"})

	expect := fmt.Errorf("failed to load config.json")

	assert.Equal(t, expect, result)
}

func TestCreate_CreateFifoError(t *testing.T) {
	dummyContainerInitExecutor := &ContainerCreator{
		specLoader:  &dummyFileSpecLoader{spec: dummySpec()},
		fifoCreator: &dummyFifoHandler{createErr: fmt.Errorf("failed to create FIFO")},
	}

	result := dummyContainerInitExecutor.Create(CreateOption{ContainerId: "123456"})

	expect := fmt.Errorf("failed to create FIFO")

	assert.Equal(t, expect, result)
}

func TestCreate_ExecuteInitError(t *testing.T) {
	dummyContainerInitExecutor := &ContainerCreator{
		specLoader:      &dummyFileSpecLoader{spec: dummySpec()},
		fifoCreator:     &dummyFifoHandler{},
		processExecutor: &dummyContainerInitExecutor{Err: fmt.Errorf("failed to execute process")},
	}

	result := dummyContainerInitExecutor.Create(CreateOption{ContainerId: "123456"})

	expect := fmt.Errorf("failed to execute process")

	assert.Equal(t, expect, result)
}

func TestExecuteInit_Success(t *testing.T) {
	spec := dummySpec()
	fifo := "/tmp/exec.fifo"

	dummyCmd := &dummyCmd{
		pid: 11111,
		err: nil,
	}
	dummyCmdFactory := &dummyCommandFactory{
		cmd: dummyCmd,
	}

	dummyInitExecutor := &containerInitExecutor{
		commandFactory: dummyCmdFactory,
	}

	pid, err := dummyInitExecutor.executeInit(spec, fifo)
	if err != nil {
		t.Fatalf("executeInit returned error: %v", err)
	}

	// assert
	// 1. the args is set to "init <fifo-path> <entrypoint>"
	expectArgs := []string{"init", "/tmp/exec.fifo", "/bin/sh"}
	resultArgs := dummyCmdFactory.commandArgs
	assert.Equal(t, expectArgs, resultArgs)

	// 2. Start() is being called
	expectStartFlag := true
	resultStartFlag := dummyCmd.startFlag
	assert.Equal(t, expectStartFlag, resultStartFlag)

	// 3. PID is returned
	expectPid := 11111
	resultPid := pid
	assert.Equal(t, expectPid, resultPid)
}

func TestExecuteInit_StartError(t *testing.T) {
	spec := dummySpec()
	fifo := "/tmp/exec.fifo"

	dummyCmd := &dummyCmd{
		pid: 11111,
		err: fmt.Errorf("start error"),
	}
	dummyCmdFactory := &dummyCommandFactory{
		cmd: dummyCmd,
	}

	dummyInitExecutor := &containerInitExecutor{
		commandFactory: dummyCmdFactory,
	}

	pid, err := dummyInitExecutor.executeInit(spec, fifo)

	// assert
	// 1. pid is -1
	expectPid := -1
	resultPid := pid
	assert.Equal(t, expectPid, resultPid)

	// 2. error message is returned
	expectErr := fmt.Errorf("start error")
	resultErr := err
	assert.Equal(t, expectErr, resultErr)
}
