package container

import (
	"droplet/internal/status"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewContainerExec_Success(t *testing.T) {
	// == arrange ==

	// == act ==
	containerExec := NewContainerExec()

	// == assert ==
	assert.NotNil(t, containerExec)
	assert.NotNil(t, containerExec.commandFactory)
	assert.NotNil(t, containerExec.containerStatusManager)
}

func TestContainerExec_Exec_InteractiveSuccess(t *testing.T) {
	// == arrange ==
	opt := ExecOption{
		ContainerId: "12345",
		Tty:         true,
		Entrypoint:  []string{"/bin/sh"},
	}
	mockExecCmd := &mockExecCmd{}
	mockExecCommandFactory := &mockExecCommandFactory{
		commandExecutor: mockExecCmd,
	}
	mockContainerStatusManager := &mockStatusHandler{
		getStatusFromIdStatus: status.RUNNING,
		getStatusFromIdErr:    nil,
		getPidFromIdPid:       11111,
	}
	containerExec := &ContainerExec{
		commandFactory:         mockExecCommandFactory,
		containerStatusManager: mockContainerStatusManager,
	}

	// == act ==
	err := containerExec.Exec(opt)

	// == assert ==
	// Command() parameter: nsenter -t <pid> --all <entrypoint>
	commandCall := mockExecCommandFactory.commandCalls[0]
	assert.Equal(t, "nsenter", commandCall.name)
	assert.Equal(t, []string{"-t", "11111", "--all", "/bin/sh"}, commandCall.args)

	// SetStdout/SetStderr/SetStdin is called
	assert.True(t, mockExecCmd.setStdoutCallFlag)
	assert.True(t, mockExecCmd.setStderrCallFlag)
	assert.True(t, mockExecCmd.setStdinCallFlag)

	// Wait() is called
	assert.True(t, mockExecCmd.waitCallFlag)

	// error is nil
	assert.Nil(t, err)
}

func TestContainerExec_Exec_NonInteractiveSuccess(t *testing.T) {
	// == arrange ==
	opt := ExecOption{
		ContainerId: "12345",
		Tty:         false,
		Entrypoint:  []string{"/bin/sh"},
	}
	mockExecCmd := &mockExecCmd{}
	mockExecCommandFactory := &mockExecCommandFactory{
		commandExecutor: mockExecCmd,
	}
	mockContainerStatusManager := &mockStatusHandler{
		getStatusFromIdStatus: status.RUNNING,
		getStatusFromIdErr:    nil,
		getPidFromIdPid:       11111,
	}
	containerExec := &ContainerExec{
		commandFactory:         mockExecCommandFactory,
		containerStatusManager: mockContainerStatusManager,
	}

	// == act ==
	err := containerExec.Exec(opt)

	// == assert ==
	// Command() parameter: nsenter -t <pid> --all <entrypoint>
	commandCall := mockExecCommandFactory.commandCalls[0]
	assert.Equal(t, "nsenter", commandCall.name)
	assert.Equal(t, []string{"-t", "11111", "--all", "/bin/sh"}, commandCall.args)

	// SetStdout/SetStderr/SetStdin is not called
	assert.False(t, mockExecCmd.setStdoutCallFlag)
	assert.False(t, mockExecCmd.setStderrCallFlag)
	assert.False(t, mockExecCmd.setStdinCallFlag)

	// Wait() is not called
	assert.False(t, mockExecCmd.waitCallFlag)

	// error is nil
	assert.Nil(t, err)
}

func TestContainerExec_Exec_GetStatusError(t *testing.T) {
	// == arrange ==
	opt := ExecOption{
		ContainerId: "12345",
		Tty:         false,
		Entrypoint:  []string{"/bin/sh"},
	}
	mockExecCmd := &mockExecCmd{}
	mockExecCommandFactory := &mockExecCommandFactory{
		commandExecutor: mockExecCmd,
	}
	mockContainerStatusManager := &mockStatusHandler{
		getStatusFromIdErr: errors.New("get status failed"),
	}
	containerExec := &ContainerExec{
		commandFactory:         mockExecCommandFactory,
		containerStatusManager: mockContainerStatusManager,
	}

	// == act ==
	err := containerExec.Exec(opt)

	// == assert ==
	// error is not nil
	assert.NotNil(t, err)
	assert.Equal(t, errors.New("get status failed"), err)
}

func TestContainerExec_Exec_StatusNotRunningError(t *testing.T) {
	// == arrange ==
	opt := ExecOption{
		ContainerId: "12345",
		Tty:         false,
		Entrypoint:  []string{"/bin/sh"},
	}
	mockExecCmd := &mockExecCmd{}
	mockExecCommandFactory := &mockExecCommandFactory{
		commandExecutor: mockExecCmd,
	}
	mockContainerStatusManager := &mockStatusHandler{
		getStatusFromIdStatus: status.STOPPED,
	}
	containerExec := &ContainerExec{
		commandFactory:         mockExecCommandFactory,
		containerStatusManager: mockContainerStatusManager,
	}

	// == act ==
	err := containerExec.Exec(opt)

	// == assert ==
	// error is not nil
	assert.NotNil(t, err)
	assert.Equal(t, errors.New("container: 12345 not running."), err)
}

func TestContainerExec_Exec_GetPidError(t *testing.T) {
	// == arrange ==
	opt := ExecOption{
		ContainerId: "12345",
		Tty:         false,
		Entrypoint:  []string{"/bin/sh"},
	}
	mockExecCmd := &mockExecCmd{}
	mockExecCommandFactory := &mockExecCommandFactory{
		commandExecutor: mockExecCmd,
	}
	mockContainerStatusManager := &mockStatusHandler{
		getStatusFromIdStatus: status.RUNNING,
		getStatusFromIdErr:    nil,
		getPidFromIdPid:       11111,

		getPidFromIdErr: errors.New("get pid failed"),
	}
	containerExec := &ContainerExec{
		commandFactory:         mockExecCommandFactory,
		containerStatusManager: mockContainerStatusManager,
	}

	// == act ==
	err := containerExec.Exec(opt)

	// == assert ==
	// error is not nil
	assert.NotNil(t, err)
	assert.Equal(t, errors.New("get pid failed"), err)
}

func TestContainerExec_Exec_StartError(t *testing.T) {
	// == arrange ==
	opt := ExecOption{
		ContainerId: "12345",
		Tty:         false,
		Entrypoint:  []string{"/bin/sh"},
	}
	mockExecCmd := &mockExecCmd{
		startErr: errors.New("Start() failed"),
	}
	mockExecCommandFactory := &mockExecCommandFactory{
		commandExecutor: mockExecCmd,
	}
	mockContainerStatusManager := &mockStatusHandler{
		getStatusFromIdStatus: status.RUNNING,
		getStatusFromIdErr:    nil,
		getPidFromIdPid:       11111,
	}
	containerExec := &ContainerExec{
		commandFactory:         mockExecCommandFactory,
		containerStatusManager: mockContainerStatusManager,
	}

	// == act ==
	err := containerExec.Exec(opt)

	// == assert ==
	// error is not nil
	assert.NotNil(t, err)
	assert.Equal(t, errors.New("Start() failed"), err)
}

func TestContainerExec_Exec_WaitError(t *testing.T) {
	// == arrange ==
	opt := ExecOption{
		ContainerId: "12345",
		Tty:         true,
		Entrypoint:  []string{"/bin/sh"},
	}
	mockExecCmd := &mockExecCmd{
		waitErr: errors.New("Wait() failed"),
	}
	mockExecCommandFactory := &mockExecCommandFactory{
		commandExecutor: mockExecCmd,
	}
	mockContainerStatusManager := &mockStatusHandler{
		getStatusFromIdStatus: status.RUNNING,
		getStatusFromIdErr:    nil,
		getPidFromIdPid:       11111,
	}
	containerExec := &ContainerExec{
		commandFactory:         mockExecCommandFactory,
		containerStatusManager: mockContainerStatusManager,
	}

	// == act ==
	err := containerExec.Exec(opt)

	// == assert ==
	// error is not nil
	assert.NotNil(t, err)
	assert.Equal(t, errors.New("Wait() failed"), err)
}
