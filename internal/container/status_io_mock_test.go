package container

import (
	"droplet/internal/spec"
	"droplet/internal/status"
)

type mockStatusHandler struct {
	// CreateStatusFile()
	createStatusFileCallFlag    bool
	createStatusFileContainerId string
	createStatusFilePid         int
	createStatusFileStatus      status.ContainerStatus
	createStatusRootfs          string
	createStatusFileBundle      string
	createStatusAnnotation      spec.AnnotationObject
	createStatusFileErr         error

	// RemoveStatusFile()
	removeStatusFileCallFlag    bool
	removeStatusFileContainerId string
	removeStatusFileErr         error

	// ReadStatusFile()
	readStatusFileCallFlag bool
	readStatusContainerId  string
	readStatusDataStr      string
	readStatusErr          error

	// UpdateStatus()
	updateStatusCallFlag    bool
	updateStatusContainerId string
	updateStatusStatus      status.ContainerStatus
	updateStatusPid         int
	updateStatusShimPid     int
	updateStatusErr         error

	// GetPidFromId()
	getPidFromIdCallFlag    bool
	getPirFromIdContainerId string
	getPidFromIdPid         int
	getPidFromIdErr         error

	// GetShimPidFromId()
	getShimFromIdCallFlag    bool
	getShimFromIdContainerId string
	getShimFromIdPid         int
	getShimFromIdErr         error

	// GetStatusFromId()
	getStatusFromIdCallFlag    bool
	getStatusFromIdContainerId string
	getStatusFromIdStatus      status.ContainerStatus
	getStatusFromIdErr         error

	// ListContainers()
	listContainersCallFlag bool
	listContainersList     []status.StatusObject
	listConainersErr       error
}

func (m *mockStatusHandler) CreateStatusFile(containerId string, pid int, status status.ContainerStatus, rootfs string, bundle string, annotation spec.AnnotationObject) error {
	m.createStatusFileCallFlag = true
	m.createStatusFileContainerId = containerId
	m.createStatusFilePid = pid
	m.createStatusFileStatus = status
	m.createStatusFileBundle = bundle
	return m.createStatusFileErr
}

func (m *mockStatusHandler) RemoveStatusFile(containerId string) error {
	m.removeStatusFileCallFlag = true
	m.removeStatusFileContainerId = containerId
	return m.removeStatusFileErr
}

func (m *mockStatusHandler) ReadStatusFile(containerId string) (string, error) {
	m.readStatusFileCallFlag = true
	m.readStatusContainerId = containerId
	return m.readStatusDataStr, m.readStatusErr
}

func (m *mockStatusHandler) UpdateStatus(containerId string, status status.ContainerStatus, pid int, shimPid int) error {
	m.updateStatusCallFlag = true
	m.updateStatusContainerId = containerId
	m.updateStatusStatus = status
	m.updateStatusPid = pid
	m.updateStatusShimPid = shimPid
	return m.updateStatusErr
}

func (m *mockStatusHandler) GetPidFromId(containerId string) (int, error) {
	m.getPidFromIdCallFlag = true
	m.getPirFromIdContainerId = containerId
	return m.getPidFromIdPid, m.getPidFromIdErr
}

func (m *mockStatusHandler) GetShimPidFromId(containerId string) (int, error) {
	m.getShimFromIdCallFlag = true
	m.getShimFromIdContainerId = containerId
	return m.getShimFromIdPid, m.getShimFromIdErr
}

func (m *mockStatusHandler) GetStatusFromId(containerId string) (status.ContainerStatus, error) {
	m.getStatusFromIdCallFlag = true
	m.getStatusFromIdContainerId = containerId
	return m.getStatusFromIdStatus, m.getStatusFromIdErr
}

func (m *mockStatusHandler) ListContainers() ([]status.StatusObject, error) {
	m.listContainersCallFlag = true
	return m.listContainersList, m.listConainersErr
}
