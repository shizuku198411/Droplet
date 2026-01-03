package container

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func buildCreateOption(t *testing.T) CreateOption {
	t.Helper()

	return CreateOption{
		ContainerId: "123456",
	}
}

func TestNewContainerCreator_Success(t *testing.T) {
	// == arrange ==

	// == act ==
	containerCreator := NewContainerCreator()

	// == assert ==
	assert.NotNil(t, containerCreator)
	assert.NotNil(t, containerCreator.specLoader)
	assert.NotNil(t, containerCreator.fifoCreator)
	assert.NotNil(t, containerCreator.processExecutor)
}

func TestNewContainerInitExecutor_Success(t *testing.T) {
	// == arrange ==

	// == act ==
	containerInitExecutor := newContainerInitExecutor()

	// == assert ==
	assert.NotNil(t, containerInitExecutor)
	assert.NotNil(t, containerInitExecutor.commandFactory)
}

func TestContainerCreator_Create_Success(t *testing.T) {
	// == arrange ==
	opts := buildCreateOption(t)
	mockSpecLoader := &mockFileSpecLoader{}
	mockFifoCreator := &mockCotainerFifoHandler{}
	mockProcessExecutor := &mockContainerInitExecutor{}
	mockContainerCreator := ContainerCreator{
		specLoader:      mockSpecLoader,
		fifoCreator:     mockFifoCreator,
		processExecutor: mockProcessExecutor,
	}

	// == act ==
	err := mockContainerCreator.Create(opts)

	// == assert ==
	// loadFile() is called
	assert.True(t, mockSpecLoader.loadFileCallFlag)

	// createFifo() is called
	assert.True(t, mockFifoCreator.createFifoCallFlag)

	// executeInit() is called
	assert.True(t, mockProcessExecutor.executeInitCallFlag)

	// err is nil
	assert.Nil(t, err)
}

func TestContainerCreator_Create_LoadFileError(t *testing.T) {
	// == arrange ==
	opts := buildCreateOption(t)
	mockSpecLoader := &mockFileSpecLoader{
		loadFileErr: errors.New("loadFile() failed"),
	}
	mockFifoCreator := &mockCotainerFifoHandler{}
	mockProcessExecutor := &mockContainerInitExecutor{}
	mockContainerCreator := ContainerCreator{
		specLoader:      mockSpecLoader,
		fifoCreator:     mockFifoCreator,
		processExecutor: mockProcessExecutor,
	}

	// == act ==
	err := mockContainerCreator.Create(opts)

	// == assert ==
	assert.NotNil(t, err)
	assert.Equal(t, errors.New("loadFile() failed"), err)
}

func TestContainerCreator_Create_CreateFifoError(t *testing.T) {
	// == arrange ==
	opts := buildCreateOption(t)
	mockSpecLoader := &mockFileSpecLoader{}
	mockFifoCreator := &mockCotainerFifoHandler{
		createFifoErr: errors.New("createFifo() failed"),
	}
	mockProcessExecutor := &mockContainerInitExecutor{}
	mockContainerCreator := ContainerCreator{
		specLoader:      mockSpecLoader,
		fifoCreator:     mockFifoCreator,
		processExecutor: mockProcessExecutor,
	}

	// == act ==
	err := mockContainerCreator.Create(opts)

	// == assert ==
	assert.NotNil(t, err)
	assert.Equal(t, errors.New("createFifo() failed"), err)
}

func TestContainerCreator_Create_ExecuteInitError(t *testing.T) {
	// == arrange ==
	opts := buildCreateOption(t)
	mockSpecLoader := &mockFileSpecLoader{}
	mockFifoCreator := &mockCotainerFifoHandler{}
	mockProcessExecutor := &mockContainerInitExecutor{
		executeInitErr: errors.New("executeInit() failed"),
	}
	mockContainerCreator := ContainerCreator{
		specLoader:      mockSpecLoader,
		fifoCreator:     mockFifoCreator,
		processExecutor: mockProcessExecutor,
	}

	// == act ==
	err := mockContainerCreator.Create(opts)

	// == assert ==
	assert.NotNil(t, err)
	assert.Equal(t, errors.New("executeInit() failed"), err)
}

func TestContainerInitExecutor_ExecuteInit_Success(t *testing.T) {
	// == arrange ==
	mockSpec := buildMockSpec(t)
	mockExecCmd := &mockExecCmd{
		pidPid: 12345,
	}
	mockExecCommandFactory := &mockExecCommandFactory{
		commandExecutor: mockExecCmd,
	}
	mockContainerInitExecutor := &containerInitExecutor{
		commandFactory: mockExecCommandFactory,
	}
	containerId := "12345"
	fifo := "exec.fifo"

	// == act ==
	pid, err := mockContainerInitExecutor.executeInit(containerId, mockSpec, fifo)

	// == assert ==
	// Command() is called
	assert.True(t, mockExecCommandFactory.commandCallFlag)

	// command args is "init <container-id> <fifo> <entrypoint>"
	expectArgs := []string{"init", "12345", "exec.fifo", "/bin/sh"}
	assert.Equal(t, mockExecCommandFactory.commandArgs, expectArgs)

	// SetSysProcAttr() is called
	assert.True(t, mockExecCmd.setSysProcAttrCallFlag)

	// Start() is called
	assert.True(t, mockExecCmd.startCallFlag)

	// pid: 12345 is returned
	assert.Equal(t, 12345, pid)

	// error is nil
	assert.Nil(t, err)
}

func TestContainerInitExecutor_ExecuteInit_StartError(t *testing.T) {
	// == arrange ==
	mockSpec := buildMockSpec(t)
	mockExecCmd := &mockExecCmd{
		startErr: errors.New("Start() failed"),
	}
	mockExecCommandFactory := &mockExecCommandFactory{
		commandExecutor: mockExecCmd,
	}
	mockContainerInitExecutor := &containerInitExecutor{
		commandFactory: mockExecCommandFactory,
	}
	containerId := "12345"
	fifo := "exec.fifo"

	// == act ==
	pid, err := mockContainerInitExecutor.executeInit(containerId, mockSpec, fifo)

	// == assert ==
	// pid: -1 is returned
	assert.Equal(t, -1, pid)
	// error is returned
	assert.NotNil(t, err)
	assert.Equal(t, errors.New("Start() failed"), err)
}
