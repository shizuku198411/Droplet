package container

import (
	"droplet/internal/status"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func buildStartOption(t *testing.T) StartOption {
	t.Helper()

	return StartOption{
		ContainerId: "12345",
	}
}

func TestNewContainerStart_Success(t *testing.T) {
	// == arrange ==

	// == act ==
	containerStart := NewContainerStart()

	// == assert ==
	assert.NotNil(t, containerStart)
	assert.NotNil(t, containerStart.fifoHandler)
}

func TestContainerStart_Success(t *testing.T) {
	// == arrange ==
	opts := buildStartOption(t)
	mockSpecLoader := &mockFileSpecLoader{}
	mockCotainerFifoHandler := &mockCotainerFifoHandler{}
	mockContainerStatusManager := &mockStatusHandler{
		getStatusFromIdStatus: status.CREATED,
	}
	mockHookController := &mockHookController{}
	containerStart := &ContainerStart{
		specLoader:              mockSpecLoader,
		fifoHandler:             mockCotainerFifoHandler,
		containerStatusManager:  mockContainerStatusManager,
		containerHookController: mockHookController,
	}

	// == act ==
	err := containerStart.Execute(opts)

	// == assert ==
	// writeFifo() is called
	assert.True(t, mockCotainerFifoHandler.writeFifoCallFlag)
	// writeFifo() path is /etc/raind/container/<containerId>/exec.fifo
	assert.Equal(t, "/etc/raind/container/12345/exec.fifo", mockCotainerFifoHandler.writeFifoPath)
	// removeFifo() is called
	assert.True(t, mockCotainerFifoHandler.removeFifoCallFlag)
	// removeFifo() path is /etc/raind/container/<containerId>/exec.fifo
	assert.Equal(t, "/etc/raind/container/12345/exec.fifo", mockCotainerFifoHandler.removeFifoPath)
	// error is nil
	assert.Nil(t, err)
}

func TestContainerStart_WriteFifoError(t *testing.T) {
	// == arrange ==
	opts := buildStartOption(t)
	mockSpecLoader := &mockFileSpecLoader{}
	mockCotainerFifoHandler := &mockCotainerFifoHandler{
		writeFifoErr: errors.New("writeFifo() failed"),
	}
	mockContainerStatusManager := &mockStatusHandler{
		getStatusFromIdStatus: status.CREATED,
	}
	mockHookController := &mockHookController{}
	containerStart := &ContainerStart{
		specLoader:              mockSpecLoader,
		fifoHandler:             mockCotainerFifoHandler,
		containerStatusManager:  mockContainerStatusManager,
		containerHookController: mockHookController,
	}

	// == act ==
	err := containerStart.Execute(opts)

	// == assert ==
	assert.NotNil(t, err)
	assert.Equal(t, errors.New("writeFifo() failed"), err)
}

func TestContainerStart_RemoveFifoError(t *testing.T) {
	// == arrange ==
	opts := buildStartOption(t)
	mockSpecLoader := &mockFileSpecLoader{}
	mockCotainerFifoHandler := &mockCotainerFifoHandler{
		removeFifoErr: errors.New("removeFifo() failed"),
	}
	mockContainerStatusManager := &mockStatusHandler{
		getStatusFromIdStatus: status.CREATED,
	}
	mockHookController := &mockHookController{}
	containerStart := &ContainerStart{
		specLoader:              mockSpecLoader,
		fifoHandler:             mockCotainerFifoHandler,
		containerStatusManager:  mockContainerStatusManager,
		containerHookController: mockHookController,
	}

	// == act ==
	err := containerStart.Execute(opts)

	// == assert ==
	assert.NotNil(t, err)
	assert.Equal(t, errors.New("removeFifo() failed"), err)
}
