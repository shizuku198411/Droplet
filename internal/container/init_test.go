package container

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func buildInitOption(t *testing.T) InitOption {
	t.Helper()

	return InitOption{
		ContainerId: "12345",
		Fifo:        "exec.fifo",
		Entrypoint:  []string{"/bin/sh"},
	}
}

func TestNewContainerInit_Success(t *testing.T) {
	// == arrange ==

	// == act ==
	containerInit := NewContainerInit()

	// == assert ==
	assert.NotNil(t, containerInit)
	assert.NotNil(t, containerInit.fifoReader)
	assert.NotNil(t, containerInit.specLoader)
	assert.NotNil(t, containerInit.containerEnvPreparer)
	assert.NotNil(t, containerInit.syscallHandler)
}

func TestNewRootContainerEnvPrepare_Success(t *testing.T) {
	// == arrange ==

	// == act ==
	rootContainerEnvPreparer := newRootContainerEnvPrepare()

	// == assert ==
	assert.NotNil(t, rootContainerEnvPreparer)
	assert.NotNil(t, rootContainerEnvPreparer.syscallHandler)
}

func TestContainerInit_Execute_Success(t *testing.T) {
	// == arrange ==
	opts := buildInitOption(t)
	mockFifoReader := &mockCotainerFifoHandler{}
	mockSpecLoader := &mockFileSpecLoader{}
	mockRootContainerEnvPreparer := &mockRootContainerEnvPreparer{}
	mockKernelSyscall := &mockKernelSyscall{}
	containerInit := &ContainerInit{
		fifoReader:           mockFifoReader,
		specLoader:           mockSpecLoader,
		containerEnvPreparer: mockRootContainerEnvPreparer,
		syscallHandler:       mockKernelSyscall,
	}

	// == act ==
	err := containerInit.Execute(opts)

	// == assert ==
	// readFifo() is called
	assert.True(t, mockFifoReader.readFifoCallFlag)
	// loadFile() is called
	assert.True(t, mockSpecLoader.loadFileCallFlag)
	// prepare() is called
	assert.True(t, mockRootContainerEnvPreparer.prepareCallFlag)
	// Exec() is called
	assert.True(t, mockKernelSyscall.execCallFlag)
	// error is nil
	assert.Nil(t, err)
}

func TestContainerInit_Execute_ReadFifoError(t *testing.T) {
	// == arrange ==
	opts := buildInitOption(t)
	mockFifoReader := &mockCotainerFifoHandler{
		readFifoErr: errors.New("readFifo() failed"),
	}
	mockSpecLoader := &mockFileSpecLoader{}
	mockRootContainerEnvPreparer := &mockRootContainerEnvPreparer{}
	mockKernelSyscall := &mockKernelSyscall{}
	containerInit := &ContainerInit{
		fifoReader:           mockFifoReader,
		specLoader:           mockSpecLoader,
		containerEnvPreparer: mockRootContainerEnvPreparer,
		syscallHandler:       mockKernelSyscall,
	}

	// == act ==
	err := containerInit.Execute(opts)

	// == assert ==
	assert.NotNil(t, err)
	assert.Equal(t, errors.New("readFifo() failed"), err)
}

func TestContainerInit_Execute_LoadFileError(t *testing.T) {
	// == arrange ==
	opts := buildInitOption(t)
	mockFifoReader := &mockCotainerFifoHandler{}
	mockSpecLoader := &mockFileSpecLoader{
		loadFileErr: errors.New("loadFile() failed"),
	}
	mockRootContainerEnvPreparer := &mockRootContainerEnvPreparer{}
	mockKernelSyscall := &mockKernelSyscall{}
	containerInit := &ContainerInit{
		fifoReader:           mockFifoReader,
		specLoader:           mockSpecLoader,
		containerEnvPreparer: mockRootContainerEnvPreparer,
		syscallHandler:       mockKernelSyscall,
	}

	// == act ==
	err := containerInit.Execute(opts)

	// == assert ==
	assert.NotNil(t, err)
	assert.Equal(t, errors.New("loadFile() failed"), err)
}

func TestContainerInit_Execute_PrepareError(t *testing.T) {
	// == arrange ==
	opts := buildInitOption(t)
	mockFifoReader := &mockCotainerFifoHandler{}
	mockSpecLoader := &mockFileSpecLoader{}
	mockRootContainerEnvPreparer := &mockRootContainerEnvPreparer{
		prepareErr: errors.New("prepare() failed"),
	}
	mockKernelSyscall := &mockKernelSyscall{}
	containerInit := &ContainerInit{
		fifoReader:           mockFifoReader,
		specLoader:           mockSpecLoader,
		containerEnvPreparer: mockRootContainerEnvPreparer,
		syscallHandler:       mockKernelSyscall,
	}

	// == act ==
	err := containerInit.Execute(opts)

	// == assert ==
	assert.NotNil(t, err)
	assert.Equal(t, errors.New("prepare() failed"), err)
}

func TestContainerInit_Execute_ExecError(t *testing.T) {
	// == arrange ==
	opts := buildInitOption(t)
	mockFifoReader := &mockCotainerFifoHandler{}
	mockSpecLoader := &mockFileSpecLoader{}
	mockRootContainerEnvPreparer := &mockRootContainerEnvPreparer{}
	mockKernelSyscall := &mockKernelSyscall{
		execErr: errors.New("Exec() failed"),
	}
	containerInit := &ContainerInit{
		fifoReader:           mockFifoReader,
		specLoader:           mockSpecLoader,
		containerEnvPreparer: mockRootContainerEnvPreparer,
		syscallHandler:       mockKernelSyscall,
	}

	// == act ==
	err := containerInit.Execute(opts)

	// == assert ==
	assert.NotNil(t, err)
	assert.Equal(t, errors.New("Exec() failed"), err)
}

func TestRootContainerEnvPrepare_Success(t *testing.T) {
	// == arrange ==
	spec := buildMockSpec(t)
	mockKernelSyscall := &mockKernelSyscall{}
	rootContainerEnvPreparer := &rootContainerEnvPreparer{
		syscallHandler: mockKernelSyscall,
	}

	// == act ==
	err := rootContainerEnvPreparer.prepare(spec)

	// == assert ==
	// Setresgid() is called from switchToUserNamespaceRoot()
	assert.True(t, mockKernelSyscall.setresgidCallFlag)
	// Setresuid() is called from switchToUserNamespaceRoot()
	assert.True(t, mockKernelSyscall.setresuidCallFlag)
	// Sethostname() is called from setHostnameToContainerId
	assert.True(t, mockKernelSyscall.sethostnameCallFlag)
	// Sethostname() is recieved hostname from config.json
	assert.Equal(t, mockKernelSyscall.sethostnameP, []byte(spec.Hostname))
	// error is nil
	assert.Nil(t, err)
}

func TestRootContainerEnvPrepare_SetresgidError(t *testing.T) {
	// == arrange ==
	spec := buildMockSpec(t)
	mockKernelSyscall := &mockKernelSyscall{
		setresgidErr: errors.New("Setresgid() failed"),
	}
	rootContainerEnvPreparer := &rootContainerEnvPreparer{
		syscallHandler: mockKernelSyscall,
	}

	// == act ==
	err := rootContainerEnvPreparer.prepare(spec)

	// == assert ==
	assert.NotNil(t, err)
	assert.Equal(t, errors.New("Setresgid() failed"), err)
}

func TestRootContainerEnvPrepare_SetresuidError(t *testing.T) {
	// == arrange ==
	spec := buildMockSpec(t)
	mockKernelSyscall := &mockKernelSyscall{
		setresuidErr: errors.New("Setresuid() failed"),
	}
	rootContainerEnvPreparer := &rootContainerEnvPreparer{
		syscallHandler: mockKernelSyscall,
	}

	// == act ==
	err := rootContainerEnvPreparer.prepare(spec)

	// == assert ==
	assert.NotNil(t, err)
	assert.Equal(t, errors.New("Setresuid() failed"), err)
}

func TestRootContainerEnvPrepare_SethostnameError(t *testing.T) {
	// == arrange ==
	spec := buildMockSpec(t)
	mockKernelSyscall := &mockKernelSyscall{
		sethostnameErr: errors.New("Sethostname() failed"),
	}
	rootContainerEnvPreparer := &rootContainerEnvPreparer{
		syscallHandler: mockKernelSyscall,
	}

	// == act ==
	err := rootContainerEnvPreparer.prepare(spec)

	// == assert ==
	assert.NotNil(t, err)
	assert.Equal(t, errors.New("Sethostname() failed"), err)
}
