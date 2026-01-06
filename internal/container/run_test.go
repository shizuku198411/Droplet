package container

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewContainerRun_Success(t *testing.T) {
	// == arrange ==

	// == act ==
	containerRun := NewContainerRun()

	// == assert ==
	assert.NotNil(t, containerRun)
	assert.NotNil(t, containerRun.specLoader)
	assert.NotNil(t, containerRun.fifoCreator)
	assert.NotNil(t, containerRun.commandFactory)
	assert.NotNil(t, containerRun.containerStart)
}

func TestContainerRun_Run_InteractiveSuccess(t *testing.T) {
	// == arrange ==
	opts := RunOption{
		ContainerId: "12345",
		Interactive: true,
	}
	mockFileSpecLoader := &mockFileSpecLoader{
		loadFileSpec: buildMockSpec(t),
	}
	mockCotainerFifoHandler := &mockCotainerFifoHandler{}
	mockExecCmd := &mockExecCmd{}
	mockExecCommandFactory := &mockExecCommandFactory{
		commandExecutor: mockExecCmd,
	}
	mockContainerStart := &ContainerStart{
		fifoHandler: mockCotainerFifoHandler,
	}
	mockContainerNetworkController := &mockContainerNetworkController{}
	mockeContainerCgroupController := &mockeContainerCgroupController{}
	containerRun := &ContainerRun{
		specLoader:               mockFileSpecLoader,
		fifoCreator:              mockCotainerFifoHandler,
		commandFactory:           mockExecCommandFactory,
		containerStart:           mockContainerStart,
		containerNetworkPreparer: mockContainerNetworkController,
		containerCgroupPreparer:  mockeContainerCgroupController,
	}

	// == act ==
	err := containerRun.Run(opts)

	// == assert ==
	// loadFile() is called
	assert.True(t, mockFileSpecLoader.loadFileCallFlag)
	// createFifo() is called
	assert.True(t, mockCotainerFifoHandler.createFifoCallFlag)
	// Command() args is "init <container-id> /etc/raind/container/12345/exec.fifo /bin/sh"
	assert.Equal(t, []string{"init", "12345", "/etc/raind/container/12345/exec.fifo", "/bin/sh"}, mockExecCommandFactory.commandCalls[0].args)
	// SetStdout() is called
	assert.True(t, mockExecCmd.setStdoutCallFlag)
	// SetStderr() is called
	assert.True(t, mockExecCmd.setStderrCallFlag)
	// SetStdin() is called
	assert.True(t, mockExecCmd.setStdinCallFlag)
	// SetSysProcAttr() is called
	assert.True(t, mockExecCmd.setSysProcAttrCallFlag)
	// Start() is called
	assert.True(t, mockExecCmd.startCallFlag)
	// Wait() is called
	assert.True(t, mockExecCmd.waitCallFlag)
	// error is nil
	assert.Nil(t, err)
}

func TestContainerRun_Run_NonInteractiveSuccess(t *testing.T) {
	// == arrange ==
	opts := RunOption{
		ContainerId: "12345",
		Interactive: false,
	}
	mockFileSpecLoader := &mockFileSpecLoader{
		loadFileSpec: buildMockSpec(t),
	}
	mockCotainerFifoHandler := &mockCotainerFifoHandler{}
	mockExecCmd := &mockExecCmd{}
	mockExecCommandFactory := &mockExecCommandFactory{
		commandExecutor: mockExecCmd,
	}
	mockContainerStart := &ContainerStart{
		fifoHandler: mockCotainerFifoHandler,
	}
	mockContainerNetworkController := &mockContainerNetworkController{}
	mockeContainerCgroupController := &mockeContainerCgroupController{}
	containerRun := &ContainerRun{
		specLoader:               mockFileSpecLoader,
		fifoCreator:              mockCotainerFifoHandler,
		commandFactory:           mockExecCommandFactory,
		containerStart:           mockContainerStart,
		containerNetworkPreparer: mockContainerNetworkController,
		containerCgroupPreparer:  mockeContainerCgroupController,
	}

	// == act ==
	err := containerRun.Run(opts)

	// == assert ==
	// loadFile() is called
	assert.True(t, mockFileSpecLoader.loadFileCallFlag)
	// createFifo() is called
	assert.True(t, mockCotainerFifoHandler.createFifoCallFlag)
	// Command() args is "init <container-id> /etc/raind/container/12345/exec.fifo /bin/sh"
	assert.Equal(t, []string{"init", "12345", "/etc/raind/container/12345/exec.fifo", "/bin/sh"}, mockExecCommandFactory.commandCalls[0].args)
	// SetSysProcAttr() is called
	assert.True(t, mockExecCmd.setSysProcAttrCallFlag)
	// Start() is called
	assert.True(t, mockExecCmd.startCallFlag)
	// error is nil
	assert.Nil(t, err)
}

func TestContainerRun_Run_LoadFileError(t *testing.T) {
	// == arrange ==
	opts := RunOption{
		ContainerId: "12345",
		Interactive: true,
	}
	mockFileSpecLoader := &mockFileSpecLoader{
		loadFileErr: errors.New("loadFile() failed"),
	}
	mockCotainerFifoHandler := &mockCotainerFifoHandler{}
	mockExecCmd := &mockExecCmd{}
	mockExecCommandFactory := &mockExecCommandFactory{
		commandExecutor: mockExecCmd,
	}
	mockContainerStart := &ContainerStart{
		fifoHandler: mockCotainerFifoHandler,
	}
	mockContainerNetworkController := &mockContainerNetworkController{}
	mockeContainerCgroupController := &mockeContainerCgroupController{}
	containerRun := &ContainerRun{
		specLoader:               mockFileSpecLoader,
		fifoCreator:              mockCotainerFifoHandler,
		commandFactory:           mockExecCommandFactory,
		containerStart:           mockContainerStart,
		containerNetworkPreparer: mockContainerNetworkController,
		containerCgroupPreparer:  mockeContainerCgroupController,
	}

	// == act ==
	err := containerRun.Run(opts)

	// == assert ==
	assert.NotNil(t, err)
	assert.Equal(t, errors.New("loadFile() failed"), err)
}

func TestContainerRun_Run_CreateFifoError(t *testing.T) {
	// == arrange ==
	opts := RunOption{
		ContainerId: "12345",
		Interactive: true,
	}
	mockFileSpecLoader := &mockFileSpecLoader{
		loadFileSpec: buildMockSpec(t),
	}
	mockCotainerFifoHandler := &mockCotainerFifoHandler{
		createFifoErr: errors.New("createFifo() failed"),
	}
	mockExecCmd := &mockExecCmd{}
	mockExecCommandFactory := &mockExecCommandFactory{
		commandExecutor: mockExecCmd,
	}
	mockContainerStart := &ContainerStart{
		fifoHandler: mockCotainerFifoHandler,
	}
	mockContainerNetworkController := &mockContainerNetworkController{}
	mockeContainerCgroupController := &mockeContainerCgroupController{}
	containerRun := &ContainerRun{
		specLoader:               mockFileSpecLoader,
		fifoCreator:              mockCotainerFifoHandler,
		commandFactory:           mockExecCommandFactory,
		containerStart:           mockContainerStart,
		containerNetworkPreparer: mockContainerNetworkController,
		containerCgroupPreparer:  mockeContainerCgroupController,
	}

	// == act ==
	err := containerRun.Run(opts)

	// == assert ==
	assert.NotNil(t, err)
	assert.Equal(t, errors.New("createFifo() failed"), err)
}

func TestContainerRun_Run_CgroupPrepareError(t *testing.T) {
	// == arrange ==
	opts := RunOption{
		ContainerId: "12345",
		Interactive: true,
	}
	mockFileSpecLoader := &mockFileSpecLoader{
		loadFileSpec: buildMockSpec(t),
	}
	mockCotainerFifoHandler := &mockCotainerFifoHandler{}
	mockExecCmd := &mockExecCmd{}
	mockExecCommandFactory := &mockExecCommandFactory{
		commandExecutor: mockExecCmd,
	}
	mockContainerStart := &ContainerStart{
		fifoHandler: mockCotainerFifoHandler,
	}
	mockContainerNetworkController := &mockContainerNetworkController{}
	mockeContainerCgroupController := &mockeContainerCgroupController{
		prepareErr: errors.New("prepare() failed"),
	}
	containerRun := &ContainerRun{
		specLoader:               mockFileSpecLoader,
		fifoCreator:              mockCotainerFifoHandler,
		commandFactory:           mockExecCommandFactory,
		containerStart:           mockContainerStart,
		containerNetworkPreparer: mockContainerNetworkController,
		containerCgroupPreparer:  mockeContainerCgroupController,
	}

	// == act ==
	err := containerRun.Run(opts)

	// == assert ==
	// error is not nil
	assert.NotNil(t, err)
	assert.Equal(t, errors.New("prepare() failed"), err)
}

func TestContainerRun_Run_NetworkPrepareError(t *testing.T) {
	// == arrange ==
	opts := RunOption{
		ContainerId: "12345",
		Interactive: true,
	}
	mockFileSpecLoader := &mockFileSpecLoader{
		loadFileSpec: buildMockSpec(t),
	}
	mockCotainerFifoHandler := &mockCotainerFifoHandler{}
	mockExecCmd := &mockExecCmd{}
	mockExecCommandFactory := &mockExecCommandFactory{
		commandExecutor: mockExecCmd,
	}
	mockContainerStart := &ContainerStart{
		fifoHandler: mockCotainerFifoHandler,
	}
	mockContainerNetworkController := &mockContainerNetworkController{
		prepareErr: errors.New("prepare() failed"),
	}
	mockeContainerCgroupController := &mockeContainerCgroupController{}
	containerRun := &ContainerRun{
		specLoader:               mockFileSpecLoader,
		fifoCreator:              mockCotainerFifoHandler,
		commandFactory:           mockExecCommandFactory,
		containerStart:           mockContainerStart,
		containerNetworkPreparer: mockContainerNetworkController,
		containerCgroupPreparer:  mockeContainerCgroupController,
	}

	// == act ==
	err := containerRun.Run(opts)

	// == assert ==
	// error is not nil
	assert.NotNil(t, err)
	assert.Equal(t, errors.New("prepare() failed"), err)
}

func TestContainerRun_Run_StartError(t *testing.T) {
	// == arrange ==
	opts := RunOption{
		ContainerId: "12345",
		Interactive: true,
	}
	mockFileSpecLoader := &mockFileSpecLoader{
		loadFileSpec: buildMockSpec(t),
	}
	mockCotainerFifoHandler := &mockCotainerFifoHandler{}
	mockExecCmd := &mockExecCmd{
		startErr: errors.New("Start() failed"),
	}
	mockExecCommandFactory := &mockExecCommandFactory{
		commandExecutor: mockExecCmd,
	}
	mockContainerStart := &ContainerStart{
		fifoHandler: mockCotainerFifoHandler,
	}
	mockContainerNetworkController := &mockContainerNetworkController{}
	mockeContainerCgroupController := &mockeContainerCgroupController{}
	containerRun := &ContainerRun{
		specLoader:               mockFileSpecLoader,
		fifoCreator:              mockCotainerFifoHandler,
		commandFactory:           mockExecCommandFactory,
		containerStart:           mockContainerStart,
		containerNetworkPreparer: mockContainerNetworkController,
		containerCgroupPreparer:  mockeContainerCgroupController,
	}

	// == act ==
	err := containerRun.Run(opts)

	// == assert ==
	assert.NotNil(t, err)
	assert.Equal(t, errors.New("Start() failed"), err)
}

func TestContainerRun_Run_WaitError(t *testing.T) {
	// == arrange ==
	opts := RunOption{
		ContainerId: "12345",
		Interactive: true,
	}
	mockFileSpecLoader := &mockFileSpecLoader{
		loadFileSpec: buildMockSpec(t),
	}
	mockCotainerFifoHandler := &mockCotainerFifoHandler{}
	mockExecCmd := &mockExecCmd{
		waitErr: errors.New("Wait() failed"),
	}
	mockExecCommandFactory := &mockExecCommandFactory{
		commandExecutor: mockExecCmd,
	}
	mockContainerStart := &ContainerStart{
		fifoHandler: mockCotainerFifoHandler,
	}
	mockContainerNetworkController := &mockContainerNetworkController{}
	mockeContainerCgroupController := &mockeContainerCgroupController{}
	containerRun := &ContainerRun{
		specLoader:               mockFileSpecLoader,
		fifoCreator:              mockCotainerFifoHandler,
		commandFactory:           mockExecCommandFactory,
		containerStart:           mockContainerStart,
		containerNetworkPreparer: mockContainerNetworkController,
		containerCgroupPreparer:  mockeContainerCgroupController,
	}

	// == act ==
	err := containerRun.Run(opts)

	// == assert ==
	assert.NotNil(t, err)
	assert.Equal(t, errors.New("Wait() failed"), err)
}
