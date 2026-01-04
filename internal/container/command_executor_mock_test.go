package container

import (
	"io"
	"os"
	"syscall"
	"time"
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

	// Mount()
	mountCallFlag bool
	mountSource   string
	mountTarget   string
	mountFstype   string
	mountFlags    uintptr
	mountData     string
	mountErr      error

	// Unmount()
	unmountCallFlag bool
	unmountTarget   string
	unmountFlags    int
	unmountErr      error

	// PivotRoot()
	pivotRootCallFlag bool
	pivotRootNewroot  string
	pivotRootPutold   string
	pivotRootErr      error

	// Chdir()
	chdirCallFlag bool
	chdirPath     string
	chdirErr      error

	// Mkdir()
	mkdirCallFlag bool
	mkdirPath     string
	mkdirMode     uint32
	mkdirErr      error

	// MkdirAll()
	mkdirAllCallFlag bool
	mkdirAllPath     string
	mkdirAllPerm     os.FileMode
	mkdirAllErr      error

	// Rmdir()
	rmdirCallFlag bool
	rmdirPath     string
	rmdirErr      error

	// Stat()
	statCallFlag bool
	statName     string
	statFileInfo os.FileInfo
	statErr      error

	// Create()
	createCallFlag bool
	createName     string
	createFile     *os.File
	createErr      error

	// IsNotExist()
	isNotExistCallFlag bool
	isNotExistErr      error
	isNotExistBool     bool

	// Symlink()
	symlinkCallFlag bool
	symlinkOldname  string
	symlinkNewname  string
	symlinkErr      error

	// Lstat()
	lstatCallFlag bool
	lstatName     string
	lstatFileInfo os.FileInfo
	lstatErr      error

	// OpenFile()
	openFileCallFlag bool
	openFileName     string
	openFileFlag     int
	openFilePerm     os.FileMode
	openFileFile     *os.File
	openFileErr      error
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

func (m *mockKernelSyscall) Mount(source string, target string, fstype string, flags uintptr, data string) error {
	m.mountCallFlag = true
	m.mountSource = source
	m.mountTarget = target
	m.mountFstype = fstype
	m.mountFlags = flags
	m.mountData = data
	return m.mountErr
}

func (m *mockKernelSyscall) Unmount(target string, flags int) error {
	m.unmountCallFlag = true
	m.unmountTarget = target
	m.unmountFlags = flags
	return m.unmountErr
}

func (m *mockKernelSyscall) PivotRoot(newroot string, putold string) error {
	m.pivotRootCallFlag = true
	m.pivotRootNewroot = newroot
	m.pivotRootPutold = putold
	return m.pivotRootErr
}

func (m *mockKernelSyscall) Chdir(path string) error {
	m.chdirCallFlag = true
	m.chdirPath = path
	return m.chdirErr
}

func (m *mockKernelSyscall) Mkdir(path string, mode uint32) error {
	m.mkdirCallFlag = true
	m.mkdirPath = path
	m.mkdirMode = mode
	return m.mkdirErr
}

func (m *mockKernelSyscall) MkdirAll(path string, perm os.FileMode) error {
	m.mkdirAllCallFlag = true
	m.mkdirAllPath = path
	m.mkdirAllPerm = perm
	return m.mkdirAllErr
}

func (m *mockKernelSyscall) Rmdir(path string) error {
	m.rmdirCallFlag = true
	m.rmdirPath = path
	return m.rmdirErr
}

func (m *mockKernelSyscall) Stat(name string) (os.FileInfo, error) {
	m.statCallFlag = true
	m.statName = name
	return m.statFileInfo, m.statErr
}

func (m *mockKernelSyscall) Create(name string) (*os.File, error) {
	m.createCallFlag = true
	m.createName = name
	return m.createFile, m.createErr
}

func (m *mockKernelSyscall) IsNotExist(err error) bool {
	m.isNotExistCallFlag = true
	m.isNotExistErr = err
	return m.isNotExistBool
}

func (m *mockKernelSyscall) Symlink(oldname string, newname string) error {
	m.symlinkCallFlag = true
	m.symlinkOldname = oldname
	m.symlinkNewname = newname
	return m.symlinkErr
}

func (m *mockKernelSyscall) Lstat(name string) (os.FileInfo, error) {
	m.lstatCallFlag = true
	m.lstatName = name
	return m.lstatFileInfo, m.lstatErr
}

func (m *mockKernelSyscall) OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	m.openFileCallFlag = true
	m.openFileName = name
	m.openFileFlag = flag
	m.openFilePerm = perm
	return m.openFileFile, m.openFileErr
}

type mockFileInfo struct {
	dir bool
}

func (m mockFileInfo) Name() string       { return "" }
func (m mockFileInfo) Size() int64        { return 0 }
func (m mockFileInfo) Mode() os.FileMode  { return 0 }
func (m mockFileInfo) ModTime() time.Time { return time.Time{} }
func (m mockFileInfo) IsDir() bool        { return m.dir }
func (m mockFileInfo) Sys() interface{}   { return nil }
