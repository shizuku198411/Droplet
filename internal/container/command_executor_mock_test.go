package container

import (
	"io"
	"syscall"
)

type mockExecCommandFactory struct {
	// Command()
	commandCallFlag bool
	commandName     string
	commandArgs     []string
	commandExecutor commandExecutor
}

func (m *mockExecCommandFactory) Command(name string, args ...string) commandExecutor {
	m.commandCallFlag = true
	m.commandName = name
	m.commandArgs = args
	return m.commandExecutor
}

type mockExecCmd struct {
	// Start()
	startCallFlag bool
	startErr      error

	// Wait()
	waitCallFlag bool
	waitErr      error

	// Pid()
	pidCallFlag bool
	pidPid      int

	// SetStdout()
	setStdoutCallFlag bool
	setStdoutWritr    io.Writer

	// SetStderr()
	setStderrCallFlag bool
	setStderrWriter   io.Writer

	// SetStdin()
	setStdinCallFlag bool
	setStdinReader   io.Reader

	// SetSysProcAttr
	setSysProcAttrCallFlag bool
	setSysProcAttrAttr     *syscall.SysProcAttr
}

func (m *mockExecCmd) Start() error {
	m.startCallFlag = true
	return m.startErr
}

func (m *mockExecCmd) Wait() error {
	m.waitCallFlag = true
	return m.waitErr
}

func (m *mockExecCmd) Pid() int {
	m.pidCallFlag = true
	return m.pidPid
}

func (m *mockExecCmd) SetStdout(w io.Writer) {
	m.setStdoutCallFlag = true
	m.setStdoutWritr = w
}

func (m *mockExecCmd) SetStderr(w io.Writer) {
	m.setStderrCallFlag = true
	m.setStderrWriter = w
}

func (m *mockExecCmd) SetStdin(r io.Reader) {
	m.setStdinCallFlag = true
	m.setStdinReader = r
}

func (m *mockExecCmd) SetSysProcAttr(attr *syscall.SysProcAttr) {
	m.setSysProcAttrCallFlag = true
	m.setSysProcAttrAttr = attr
}

type mockKernelSyscall struct {
	// Exec()
	execCallFlag bool
	execArgv0    string
	execArgv     []string
	execEnvv     []string
	execErr      error

	// Setresgid()
	setresgidCallFlag bool
	setresgidRgid     int
	setresgidEgid     int
	setresgidSgid     int
	setresgidErr      error

	// Setresuid()
	setresuidCallFlag bool
	setresuidRuid     int
	setresuidEuid     int
	setresuidSuid     int
	setresuidErr      error

	// Sethostname()
	sethostnameCallFlag bool
	sethostnameP        []byte
	sethostnameErr      error
}

func (m *mockKernelSyscall) Exec(argv0 string, argv []string, envv []string) error {
	m.execCallFlag = true
	m.execArgv0 = argv0
	m.execArgv = argv
	m.execEnvv = envv
	return m.execErr
}

func (m *mockKernelSyscall) Setresgid(rgid int, egid int, sgid int) error {
	m.setresgidCallFlag = true
	m.setresgidRgid = rgid
	m.setresgidEgid = egid
	m.setresgidSgid = sgid
	return m.setresgidErr
}

func (m *mockKernelSyscall) Setresuid(ruid int, euid int, suid int) error {
	m.setresuidCallFlag = true
	m.setresuidRuid = ruid
	m.setresuidEuid = euid
	m.setresuidSuid = suid
	return m.setresuidErr
}

func (m *mockKernelSyscall) Sethostname(p []byte) error {
	m.sethostnameCallFlag = true
	m.sethostnameP = p
	return m.sethostnameErr
}
