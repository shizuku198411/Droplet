package container

import (
	"droplet/internal/spec"
	"droplet/internal/utils"
	"fmt"
	"path/filepath"
	"strconv"
)

func newContainerCgroupController() *containerCgroupController {
	return &containerCgroupController{
		syscallHandler: utils.NewSyscallHandler(),
	}
}

type containerCgroupPreparer interface {
	prepare(containerId string, spec spec.Spec, pid int) error
}

type containerCgroupController struct {
	syscallHandler utils.KernelSyscallHandler
}

func (c *containerCgroupController) prepare(containerId string, spec spec.Spec, pid int) error {
	// 1. set memory limit
	if err := c.setMemoryLimit(containerId, spec.LinuxSpec.Resources.Memory); err != nil {
		return err
	}

	// 2. set cpu limit
	if err := c.setCpuLimit(containerId, spec.LinuxSpec.Resources.Cpu); err != nil {
		return err
	}

	// 3. set pid to cgroup.procs
	if err := c.setProcessToCgroup(containerId, pid); err != nil {
		return err
	}

	return nil
}

func (c *containerCgroupController) setMemoryLimit(containerId string, memoryObject spec.MemoryObject) error {
	cgroupPath := utils.CgroupPath(containerId)
	memoryPath := filepath.Join(cgroupPath, "memory.max")
	memoryLimit := strconv.FormatInt(int64(memoryObject.Limit), 10)

	if err := c.syscallHandler.WriteFile(memoryPath, []byte(memoryLimit+"\n"), 0644); err != nil {
		return err
	}

	return nil
}

func (c *containerCgroupController) setCpuLimit(containerId string, cpuObject spec.CpuObject) error {
	cgroupPath := utils.CgroupPath(containerId)
	cpuPath := filepath.Join(cgroupPath, "cpu.max")
	cpuLimit := fmt.Sprintf("%d %d\n", cpuObject.Quota, cpuObject.Period)

	if err := c.syscallHandler.WriteFile(cpuPath, []byte(cpuLimit), 0644); err != nil {
		return err
	}
	return nil
}

func (c *containerCgroupController) setProcessToCgroup(containerId string, pid int) error {
	cgroupPath := utils.CgroupPath(containerId)
	cgroupProcs := filepath.Join(cgroupPath, "cgroup.procs")
	data := strconv.Itoa(pid) + "\n"

	if err := c.syscallHandler.WriteFile(cgroupProcs, []byte(data), 0644); err != nil {
		return err
	}

	return nil
}
