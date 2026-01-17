package container

import (
	"droplet/internal/spec"
)

type mockContainerInitExecutor struct {
	// executeInit()
	executeInitCallFlag    bool
	executeInitContainerId string
	executeInitSpec        spec.Spec
	executeInitFifo        string
	executeInitPid         int
	executeInitErr         error

	// executeShim()
	executeShimCallFlag    bool
	executeShimContainerId string
	executeShimSpec        spec.Spec
	executeShimFifo        string
	executeShimPid         int
	executeShimErr         error
}

func (m *mockContainerInitExecutor) executeInit(containerId string, spec spec.Spec, fifo string) (int, error) {
	m.executeInitCallFlag = true
	m.executeInitContainerId = containerId
	m.executeInitSpec = spec
	m.executeInitFifo = fifo
	return m.executeInitPid, m.executeInitErr
}

func (m *mockContainerInitExecutor) executeShim(containerId string, spec spec.Spec, fifo string) (int, error) {
	m.executeShimCallFlag = true
	m.executeShimContainerId = containerId
	m.executeShimSpec = spec
	m.executeShimFifo = fifo
	return m.executeShimPid, m.executeShimErr
}
