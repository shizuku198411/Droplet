package container

import (
	"droplet/internal/spec"
)

type mockRootContainerEnvPreparer struct {
	// prepare()
	prepareCallFlag    bool
	prepareContainerId string
	prepareSpec        spec.Spec
	prepareErr         error
}

func (m *mockRootContainerEnvPreparer) prepare(containerId string, spec spec.Spec) error {
	m.prepareCallFlag = true
	m.prepareContainerId = containerId
	m.prepareSpec = spec
	return m.prepareErr
}
