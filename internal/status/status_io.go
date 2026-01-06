package status

import (
	"droplet/internal/utils"
	"path/filepath"
)

type ContainerStatusManager interface {
	CreateStatusFile(path string, containerId string, pid int, status ContainerStatus, bundle string) error
	UpdateStatus(path string, containerId string, status ContainerStatus, pid int) error
}

func NewStatusHandler() *StatusHandler {
	return &StatusHandler{}
}

type StatusHandler struct{}

func (h *StatusHandler) CreateStatusFile(path string, containerId string, pid int, status ContainerStatus, bundle string) error {
	stateFilePath := filepath.Join(path, "state.json")
	statusObject := StatusObject{
		OciVersion: "1.3.0",
		Id:         containerId,
		Status:     status.String(),
		Pid:        pid,
		Bundle:     bundle,
	}

	if err := utils.WriteJsonToFile(stateFilePath, statusObject); err != nil {
		return err
	}

	return nil
}

func (h *StatusHandler) UpdateStatus(path string, containerId string, status ContainerStatus, pid int) error {
	stateFilePath := filepath.Join(path, "state.json")
	// load status file
	var statusObject StatusObject
	if err := utils.ReadJsonFile(stateFilePath, &statusObject); err != nil {
		return err
	}

	// update
	if status >= 0 && status <= 3 {
		statusObject.Status = status.String()
	}
	if pid > 0 {
		statusObject.Pid = pid
	}
	statusObject.Status = status.String()

	// write status file
	if err := utils.WriteJsonToFile(stateFilePath, statusObject); err != nil {
		return err
	}

	return nil
}
