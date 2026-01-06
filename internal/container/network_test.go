package container

import (
	"droplet/internal/spec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainerNetworkController_Prepare_Success(t *testing.T) {
	// == arrange ==
	spec := buildMockSpec(t)
	containerId := "12345"
	pid := 11111
	mockExecCmd := &mockExecCmd{}
	mockExecCommandFactory := &mockExecCommandFactory{
		commandExecutor: mockExecCmd,
	}
	containerNetworkController := &containerNetworkController{
		commandFactory: mockExecCommandFactory,
	}

	// == act ==
	err := containerNetworkController.prepare(containerId, pid, spec.Annotations)

	// == assert ==
	// error is nil
	assert.Nil(t, err)
}

func TestContainerNetworkController_CreateVethPair_Success(t *testing.T) {
	// == arrange ==
	containerId := "12345"
	pid := 11111
	networkConfig := spec.NetConfigObject{
		HostInterface:   "eth0",
		BridgeInterface: "raind0",
		Interface: spec.InterfaceObject{
			Name: "eth0",
			IPv4: spec.IPv4Object{
				Address: "10.166.0.1/24",
				Gateway: "10.166.0.254",
			},
			Dns: spec.DnsObject{
				Servers: []string{
					"8.8.8.8",
				},
			},
		},
	}
	mockExecCmd := &mockExecCmd{}
	mockExecCommandFactory := &mockExecCommandFactory{
		commandExecutor: mockExecCmd,
	}
	containerNetworkController := &containerNetworkController{
		commandFactory: mockExecCommandFactory,
	}

	// == act ==
	err := containerNetworkController.createVethPair(containerId, pid, networkConfig)

	// == act ==
	// Command() call time: 3
	assert.Equal(t, 3, len(mockExecCommandFactory.commandCalls))

	// Command() call 1: ip link add name raind<pid> type veth peer name <host-i/f> netns <pid>
	commandCall1 := mockExecCommandFactory.commandCalls[0]
	assert.Equal(t, "ip", commandCall1.name)
	assert.Equal(t,
		[]string{"link", "add", "name", "raind11111", "type", "veth", "peer", "name", "eth0", "netns", "11111"},
		commandCall1.args,
	)

	// Command() call 2: ip link set raind<pid> master <bridge>
	commandCall2 := mockExecCommandFactory.commandCalls[1]
	assert.Equal(t, "ip", commandCall1.name)
	assert.Equal(t,
		[]string{"link", "set", "raind11111", "master", "raind0"},
		commandCall2.args,
	)

	// Command() call 3: ip link set raind<pid> up
	commandCall3 := mockExecCommandFactory.commandCalls[2]
	assert.Equal(t, "ip", commandCall3.name)
	assert.Equal(t,
		[]string{"link", "set", "raind11111", "up"},
		commandCall3.args,
	)

	// error is nil
	assert.Nil(t, err)
}

func TestContainerNetworkController_SetupContainerNetns_Success(t *testing.T) {
	// == arrange ==
	pid := 11111
	networkConfig := spec.NetConfigObject{
		HostInterface:   "eth0",
		BridgeInterface: "raind0",
		Interface: spec.InterfaceObject{
			Name: "eth0",
			IPv4: spec.IPv4Object{
				Address: "10.166.0.1/24",
				Gateway: "10.166.0.254",
			},
			Dns: spec.DnsObject{
				Servers: []string{
					"8.8.8.8",
				},
			},
		},
	}
	mockExecCmd := &mockExecCmd{}
	mockExecCommandFactory := &mockExecCommandFactory{
		commandExecutor: mockExecCmd,
	}
	containerNetworkController := &containerNetworkController{
		commandFactory: mockExecCommandFactory,
	}

	// == act ==
	err := containerNetworkController.setupContainerNetns(pid, networkConfig)

	// == act ==
	// Command() call time: 5
	assert.Equal(t, 5, len(mockExecCommandFactory.commandCalls))

	// Command() call 1: nsenter -t <pid> -n ip link set lo up
	commandCall1 := mockExecCommandFactory.commandCalls[0]
	assert.Equal(t, "nsenter", commandCall1.name)
	assert.Equal(t,
		[]string{"-t", "11111", "-n", "ip", "link", "set", "lo", "up"},
		commandCall1.args,
	)

	// Command() call 2: nsetner -t <pid> -n ip link set <host-i/f> name <container-i/f>
	commandCall2 := mockExecCommandFactory.commandCalls[1]
	assert.Equal(t, "nsenter", commandCall1.name)
	assert.Equal(t,
		[]string{"-t", "11111", "-n", "ip", "link", "set", "eth0", "name", "eth0"},
		commandCall2.args,
	)

	// Command() call 3: nsetner -t <pid> -n ip addr add <address> dev <container-i/f>
	commandCall3 := mockExecCommandFactory.commandCalls[2]
	assert.Equal(t, "nsenter", commandCall3.name)
	assert.Equal(t,
		[]string{"-t", "11111", "-n", "ip", "addr", "add", "10.166.0.1/24", "dev", "eth0"},
		commandCall3.args,
	)

	// Command() call 4: nsetner -t <pid> -n ip link set <container-i/f> up
	commandCall4 := mockExecCommandFactory.commandCalls[3]
	assert.Equal(t, "nsenter", commandCall4.name)
	assert.Equal(t,
		[]string{"-t", "11111", "-n", "ip", "link", "set", "eth0", "up"},
		commandCall4.args,
	)

	// Command() call 5: nsetner -t <pid> -n ip route add default via <gateway>
	commandCall5 := mockExecCommandFactory.commandCalls[4]
	assert.Equal(t, "nsenter", commandCall5.name)
	assert.Equal(t,
		[]string{"-t", "11111", "-n", "ip", "route", "add", "default", "via", "10.166.0.254"},
		commandCall5.args,
	)

	// error is nil
	assert.Nil(t, err)
}
