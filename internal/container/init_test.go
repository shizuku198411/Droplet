package container

import (
	"errors"
	"syscall"
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

func TestPrepare_Success(t *testing.T) {
	// == arrange ==
	spec := buildMockSpec(t)
	mockKernelSyscall := &mockKernelSyscall{
		statFileInfo: mockFileInfo{dir: true},
	}
	rootContainerEnvPreparer := &rootContainerEnvPreparer{
		syscallHandler: mockKernelSyscall,
	}

	// == act ==
	err := rootContainerEnvPreparer.prepare(spec)

	// == assert ==
	// error is nil
	assert.Nil(t, err)
}

func TestPrepare_SetresgidError(t *testing.T) {
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

func TestPrepare_SetresuidError(t *testing.T) {
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

func TestPrepare_SethostnameError(t *testing.T) {
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

func TestSwitchToUserNamespaceRoot_Success(t *testing.T) {
	// == arrange ==
	mockKernelSyscall := &mockKernelSyscall{
		statFileInfo: mockFileInfo{dir: true},
	}
	rootContainerEnvPreparer := &rootContainerEnvPreparer{
		syscallHandler: mockKernelSyscall,
	}

	// == act ==
	err := rootContainerEnvPreparer.switchToUserNamespaceRoot()

	// == assert ==
	// Setresgid() is called
	assert.True(t, mockKernelSyscall.setresgidCallFlag)
	// Setresgid() parameter is 0, 0, 0
	assert.Equal(t, 0, mockKernelSyscall.setresgidRgid)
	assert.Equal(t, 0, mockKernelSyscall.setresgidEgid)
	assert.Equal(t, 0, mockKernelSyscall.setresgidSgid)
	// Setresuid() is called
	assert.True(t, mockKernelSyscall.setresuidCallFlag)
	// Setresuid() parameter is 0, 0, 0
	assert.Equal(t, 0, mockKernelSyscall.setresuidRuid)
	assert.Equal(t, 0, mockKernelSyscall.setresuidEuid)
	assert.Equal(t, 0, mockKernelSyscall.setresuidSuid)
	// error is nil
	assert.Nil(t, err)
}

func TestSwitcToUserNamespaceRoot_SetresgidError(t *testing.T) {
	// == arrange ==
	mockKernelSyscall := &mockKernelSyscall{
		statFileInfo: mockFileInfo{dir: true},
		setresgidErr: errors.New("Setresgid() failed"),
	}
	rootContainerEnvPreparer := &rootContainerEnvPreparer{
		syscallHandler: mockKernelSyscall,
	}

	// == act ==
	err := rootContainerEnvPreparer.switchToUserNamespaceRoot()

	// == assert ==
	// error is not nil
	assert.NotNil(t, err)
	assert.Equal(t, errors.New("Setresgid() failed"), err)
}

func TestSwitcToUserNamespaceRoot_SetresuidError(t *testing.T) {
	// == arrange ==
	mockKernelSyscall := &mockKernelSyscall{
		statFileInfo: mockFileInfo{dir: true},
		setresuidErr: errors.New("Setresuid() failed"),
	}
	rootContainerEnvPreparer := &rootContainerEnvPreparer{
		syscallHandler: mockKernelSyscall,
	}

	// == act ==
	err := rootContainerEnvPreparer.switchToUserNamespaceRoot()

	// == assert ==
	// error is not nil
	assert.NotNil(t, err)
	assert.Equal(t, errors.New("Setresuid() failed"), err)
}

func TestSetHostnameToContainerId_Success(t *testing.T) {
	// == arrange ==
	spec := buildMockSpec(t)
	mockKernelSyscall := &mockKernelSyscall{
		statFileInfo: mockFileInfo{dir: true},
	}
	rootContainerEnvPreparer := &rootContainerEnvPreparer{
		syscallHandler: mockKernelSyscall,
	}

	// == act ==
	err := rootContainerEnvPreparer.setHostnameToContainerId(spec.Hostname)

	// == assert ==
	// Sethostname() is called
	assert.True(t, mockKernelSyscall.sethostnameCallFlag)
	// Sethostname() parameter is spec.Hostanme
	assert.Equal(t, []byte(spec.Hostname), mockKernelSyscall.sethostnameP)
	// error is nil
	assert.Nil(t, err)
}

func TestSetHostnameToContainerId_SethostnameError(t *testing.T) {
	// == arrange ==
	spec := buildMockSpec(t)
	mockKernelSyscall := &mockKernelSyscall{
		statFileInfo:   mockFileInfo{dir: true},
		sethostnameErr: errors.New("Sethostname() failed"),
	}
	rootContainerEnvPreparer := &rootContainerEnvPreparer{
		syscallHandler: mockKernelSyscall,
	}

	// == act ==
	err := rootContainerEnvPreparer.setHostnameToContainerId(spec.Hostname)

	// == assert ==
	// error is not nil
	assert.NotNil(t, err)
	assert.Equal(t, errors.New("Sethostname() failed"), err)
}

func TestSetupOverlay_Success(t *testing.T) {
	// == arrange ==
	spec := buildMockSpec(t)
	mockKernelSyscall := &mockKernelSyscall{
		statFileInfo: mockFileInfo{dir: true},
	}
	rootContainerEnvPreparer := &rootContainerEnvPreparer{
		syscallHandler: mockKernelSyscall,
	}

	// == act ==
	err := rootContainerEnvPreparer.setupOverlay(spec.Root.Path, spec.Annotations.Image)

	// == assert ==
	// Mount() call time: 1
	assert.Equal(t, 1, len(mockKernelSyscall.mountCalls))
	// Mount() parameter is:
	//   source: "overlay"
	//   destination: spec.root.path
	//   type: image.rootfstype
	//   flags: 0
	//   data: lowerdir=<image.imagelayer>,upperdir=<image.upperdir>,workdir=<image.workdir>
	mountCall1 := mockKernelSyscall.mountCalls[0]
	assert.Equal(t, "overlay", mountCall1.source)
	assert.Equal(t, spec.Root.Path, mountCall1.target)
	assert.Equal(t, "overlay", mountCall1.fstype)
	assert.Equal(t, uintptr(0), mountCall1.flags)
	assert.Equal(t, "lowerdir=/image/path,upperdir=/upper/path,workdir=/work/path", mountCall1.data)
	// error is nil
	assert.Nil(t, err)
}

func TestSetupOverlay_MountError(t *testing.T) {
	// == arrange ==
	spec := buildMockSpec(t)
	mockKernelSyscall := &mockKernelSyscall{
		statFileInfo: mockFileInfo{dir: true},
		mountErr:     errors.New("Mount() failed"),
	}
	rootContainerEnvPreparer := &rootContainerEnvPreparer{
		syscallHandler: mockKernelSyscall,
	}

	// == act ==
	err := rootContainerEnvPreparer.setupOverlay(spec.Root.Path, spec.Annotations.Image)

	// == assert ==
	// error is not nil
	assert.NotNil(t, err)
	assert.Equal(t, errors.New("Mount() failed"), err)
}

func TestMountFilesystem_Success(t *testing.T) {
	// == arrange ==
	spec := buildMockSpec(t)
	mockKernelSyscall := &mockKernelSyscall{
		statFileInfo: mockFileInfo{dir: true},
	}
	rootContainerEnvPreparer := &rootContainerEnvPreparer{
		syscallHandler: mockKernelSyscall,
	}

	// == act ==
	err := rootContainerEnvPreparer.mountFilesystem(spec.Root.Path, spec.Mounts)

	// == assert ==
	// Mount() call time: 14
	assert.Equal(t, 14, len(mockKernelSyscall.mountCalls))

	// Mount() call 1: /proc
	mountCall1 := mockKernelSyscall.mountCalls[0]
	assert.Equal(t, "proc", mountCall1.source)
	assert.Equal(t, spec.Root.Path+"/proc", mountCall1.target)
	assert.Equal(t, "proc", mountCall1.fstype)
	assert.Equal(t, uintptr(syscall.MS_NOSUID|syscall.MS_NOEXEC|syscall.MS_NODEV), mountCall1.flags)
	assert.Equal(t, "", mountCall1.data)

	// Mount() call 2: /dev
	mountCall2 := mockKernelSyscall.mountCalls[1]
	assert.Equal(t, "tmpfs", mountCall2.source)
	assert.Equal(t, spec.Root.Path+"/dev", mountCall2.target)
	assert.Equal(t, "tmpfs", mountCall2.fstype)
	assert.Equal(t, uintptr(syscall.MS_NOSUID|syscall.MS_STRICTATIME), mountCall2.flags)
	assert.Equal(t, "mode=755,size=65536k", mountCall2.data)

	// Mount() call 3: /dev/pts
	mountCall3 := mockKernelSyscall.mountCalls[2]
	assert.Equal(t, "devpts", mountCall3.source)
	assert.Equal(t, spec.Root.Path+"/dev/pts", mountCall3.target)
	assert.Equal(t, "devpts", mountCall3.fstype)
	assert.Equal(t, uintptr(syscall.MS_NOSUID|syscall.MS_NOEXEC), mountCall3.flags)
	assert.Equal(t, "newinstance,ptmxmode=0666,mode=0620,gid=5", mountCall3.data)

	// Mount() call 4: /sys
	mountCall4 := mockKernelSyscall.mountCalls[3]
	assert.Equal(t, "sysfs", mountCall4.source)
	assert.Equal(t, spec.Root.Path+"/sys", mountCall4.target)
	assert.Equal(t, "sysfs", mountCall4.fstype)
	assert.Equal(t, uintptr(syscall.MS_NOSUID|syscall.MS_NOEXEC|syscall.MS_NODEV|syscall.MS_RDONLY), mountCall4.flags)
	assert.Equal(t, "", mountCall4.data)

	// Mount() call 5: /sys/fs/cgroup
	mountCall5 := mockKernelSyscall.mountCalls[4]
	assert.Equal(t, "cgroup", mountCall5.source)
	assert.Equal(t, spec.Root.Path+"/sys/fs/cgroup", mountCall5.target)
	assert.Equal(t, "cgroup", mountCall5.fstype)
	assert.Equal(t, uintptr(syscall.MS_NOSUID|syscall.MS_NOEXEC|syscall.MS_NODEV|syscall.MS_RDONLY), mountCall5.flags)
	assert.Equal(t, "", mountCall5.data)

	// Mount() call 6: /dev/mqueue
	mountCall6 := mockKernelSyscall.mountCalls[5]
	assert.Equal(t, "mqueue", mountCall6.source)
	assert.Equal(t, spec.Root.Path+"/dev/mqueue", mountCall6.target)
	assert.Equal(t, "mqueue", mountCall6.fstype)
	assert.Equal(t, uintptr(syscall.MS_NOSUID|syscall.MS_NOEXEC|syscall.MS_NODEV), mountCall6.flags)
	assert.Equal(t, "", mountCall6.data)

	// Mount() call 7: /dev/shm
	mountCall7 := mockKernelSyscall.mountCalls[6]
	assert.Equal(t, "shm", mountCall7.source)
	assert.Equal(t, spec.Root.Path+"/dev/shm", mountCall7.target)
	assert.Equal(t, "tmpfs", mountCall7.fstype)
	assert.Equal(t, uintptr(syscall.MS_NOSUID|syscall.MS_NOEXEC|syscall.MS_NODEV), mountCall7.flags)
	assert.Equal(t, "mode=1777,size=67108864", mountCall7.data)

	// Mount() call 8: /etc/resolv.conf
	mountCall8 := mockKernelSyscall.mountCalls[7]
	assert.Equal(t, "/etc/raind/container/mycontainer/etc/resolv.conf", mountCall8.source)
	assert.Equal(t, spec.Root.Path+"/etc/resolv.conf", mountCall8.target)
	assert.Equal(t, "bind", mountCall8.fstype)
	assert.Equal(t, uintptr(syscall.MS_BIND|syscall.MS_REC), mountCall8.flags)
	assert.Equal(t, "", mountCall8.data)

	// Mount() call 9: /etc/resolv.conf remount
	mountCall9 := mockKernelSyscall.mountCalls[8]
	assert.Equal(t, "", mountCall9.source)
	assert.Equal(t, spec.Root.Path+"/etc/resolv.conf", mountCall9.target)
	assert.Equal(t, "", mountCall9.fstype)
	assert.Equal(t, uintptr(syscall.MS_PRIVATE|syscall.MS_REC), mountCall9.flags)
	assert.Equal(t, "", mountCall9.data)

	// Mount() call 10: /etc/hostname
	mountCall10 := mockKernelSyscall.mountCalls[9]
	assert.Equal(t, "/etc/raind/container/mycontainer/etc/hostname", mountCall10.source)
	assert.Equal(t, spec.Root.Path+"/etc/hostname", mountCall10.target)
	assert.Equal(t, "bind", mountCall10.fstype)
	assert.Equal(t, uintptr(syscall.MS_BIND|syscall.MS_REC), mountCall10.flags)
	assert.Equal(t, "", mountCall10.data)

	// Mount() call 11: /etc/hostname remount
	mountCall11 := mockKernelSyscall.mountCalls[10]
	assert.Equal(t, "", mountCall11.source)
	assert.Equal(t, spec.Root.Path+"/etc/hostname", mountCall11.target)
	assert.Equal(t, "", mountCall11.fstype)
	assert.Equal(t, uintptr(syscall.MS_PRIVATE|syscall.MS_REC), mountCall11.flags)
	assert.Equal(t, "", mountCall11.data)

	// Mount() call 12: /etc/hosts
	mountCall12 := mockKernelSyscall.mountCalls[11]
	assert.Equal(t, "/etc/raind/container/mycontainer/etc/hosts", mountCall12.source)
	assert.Equal(t, spec.Root.Path+"/etc/hosts", mountCall12.target)
	assert.Equal(t, "bind", mountCall12.fstype)
	assert.Equal(t, uintptr(syscall.MS_BIND|syscall.MS_REC), mountCall12.flags)
	assert.Equal(t, "", mountCall12.data)

	// Mount() call 13: /etc/hosts remount
	mountCall13 := mockKernelSyscall.mountCalls[12]
	assert.Equal(t, "", mountCall13.source)
	assert.Equal(t, spec.Root.Path+"/etc/hosts", mountCall13.target)
	assert.Equal(t, "", mountCall13.fstype)
	assert.Equal(t, uintptr(syscall.MS_PRIVATE|syscall.MS_REC), mountCall13.flags)
	assert.Equal(t, "", mountCall13.data)

	// Mount() call 14: user-specified mount
	mountCall14 := mockKernelSyscall.mountCalls[13]
	assert.Equal(t, "/src", mountCall14.source)
	assert.Equal(t, spec.Root.Path+"/dst", mountCall14.target)
	assert.Equal(t, "", mountCall14.fstype)
	assert.Equal(t, uintptr(syscall.MS_BIND), mountCall14.flags)
	assert.Equal(t, "", mountCall14.data)

	// error is nil
	assert.Nil(t, err)
}

func TestMountFilesystem_MountError(t *testing.T) {
	// == arrange ==
	spec := buildMockSpec(t)
	mockKernelSyscall := &mockKernelSyscall{
		statFileInfo: mockFileInfo{dir: true},
		mountErr:     errors.New("Mount() failed"),
	}
	rootContainerEnvPreparer := &rootContainerEnvPreparer{
		syscallHandler: mockKernelSyscall,
	}

	// == act ==
	err := rootContainerEnvPreparer.mountFilesystem(spec.Root.Path, spec.Mounts)

	// == assert ==
	// error is not nil
	assert.NotNil(t, err)
	assert.Equal(t, errors.New("Mount() failed"), err)
}

func TestMountStdDevice_Success(t *testing.T) {
	// == arrange ==
	spec := buildMockSpec(t)
	mockKernelSyscall := &mockKernelSyscall{}
	rootContainerEnvPreparer := &rootContainerEnvPreparer{
		syscallHandler: mockKernelSyscall,
	}

	// == act ==
	err := rootContainerEnvPreparer.mountStdDevice(spec.Root.Path)

	// == assert ==
	// Mount() call time: 12
	assert.Equal(t, 12, len(mockKernelSyscall.mountCalls))

	// Mount() call 1: /dev/random
	mountCall1 := mockKernelSyscall.mountCalls[0]
	assert.Equal(t, "/dev/random", mountCall1.source)
	assert.Equal(t, spec.Root.Path+"/dev/random", mountCall1.target)
	assert.Equal(t, "", mountCall1.fstype)
	assert.Equal(t, uintptr(syscall.MS_BIND), mountCall1.flags)
	assert.Equal(t, "", mountCall1.data)

	// Mount() call 2: /dev/random remount
	mountCall2 := mockKernelSyscall.mountCalls[1]
	assert.Equal(t, "", mountCall2.source)
	assert.Equal(t, spec.Root.Path+"/dev/random", mountCall2.target)
	assert.Equal(t, "", mountCall2.fstype)
	assert.Equal(t, uintptr(syscall.MS_BIND|syscall.MS_REMOUNT|syscall.MS_RDONLY|syscall.MS_NOEXEC|syscall.MS_NOSUID), mountCall2.flags)
	assert.Equal(t, "", mountCall2.data)

	// Mount() call 3: /dev/urandom
	mountCall3 := mockKernelSyscall.mountCalls[2]
	assert.Equal(t, "/dev/urandom", mountCall3.source)
	assert.Equal(t, spec.Root.Path+"/dev/urandom", mountCall3.target)
	assert.Equal(t, "", mountCall3.fstype)
	assert.Equal(t, uintptr(syscall.MS_BIND), mountCall3.flags)
	assert.Equal(t, "", mountCall3.data)

	// Mount() call 4: /dev/urandom remount
	mountCall4 := mockKernelSyscall.mountCalls[3]
	assert.Equal(t, "", mountCall4.source)
	assert.Equal(t, spec.Root.Path+"/dev/urandom", mountCall4.target)
	assert.Equal(t, "", mountCall4.fstype)
	assert.Equal(t, uintptr(syscall.MS_BIND|syscall.MS_REMOUNT|syscall.MS_RDONLY|syscall.MS_NOEXEC|syscall.MS_NOSUID), mountCall4.flags)
	assert.Equal(t, "", mountCall4.data)

	// Mount() call 5: /dev/null
	mountCall5 := mockKernelSyscall.mountCalls[4]
	assert.Equal(t, "/dev/null", mountCall5.source)
	assert.Equal(t, spec.Root.Path+"/dev/null", mountCall5.target)
	assert.Equal(t, "", mountCall5.fstype)
	assert.Equal(t, uintptr(syscall.MS_BIND), mountCall5.flags)
	assert.Equal(t, "", mountCall5.data)

	// Mount() call 6: /dev/null remount
	mountCall6 := mockKernelSyscall.mountCalls[5]
	assert.Equal(t, "", mountCall6.source)
	assert.Equal(t, spec.Root.Path+"/dev/null", mountCall6.target)
	assert.Equal(t, "", mountCall6.fstype)
	assert.Equal(t, uintptr(syscall.MS_BIND|syscall.MS_REMOUNT|syscall.MS_RDONLY|syscall.MS_NOEXEC|syscall.MS_NOSUID), mountCall6.flags)
	assert.Equal(t, "", mountCall6.data)

	// Mount() call 7: /dev/zero
	mountCall7 := mockKernelSyscall.mountCalls[6]
	assert.Equal(t, "/dev/zero", mountCall7.source)
	assert.Equal(t, spec.Root.Path+"/dev/zero", mountCall7.target)
	assert.Equal(t, "", mountCall7.fstype)
	assert.Equal(t, uintptr(syscall.MS_BIND), mountCall7.flags)
	assert.Equal(t, "", mountCall7.data)

	// Mount() call 8: /dev/zero remount
	mountCall8 := mockKernelSyscall.mountCalls[7]
	assert.Equal(t, "", mountCall8.source)
	assert.Equal(t, spec.Root.Path+"/dev/zero", mountCall8.target)
	assert.Equal(t, "", mountCall8.fstype)
	assert.Equal(t, uintptr(syscall.MS_BIND|syscall.MS_REMOUNT|syscall.MS_RDONLY|syscall.MS_NOEXEC|syscall.MS_NOSUID), mountCall8.flags)
	assert.Equal(t, "", mountCall8.data)

	// Mount() call 9: /dev/full
	mountCall9 := mockKernelSyscall.mountCalls[8]
	assert.Equal(t, "/dev/full", mountCall9.source)
	assert.Equal(t, spec.Root.Path+"/dev/full", mountCall9.target)
	assert.Equal(t, "", mountCall9.fstype)
	assert.Equal(t, uintptr(syscall.MS_BIND), mountCall9.flags)
	assert.Equal(t, "", mountCall9.data)

	// Mount() call 10: /dev/full remount
	mountCall10 := mockKernelSyscall.mountCalls[9]
	assert.Equal(t, "", mountCall10.source)
	assert.Equal(t, spec.Root.Path+"/dev/full", mountCall10.target)
	assert.Equal(t, "", mountCall10.fstype)
	assert.Equal(t, uintptr(syscall.MS_BIND|syscall.MS_REMOUNT|syscall.MS_RDONLY|syscall.MS_NOEXEC|syscall.MS_NOSUID), mountCall10.flags)
	assert.Equal(t, "", mountCall10.data)

	// Mount() call 11: /dev/tty
	mountCall11 := mockKernelSyscall.mountCalls[10]
	assert.Equal(t, "/dev/tty", mountCall11.source)
	assert.Equal(t, spec.Root.Path+"/dev/tty", mountCall11.target)
	assert.Equal(t, "", mountCall11.fstype)
	assert.Equal(t, uintptr(syscall.MS_BIND), mountCall11.flags)
	assert.Equal(t, "", mountCall11.data)

	// Mount() call 12: /dev/tty remount
	mountCall12 := mockKernelSyscall.mountCalls[11]
	assert.Equal(t, "", mountCall12.source)
	assert.Equal(t, spec.Root.Path+"/dev/tty", mountCall12.target)
	assert.Equal(t, "", mountCall12.fstype)
	assert.Equal(t, uintptr(syscall.MS_BIND|syscall.MS_REMOUNT|syscall.MS_RDONLY|syscall.MS_NOEXEC|syscall.MS_NOSUID), mountCall12.flags)
	assert.Equal(t, "", mountCall12.data)

	// error is nil
	assert.Nil(t, err)
}

func TestMountStdDevice_MountError(t *testing.T) {
	// == arrange ==
	spec := buildMockSpec(t)
	mockKernelSyscall := &mockKernelSyscall{
		mountErr: errors.New("Mount() failed"),
	}
	rootContainerEnvPreparer := &rootContainerEnvPreparer{
		syscallHandler: mockKernelSyscall,
	}

	// == act ==
	err := rootContainerEnvPreparer.mountStdDevice(spec.Root.Path)

	// == assert ==
	// error is not nil
	assert.NotNil(t, err)
	assert.Equal(t, errors.New("Mount() failed"), err)
}

func TestCreateSymbolicLink_Success(t *testing.T) {
	// == arrange ==
	spec := buildMockSpec(t)
	mockKernelSyscall := &mockKernelSyscall{
		lstatErr: errors.New("file is not symbolic link"),
	}
	rootContainerEnvPreparer := &rootContainerEnvPreparer{
		syscallHandler: mockKernelSyscall,
	}

	// == act ==
	err := rootContainerEnvPreparer.createSymbolicLink(spec.Root.Path)

	// == assert ==
	// Symlink() call time: 5
	assert.Equal(t, 5, len(mockKernelSyscall.symlinkCalls))

	// Symlink() call 1: /dev/fd -> /proc/self/fd
	symlinkCall1 := mockKernelSyscall.symlinkCalls[0]
	assert.Equal(t, spec.Root.Path+"/dev/fd", symlinkCall1.newname)
	assert.Equal(t, "/proc/self/fd", symlinkCall1.oldname)

	// Symlink() call 2: /dev/stdin -> /proc/self/fd/0
	symlinkCall2 := mockKernelSyscall.symlinkCalls[1]
	assert.Equal(t, spec.Root.Path+"/dev/stdin", symlinkCall2.newname)
	assert.Equal(t, "/proc/self/fd/0", symlinkCall2.oldname)

	// Symlink() call 3: /dev/stdout -> /proc/self/fd/1
	symlinkCall3 := mockKernelSyscall.symlinkCalls[2]
	assert.Equal(t, spec.Root.Path+"/dev/stdout", symlinkCall3.newname)
	assert.Equal(t, "/proc/self/fd/1", symlinkCall3.oldname)

	// Symlink() call 4: /dev/stderr -> /proc/self/fd/2
	symlinkCall4 := mockKernelSyscall.symlinkCalls[3]
	assert.Equal(t, spec.Root.Path+"/dev/stderr", symlinkCall4.newname)
	assert.Equal(t, "/proc/self/fd/2", symlinkCall4.oldname)

	// Symlink() call 5: /dev/ptmx -> /dev/pts/ptmx
	symlinkCall5 := mockKernelSyscall.symlinkCalls[4]
	assert.Equal(t, spec.Root.Path+"/dev/ptmx", symlinkCall5.newname)
	assert.Equal(t, "/dev/pts/ptmx", symlinkCall5.oldname)

	// error is nil
	assert.Nil(t, err)
}

func TestCreateSymbolicLink_LstatIsNil_Success(t *testing.T) {
	// == arrange ==
	spec := buildMockSpec(t)
	mockKernelSyscall := &mockKernelSyscall{
		lstatErr: nil,
	}
	rootContainerEnvPreparer := &rootContainerEnvPreparer{
		syscallHandler: mockKernelSyscall,
	}

	// == act ==
	err := rootContainerEnvPreparer.createSymbolicLink(spec.Root.Path)

	// == assert ==
	// Symlink() call time: 0
	assert.Equal(t, 0, len(mockKernelSyscall.symlinkCalls))
	// error is nil
	assert.Nil(t, err)
}

func TestCreateSymbolicLink_SymlinkError(t *testing.T) {
	// == arrange ==
	spec := buildMockSpec(t)
	mockKernelSyscall := &mockKernelSyscall{
		lstatErr:   errors.New("file is not symbolic link"),
		symlinkErr: errors.New("Symlink() failed"),
	}
	rootContainerEnvPreparer := &rootContainerEnvPreparer{
		syscallHandler: mockKernelSyscall,
	}

	// == act ==
	err := rootContainerEnvPreparer.createSymbolicLink(spec.Root.Path)

	// == assert ==
	// error is not nil
	assert.NotNil(t, err)
	assert.Equal(t, errors.New("Symlink() failed"), err)
}

func TestPivotRoot_Success(t *testing.T) {
	// == arrange ==
	spec := buildMockSpec(t)
	mockKernelSyscall := &mockKernelSyscall{}
	rootContainerEnvPreparer := &rootContainerEnvPreparer{
		syscallHandler: mockKernelSyscall,
	}

	// == act ==
	err := rootContainerEnvPreparer.pivotRoot(spec.Root.Path)

	// == assert ==
	// Mkdir() is called
	assert.True(t, mockKernelSyscall.mkdirCallFlag)
	// Mkdir() parameter: path=/etc/raind/container/<container-id>/merged/put_old
	assert.Equal(t, spec.Root.Path+"/put_old", mockKernelSyscall.mkdirPath)

	// PivotRoot() is called
	assert.True(t, mockKernelSyscall.pivotRootCallFlag)
	// PivotRoot() parameter:
	//   newroot=/etr/raind/container/<container-id>/merged
	//   putold=/etr/raind/container/<container-id>/merged/put_old
	assert.Equal(t, spec.Root.Path, mockKernelSyscall.pivotRootNewroot)
	assert.Equal(t, spec.Root.Path+"/put_old", mockKernelSyscall.pivotRootPutold)

	// Chdir() is called
	assert.True(t, mockKernelSyscall.chdirCallFlag)
	// Chdir() parameter: path=/
	assert.Equal(t, "/", mockKernelSyscall.chdirPath)

	// Unmount() is called
	assert.True(t, mockKernelSyscall.unmountCallFlag)
	// Unmount() parameter:
	//   target=/put_old
	//   flag=MNT_DETACH
	assert.Equal(t, "/put_old", mockKernelSyscall.unmountTarget)
	assert.Equal(t, syscall.MNT_DETACH, mockKernelSyscall.unmountFlags)

	// Rmdir() is called
	assert.True(t, mockKernelSyscall.rmdirCallFlag)
	// Rmdir() parameter: path=/put_old
	assert.Equal(t, "/put_old", mockKernelSyscall.rmdirPath)

	// error is nil
	assert.Nil(t, err)
}

func TestPivotRoot_MkdirError(t *testing.T) {
	// == arrange ==
	spec := buildMockSpec(t)
	mockKernelSyscall := &mockKernelSyscall{
		mkdirErr: errors.New("Mkdir() failed"),
	}
	rootContainerEnvPreparer := &rootContainerEnvPreparer{
		syscallHandler: mockKernelSyscall,
	}

	// == act ==
	err := rootContainerEnvPreparer.pivotRoot(spec.Root.Path)

	// == assert ==
	// error is not nil
	assert.NotNil(t, err)
	assert.Equal(t, errors.New("Mkdir() failed"), err)
}

func TestPivotRoot_PivotRootError(t *testing.T) {
	// == arrange ==
	spec := buildMockSpec(t)
	mockKernelSyscall := &mockKernelSyscall{
		pivotRootErr: errors.New("PivotRoot() failed"),
	}
	rootContainerEnvPreparer := &rootContainerEnvPreparer{
		syscallHandler: mockKernelSyscall,
	}

	// == act ==
	err := rootContainerEnvPreparer.pivotRoot(spec.Root.Path)

	// == assert ==
	// error is not nil
	assert.NotNil(t, err)
	assert.Equal(t, errors.New("PivotRoot() failed"), err)
}

func TestPivotRoot_ChdirError(t *testing.T) {
	// == arrange ==
	spec := buildMockSpec(t)
	mockKernelSyscall := &mockKernelSyscall{
		chdirErr: errors.New("Chdir() failed"),
	}
	rootContainerEnvPreparer := &rootContainerEnvPreparer{
		syscallHandler: mockKernelSyscall,
	}

	// == act ==
	err := rootContainerEnvPreparer.pivotRoot(spec.Root.Path)

	// == assert ==
	// error is not nil
	assert.NotNil(t, err)
	assert.Equal(t, errors.New("Chdir() failed"), err)
}

func TestPivotRoot_UnmountError(t *testing.T) {
	// == arrange ==
	spec := buildMockSpec(t)
	mockKernelSyscall := &mockKernelSyscall{
		unmountErr: errors.New("Unmount() failed"),
	}
	rootContainerEnvPreparer := &rootContainerEnvPreparer{
		syscallHandler: mockKernelSyscall,
	}

	// == act ==
	err := rootContainerEnvPreparer.pivotRoot(spec.Root.Path)

	// == assert ==
	// error is not nil
	assert.NotNil(t, err)
	assert.Equal(t, errors.New("Unmount() failed"), err)
}

func TestPivotRoot_RmdirError(t *testing.T) {
	// == arrange ==
	spec := buildMockSpec(t)
	mockKernelSyscall := &mockKernelSyscall{
		rmdirErr: errors.New("Rmdir() failed"),
	}
	rootContainerEnvPreparer := &rootContainerEnvPreparer{
		syscallHandler: mockKernelSyscall,
	}

	// == act ==
	err := rootContainerEnvPreparer.pivotRoot(spec.Root.Path)

	// == assert ==
	// error is not nil
	assert.NotNil(t, err)
	assert.Equal(t, errors.New("Rmdir() failed"), err)
}
