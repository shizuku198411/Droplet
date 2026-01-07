package status

import (
	"droplet/internal/oci"
	"droplet/internal/spec"
	"droplet/internal/utils"
	"os"
	"syscall"
)

type ContainerStatusManager interface {
	CreateStatusFile(containerId string, pid int, status ContainerStatus, rootfs string, bundle string, annotation spec.AnnotationObject) error
	ReadStatusFile(containerId string) (string, error)
	UpdateStatus(containerId string, status ContainerStatus, pid int) error
	GetPidFromId(containerId string) (int, error)
	GetStatusFromId(containerId string) (ContainerStatus, error)
}

func NewStatusHandler() *StatusHandler {
	return &StatusHandler{
		processManager: NewProcessHandler(),
	}
}

type StatusHandler struct {
	processManager ProcessManager
}

func (h *StatusHandler) CreateStatusFile(containerId string, pid int, status ContainerStatus,
	rootfs string, bundle string, annotation spec.AnnotationObject) error {
	stateFilePath := utils.ContainerStatePath(containerId)
	statusObject := StatusObject{
		OciVersion: oci.OCIVersion,
		Id:         containerId,
		Status:     status.String(),
		Pid:        pid,
		Rootfs:     rootfs,
		Bundle:     bundle,
		Annotaion:  annotation,
	}

	if err := utils.WriteJsonToFile(stateFilePath, statusObject); err != nil {
		return err
	}

	return nil
}

func (h *StatusHandler) ReadStatusFile(containerId string) (string, error) {
	stateFilePath := utils.ContainerStatePath(containerId)
	// load status file
	var statusObject StatusObject
	if err := utils.ReadJsonFile(stateFilePath, &statusObject); err != nil {
		return "", err
	}
	pid := statusObject.Pid
	currentStatus, parseErr := ParseContainerStatus(statusObject.Status)
	if parseErr != nil {
		return "", parseErr
	}

	// recompute status
	if err := h.recomputeStatus(containerId, pid, currentStatus); err != nil {
		return "", err
	}

	// read file
	data, err := os.ReadFile(stateFilePath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (h *StatusHandler) UpdateStatus(containerId string, status ContainerStatus, pid int) error {
	stateFilePath := utils.ContainerStatePath(containerId)
	// load status file
	var statusObject StatusObject
	if err := utils.ReadJsonFile(stateFilePath, &statusObject); err != nil {
		return err
	}

	// update
	if status >= 0 && status <= 3 {
		statusObject.Status = status.String()
	}
	if pid >= 0 {
		statusObject.Pid = pid
	}

	// write status file
	if err := utils.WriteJsonToFile(stateFilePath, statusObject); err != nil {
		return err
	}

	return nil
}

func (h *StatusHandler) GetPidFromId(containerId string) (int, error) {
	stateFilePath := utils.ContainerStatePath(containerId)
	// load status file
	var statusObject StatusObject
	if err := utils.ReadJsonFile(stateFilePath, &statusObject); err != nil {
		return -1, err
	}

	return statusObject.Pid, nil
}

func (h *StatusHandler) GetStatusFromId(containerId string) (ContainerStatus, error) {
	stateFilePath := utils.ContainerStatePath(containerId)
	// load status file
	var statusObject StatusObject
	if err := utils.ReadJsonFile(stateFilePath, &statusObject); err != nil {
		return -1, err
	}
	pid := statusObject.Pid
	currentStatus, parseErr := ParseContainerStatus(statusObject.Status)
	if parseErr != nil {
		return -1, parseErr
	}

	// recompute status
	if err := h.recomputeStatus(containerId, pid, currentStatus); err != nil {
		return -1, err
	}

	if err := utils.ReadJsonFile(stateFilePath, &statusObject); err != nil {
		return -1, err
	}
	newStatus, newParseErr := ParseContainerStatus(statusObject.Status)
	if newParseErr != nil {
		return -1, newParseErr
	}

	return newStatus, nil
}

func (h *StatusHandler) recomputeStatus(containerId string, pid int, currentStatus ContainerStatus) error {
	if currentStatus == RUNNING {
		alive, _ := h.pidAlive(pid)
		if !alive {
			if err := h.UpdateStatus(containerId, STOPPED, 0); err != nil {
				return err
			}
		}
	}
	return nil
}

func (h *StatusHandler) pidAlive(pid int) (bool, error) {
	if pid <= 0 {
		// process not exist
		return false, nil
	}

	// send 0 signal to process
	err := h.processManager.Kill(pid, 0)
	if err == nil {
		// process exist
		return true, nil
	}
	if err == syscall.ESRCH {
		// no such process
		return false, nil
	}
	if err == syscall.EPERM {
		// operation not permitted, but process exist
		return true, nil
	}

	return false, nil
}

type ProcessManager interface {
	Kill(pid int, sig syscall.Signal) error
}

func NewProcessHandler() *ProcessHandler {
	return &ProcessHandler{}
}

type ProcessHandler struct{}

func (h *ProcessHandler) Kill(pid int, sig syscall.Signal) error {
	return syscall.Kill(pid, sig)
}
