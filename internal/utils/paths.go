package utils

import (
	"os"
	"path/filepath"
)

const cgroupRootDir = "/sys/fs/cgroup/raind"

func DefaultRootDir() string {
	if v := os.Getenv("RAIND_ROOT_DIR"); v != "" {
		return v
	}
	return "/etc/raind/container"
}

// directory for each container
//
//	e.g. /etc/raind/container/<container-id>
func ContainerDir(containerId string) string {
	return filepath.Join(DefaultRootDir(), containerId)
}

// config.json path
//
//	e.g. /etc/raind/container/<container-id>/config.json
func ConfigFilePath(containerId string) string {
	return filepath.Join(ContainerDir(containerId), "config.json")
}

// state path
func ContainerStatePath(containerId string) string {
	return filepath.Join(ContainerDir(containerId), "state.json")
}

// fifo path
//
//	e.g. /etc/raind/container/<container-id>/exec.fifo
func FifoPath(containerId string) string {
	return filepath.Join(ContainerDir(containerId), "exec.fifo")
}

func SockPath(containerId string) string {
	return filepath.Join(ContainerDir(containerId), "tty.sock")
}

// InitPidFilePath returns pidfile path under container dir.
func InitPidFilePath(containerId string) string {
	return filepath.Join(ContainerDir(containerId), "init.pid")
}

// cgroup path
func CgroupPath(containerId string) string {
	return filepath.Join(cgroupRootDir, containerId)
}

func ShimLogPath(id string) string {
	return filepath.Join(ContainerDir(id), "shim.log")
}
func ConsoleLogPath(id string) string {
	return filepath.Join(ContainerDir(id), "console.log")
}
func InitLogPath(id string) string {
	return filepath.Join(ContainerDir(id), "init.log")
}
