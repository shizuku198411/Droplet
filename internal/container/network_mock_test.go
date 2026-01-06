package container

import (
	"droplet/internal/spec"
)

type mockContainerNetworkController struct {
	// prepare()
	prepareCallFlag    bool
	prepareContainerId string
	preparePid         int
	prepareAnnotaion   spec.AnnotationObject
	prepareErr         error
}

func (m *mockContainerNetworkController) prepare(containerId string, pid int, annotation spec.AnnotationObject) error {
	m.prepareCallFlag = true
	m.prepareContainerId = containerId
	m.preparePid = pid
	m.prepareAnnotaion = annotation
	return m.prepareErr
}
