package container

import (
	"errors"
	"fmt"
	"io/fs"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewContainerCgroupController_Success(t *testing.T) {
	// == arrange ==

	// == act ==
	got := newContainerCgroupController()

	// == assert ==
	assert.NotNil(t, got)
	assert.NotNil(t, got.syscallHandler)
}

func TestContainerCgroupController_Prepare_Success(t *testing.T) {
	// == arrange ==
	spec := buildMockSpec(t)
	containerId := "12345"
	initPid := 11111
	mockKernelSyscall := &mockKernelSyscall{}
	containerCgroupPreparer := &containerCgroupController{
		syscallHandler: mockKernelSyscall,
	}

	// == act ==
	err := containerCgroupPreparer.prepare(containerId, spec, initPid)

	// == assert ==
	// error is nil
	assert.Nil(t, err)
}

func TestContainerCgroupController_SetMemoryLimit_Success(t *testing.T) {
	// == arrange ==
	spec := buildMockSpec(t)
	containerId := "12345"
	mockKernelSyscall := &mockKernelSyscall{}
	containerCgroupPreparer := &containerCgroupController{
		syscallHandler: mockKernelSyscall,
	}

	// == act ==
	err := containerCgroupPreparer.setMemoryLimit(containerId, spec.LinuxSpec.Resources.Memory)

	// == assert ==
	// WriteFile() is called
	assert.True(t, mockKernelSyscall.writeFileCallFlag)
	// WriteFile() parameter:
	//   path=/sys/fs/cgroup/raind/<container-id>/memory.max
	//   data=536870912
	//   mode=0644
	assert.Equal(t, fmt.Sprintf("/sys/fs/cgroup/raind/%s/memory.max", containerId), mockKernelSyscall.writeFileName)
	assert.Equal(t, []byte(strconv.FormatInt(int64(536870912), 10)+"\n"), mockKernelSyscall.writeFileData)
	assert.Equal(t, fs.FileMode(0644), mockKernelSyscall.writeFilePerm)

	// error is nil
	assert.Nil(t, err)
}

func TestContainerCgroupController_SetMemoryLimit_WriteFileError(t *testing.T) {
	// == arrange ==
	spec := buildMockSpec(t)
	containerId := "12345"
	mockKernelSyscall := &mockKernelSyscall{
		writeFileErr: errors.New("WriteFile() failed"),
	}
	containerCgroupPreparer := &containerCgroupController{
		syscallHandler: mockKernelSyscall,
	}

	// == act ==
	err := containerCgroupPreparer.setMemoryLimit(containerId, spec.LinuxSpec.Resources.Memory)

	// == assert ==
	// error is not nil
	assert.NotNil(t, err)
	assert.Equal(t, errors.New("WriteFile() failed"), err)
}

func TestContainerCgroupController_SetCpuLimit_Success(t *testing.T) {
	// == arrange ==
	spec := buildMockSpec(t)
	containerId := "12345"
	mockKernelSyscall := &mockKernelSyscall{}
	containerCgroupPreparer := &containerCgroupController{
		syscallHandler: mockKernelSyscall,
	}

	// == act ==
	err := containerCgroupPreparer.setCpuLimit(containerId, spec.LinuxSpec.Resources.Cpu)

	// == assert ==
	// WriteFile() is called
	assert.True(t, mockKernelSyscall.writeFileCallFlag)
	// WriteFie() parameter:
	//   path=/sys/fs/cgroup/raind/<container-id>/cpu.max
	//   data="100000 80000"
	//   mode=0644
	assert.Equal(t, fmt.Sprintf("/sys/fs/cgroup/raind/%s/cpu.max", containerId), mockKernelSyscall.writeFileName)
	assert.Equal(t, []byte(fmt.Sprintf("%d %d\n", spec.LinuxSpec.Resources.Cpu.Quota, spec.LinuxSpec.Resources.Cpu.Period)), mockKernelSyscall.writeFileData)
	assert.Equal(t, fs.FileMode(0644), mockKernelSyscall.writeFilePerm)

	// error is nil
	assert.Nil(t, err)
}

func TestContainerCgroupController_SetCpuLimit_WriteFileError(t *testing.T) {
	// == arrange ==
	spec := buildMockSpec(t)
	containerId := "12345"
	mockKernelSyscall := &mockKernelSyscall{
		writeFileErr: errors.New("WriteFile() failed"),
	}
	containerCgroupPreparer := &containerCgroupController{
		syscallHandler: mockKernelSyscall,
	}

	// == act ==
	err := containerCgroupPreparer.setCpuLimit(containerId, spec.LinuxSpec.Resources.Cpu)

	// == assert ==
	// error is not nil
	assert.NotNil(t, err)
	assert.Equal(t, errors.New("WriteFile() failed"), err)
}

func TestContainerCgroupController_SetProcessToCgroup_Success(t *testing.T) {
	// == arrange ==
	containerId := "12345"
	initPid := 11111
	mockKernelSyscall := &mockKernelSyscall{}
	containerCgroupPreparer := &containerCgroupController{
		syscallHandler: mockKernelSyscall,
	}

	// == act ==
	err := containerCgroupPreparer.setProcessToCgroup(containerId, initPid)

	// == assert ==
	// WriteFile() is called
	assert.True(t, mockKernelSyscall.writeFileCallFlag)
	// WriteFie() parameter:
	//   path=/sys/fs/cgroup/raind/<container-id>/cgroup.procs
	//   data=<pid>
	//   mode=0644
	assert.Equal(t, fmt.Sprintf("/sys/fs/cgroup/raind/%s/cgroup.procs", containerId), mockKernelSyscall.writeFileName)
	assert.Equal(t, []byte(fmt.Sprintf("%d\n", initPid)), mockKernelSyscall.writeFileData)
	assert.Equal(t, fs.FileMode(0644), mockKernelSyscall.writeFilePerm)

	// error is nil
	assert.Nil(t, err)
}

func TestContainerCgroupController_SetProcessToCgroup_WriteFileError(t *testing.T) {
	// == arrange ==
	containerId := "12345"
	initPid := 11111
	mockKernelSyscall := &mockKernelSyscall{
		writeFileErr: errors.New("WriteFile() failed"),
	}
	containerCgroupPreparer := &containerCgroupController{
		syscallHandler: mockKernelSyscall,
	}

	// == act ==
	err := containerCgroupPreparer.setProcessToCgroup(containerId, initPid)

	// == assert ==
	// error is not nil
	assert.NotNil(t, err)
	assert.Equal(t, errors.New("WriteFile() failed"), err)
}
