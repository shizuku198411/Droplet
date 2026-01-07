package container

import (
	"droplet/internal/spec"
)

type mockHookController struct {
	// RunCreateRuntimeHooks()
	runCreateRuntimeHooksCallFlag    bool
	runCreateRuntimeHooksContainerId string
	runCreateRuntimeHooksHookList    []spec.HookObject
	runCreateRuntimeHooksErr         error

	// RunCreateContainerHooks()
	runCreateContainerHooksCallFlag    bool
	runCreateContainerHooksContainerId string
	runCreateContainerHooksHookList    []spec.HookObject
	runCreateContainerHooksErr         error

	// RunStartContainerHooks()
	runStartContainerHooksCallFlag    bool
	runStartContainerHooksContainerId string
	runStartContainerHooksHookList    []spec.HookObject
	runStartContainerHooksErr         error

	// RunPoststartHooks()
	runPoststartHooksCallFlag    bool
	runPoststartHooksContainerId string
	runPoststartHooksHookList    []spec.HookObject
	runPoststartHooksErr         error

	// RunPoststopHooks()
	runPoststopHooksCallFlag    bool
	runPoststopHooksContainerId string
	runPoststopHooksHookList    []spec.HookObject
	runPoststopHooksErr         error
}

func (m *mockHookController) RunCreateRuntimeHooks(containerId string, hookList []spec.HookObject) error {
	m.runCreateRuntimeHooksCallFlag = true
	m.runCreateRuntimeHooksContainerId = containerId
	m.runCreateRuntimeHooksHookList = hookList
	return m.runCreateRuntimeHooksErr
}

func (m *mockHookController) RunCreateContainerHooks(containerId string, hookList []spec.HookObject) error {
	m.runCreateContainerHooksCallFlag = true
	m.runCreateContainerHooksContainerId = containerId
	m.runCreateContainerHooksHookList = hookList
	return m.runCreateContainerHooksErr
}

func (m *mockHookController) RunStartContainerHooks(containerId string, hookList []spec.HookObject) error {
	m.runStartContainerHooksCallFlag = true
	m.runStartContainerHooksContainerId = containerId
	m.runStartContainerHooksHookList = hookList
	return m.runStartContainerHooksErr
}

func (m *mockHookController) RunPoststartHooks(containerId string, hookList []spec.HookObject) error {
	m.runPoststartHooksCallFlag = true
	m.runPoststartHooksContainerId = containerId
	m.runPoststartHooksHookList = hookList
	return m.runPoststartHooksErr
}

func (m *mockHookController) RunPoststopHooks(containerId string, hookList []spec.HookObject) error {
	m.runPoststopHooksCallFlag = true
	m.runPoststopHooksContainerId = containerId
	m.runPoststopHooksHookList = hookList
	return m.runPoststopHooksErr
}
