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
	updateStatusErr         error

	// GetPidFromId()
	getPidFromIdCallFlag    bool
	getPirFromIdContainerId string
	getPidFromIdPid         int
	getPidFromIdErr         error

	// GetStatusFromId()
	getStatusFromIdCallFlag    bool
	getStatusFromIdContainerId string
	getStatusFromIdStatus      status.ContainerStatus
	getStatusFromIdErr         error
}

func (m *mockStatusHandler) CreateStatusFile(containerId string, pid int, status status.ContainerStatus, rootfs string, bundle string, annotation spec.AnnotationObject) error {
	m.createStatusFileCallFlag = true
	m.createStatusFileContainerId = containerId
	m.createStatusFilePid = pid
	m.createStatusFileStatus = status
	m.createStatusFileBundle = bundle
	return m.createStatusFileErr
}

func (m *mockStatusHandler) ReadStatusFile(containerId string) (string, error) {
	m.readStatusFileCallFlag = true
	m.readStatusContainerId = containerId
	return m.readStatusDataStr, m.readStatusErr
}

func (m *mockStatusHandler) UpdateStatus(containerId string, status status.ContainerStatus, pid int) error {
	m.updateStatusCallFlag = true
	m.updateStatusContainerId = containerId
	m.updateStatusStatus = status
	m.updateStatusPid = pid
	return m.updateStatusErr
}

func (m *mockStatusHandler) GetPidFromId(containerId string) (int, error) {
	m.getPidFromIdCallFlag = true
	m.getPirFromIdContainerId = containerId
	return m.getPidFromIdPid, m.getPidFromIdErr
}

func (m *mockStatusHandler) GetStatusFromId(containerId string) (status.ContainerStatus, error) {
	m.getStatusFromIdCallFlag = true
	m.getStatusFromIdContainerId = containerId
	return m.getStatusFromIdStatus, m.getStatusFromIdErr
}
