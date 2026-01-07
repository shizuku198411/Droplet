package container

import (
	"droplet/internal/status"
	"errors"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewContainerKill_Success(t *testing.T) {
	// == arrange ==

	// == act ==
	containerKill := NewContainerKill()

	// == assert ==
	assert.NotNil(t, containerKill)
	assert.NotNil(t, containerKill.syscallHandler)
	assert.NotNil(t, containerKill.containerStatusManager)
}

func TestContainerKill_Kill_Success(t *testing.T) {
	// == arrange ==
	opt := KillOption{
		ContainerId: "12345",
		Signal:      "KILL",
	}
	mockKernelSyscall := &mockKernelSyscall{}
	mockContainerStatusManager := &mockStatusHandler{
		getStatusFromIdStatus: status.RUNNING,
		getStatusFromIdErr:    nil,
		getPidFromIdPid:       11111,
	}
	containerKill := &ContainerKill{
		syscallHandler:         mockKernelSyscall,
		containerStatusManager: mockContainerStatusManager,
	}

	// == act ==
	err := containerKill.Kill(opt)

	// == assert ==
	// Kill() is called
	assert.True(t, mockKernelSyscall.killCallFlag)
	// Kill() parameter: pid=11111, signal=KILL
	assert.Equal(t, 11111, mockKernelSyscall.killPid)
	assert.Equal(t, syscall.SIGKILL, mockKernelSyscall.killSig)

	// error is nil
	assert.Nil(t, err)
}

func TestContainerKill_Kill_GetStatusError(t *testing.T) {
	// == arrange ==
	opt := KillOption{
		ContainerId: "12345",
		Signal:      "KILL",
	}
	mockKernelSyscall := &mockKernelSyscall{}
	mockContainerStatusManager := &mockStatusHandler{
		getStatusFromIdErr: errors.New("get status failed"),
	}
	containerKill := &ContainerKill{
		syscallHandler:         mockKernelSyscall,
		containerStatusManager: mockContainerStatusManager,
	}

	// == act ==
	err := containerKill.Kill(opt)

	// == assert ==
	assert.NotNil(t, err)
	assert.Equal(t, errors.New("get status failed"), err)
}

func TestContainerKill_Kill_NotRunningError(t *testing.T) {
	// == arrange ==
	opt := KillOption{
		ContainerId: "12345",
		Signal:      "KILL",
	}
	mockKernelSyscall := &mockKernelSyscall{}
	mockContainerStatusManager := &mockStatusHandler{
		getStatusFromIdStatus: status.STOPPED,
		getStatusFromIdErr:    nil,
	}
	containerKill := &ContainerKill{
		syscallHandler:         mockKernelSyscall,
		containerStatusManager: mockContainerStatusManager,
	}

	// == act ==
	err := containerKill.Kill(opt)

	// == assert ==
	assert.NotNil(t, err)
	assert.Equal(t, errors.New("container: 12345 not running."), err)
}

func TestContainerKill_Kill_GetPidError(t *testing.T) {
	// == arrange ==
	opt := KillOption{
		ContainerId: "12345",
		Signal:      "KILL",
	}
	mockKernelSyscall := &mockKernelSyscall{}
	mockContainerStatusManager := &mockStatusHandler{
		getStatusFromIdStatus: status.RUNNING,
		getStatusFromIdErr:    nil,
		getPidFromIdErr:       errors.New("get pid failed"),
	}
	containerKill := &ContainerKill{
		syscallHandler:         mockKernelSyscall,
		containerStatusManager: mockContainerStatusManager,
	}

	// == act ==
	err := containerKill.Kill(opt)

	// == assert ==
	assert.NotNil(t, err)
	assert.Equal(t, errors.New("get pid failed"), err)
}
