package container

import (
	"path/filepath"
)

// raind container root directory
const rootDir = "/etc/raind/container"

// directory for each container
//
//	e.g. /etc/raind/container/<container-id>
func containerDir(containerId string) string {
	return filepath.Join(rootDir, containerId)
}

// config.json path
//
//	e.g. /etc/raind/container/<container-id>/config.json
func configFilePath(containerId string) string {
	return filepath.Join(containerDir(containerId), "config.json")
}

// fifo path
//
//	e.g. /etc/raind/container/<container-id>/exec.fifo
func fifoPath(containerId string) string {
	return filepath.Join(containerDir(containerId), "exec.fifo")
}
