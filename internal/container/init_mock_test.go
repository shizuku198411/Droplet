package container

import (
	"droplet/internal/spec"
)

type mockRootContainerEnvPreparer struct {
	// prepare()
	prepareCallFlag bool
	prepareSpec     spec.Spec
	prepareErr      error
}

func (m *mockRootContainerEnvPreparer) prepare(spec spec.Spec) error {
	m.prepareCallFlag = true
	m.prepareSpec = spec
	return m.prepareErr
}
