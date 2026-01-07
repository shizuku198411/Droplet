package container

import (
	"droplet/internal/utils"
	"io"
	"os"
	"syscall"
	"time"
)

type mockExecCommandFactory struct {
	// Command()
	commandCalls    []commandParameter
	commandExecutor utils.CommandExecutor
}

type commandParameter struct {
	name string
	args []string
}

func (m *mockExecCommandFactory) Command(name string, args ...string) utils.CommandExecutor {
	m.commandCalls = append(m.commandCalls, commandParameter{
		name: name,
		args: args,
	})
	return m.commandExecutor
}

type mockExecCmd struct {
	// Start()
	startCallFlag bool
	startErr      error

	// Wait()
	waitCallFlag bool
	waitErr      error

	// Run()
	runCallFlag bool
	runErr      error

	// Pid()
	pidCallFlag bool
	pidPid      int

	// SetEnv()
	setEnvCallFlag bool
	setEnvValue    []string

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

func (m *mockExecCmd) Run() error {
	m.runCallFlag = true
	return m.runErr
}

func (m *mockExecCmd) Pid() int {
	m.pidCallFlag = true
	return m.pidPid
}

func (m *mockExecCmd) SetEnv(envv []string) {
	m.setEnvCallFlag = true
	m.setEnvValue = envv
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
	mountCalls []mountCallParameter
	mountErr   error

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

	// Remove()
	removeCallFlag bool
	removeName     string
	removeErr      error

	// IsNotExist()
	isNotExistCallFlag bool
	isNotExistErr      error
	isNotExistBool     bool

	// Symlink()
	symlinkCalls []symlinkParameter
	symlinkErr   error

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

	// WriteFile()
	writeFileCallFlag bool
	writeFileName     string
	writeFileData     []byte
	writeFilePerm     os.FileMode
	writeFileErr      error

	// Kill()
	killCallFlag bool
	killPid      int
	killSig      syscall.Signal
	killErr      error
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

type mountCallParameter struct {
	source string
	target string
	fstype string
	flags  uintptr
	data   string
}

func (m *mockKernelSyscall) Mount(source string, target string, fstype string, flags uintptr, data string) error {
	m.mountCalls = append(m.mountCalls, mountCallParameter{
		source: source,
		target: target,
		fstype: fstype,
		flags:  flags,
		data:   data,
	})
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

func (m *mockKernelSyscall) Remove(name string) error {
	m.removeCallFlag = true
	m.removeName = name
	return m.removeErr
}

func (m *mockKernelSyscall) IsNotExist(err error) bool {
	m.isNotExistCallFlag = true
	m.isNotExistErr = err
	return m.isNotExistBool
}

type symlinkParameter struct {
	oldname string
	newname string
}

func (m *mockKernelSyscall) Symlink(oldname string, newname string) error {
	m.symlinkCalls = append(m.symlinkCalls, symlinkParameter{
		oldname: oldname,
		newname: newname,
	})
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

func (m *mockKernelSyscall) WriteFile(name string, data []byte, perm os.FileMode) error {
	m.writeFileCallFlag = true
	m.writeFileName = name
	m.writeFileData = data
	m.writeFilePerm = perm
	return m.writeFileErr
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

func (m *mockKernelSyscall) Kill(pid int, sig syscall.Signal) error {
	m.killCallFlag = true
	m.killPid = pid
	m.killSig = sig
	return m.killErr
}
