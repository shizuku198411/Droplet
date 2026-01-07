package container

import (
	"droplet/internal/spec"
	"droplet/internal/utils"
	"fmt"
)

func newContainerNetworkController() *containerNetworkController {
	return &containerNetworkController{
		commandFactory: &utils.ExecCommandFactory{},
	}
}

type containerNetworkPreparer interface {
	prepare(containerId string, pid int, annotation spec.AnnotationObject) error
}

type containerNetworkController struct {
	commandFactory utils.CommandFactory
}

func (c *containerNetworkController) prepare(containerId string, pid int, annotation spec.AnnotationObject) error {
	// retrieve network config from annotation
	var networkConfig spec.NetConfigObject
	if err := utils.StringToJson(annotation.Net, &networkConfig); err != nil {
		return err
	}

	// 1. create veth pair
	if err := c.createVethPair(containerId, pid, networkConfig); err != nil {
		return err
	}

	// 2. setup inside container
	if err := c.setupContainerNetns(pid, networkConfig); err != nil {
		return err
	}
	return nil
}

func (c *containerNetworkController) createVethPair(containerId string, pid int, networkConfig spec.NetConfigObject) error {
	pidStr := fmt.Sprint(pid)

	// 1. create veth
	createVeth := c.commandFactory.Command("ip", "link", "add", "name", "raind"+pidStr, "type", "veth", "peer", "name", networkConfig.HostInterface, "netns", pidStr)
	if err := createVeth.Run(); err != nil {
		return err
	}

	// 2. attach veth to bridge
	attacheVeth := c.commandFactory.Command("ip", "link", "set", "raind"+pidStr, "master", networkConfig.BridgeInterface)
	if err := attacheVeth.Run(); err != nil {
		return err
	}

	// 3. up veth
	upVeth := c.commandFactory.Command("ip", "link", "set", "raind"+pidStr, "up")
	if err := upVeth.Run(); err != nil {
		return err
	}
	return nil
}

func (c *containerNetworkController) setupContainerNetns(pid int, networkConfig spec.NetConfigObject) error {
	pidStr := fmt.Sprint(pid)

	// 1. up loopback i/f
	upLoopbackIf := c.commandFactory.Command("nsenter", "-t", pidStr, "-n", "ip", "link", "set", "lo", "up")
	if err := upLoopbackIf.Run(); err != nil {
		return err
	}

	// 2. rename veth
	renameVeth := c.commandFactory.Command("nsenter", "-t", pidStr, "-n", "ip", "link", "set", networkConfig.HostInterface, "name", networkConfig.Interface.Name)
	if err := renameVeth.Run(); err != nil {
		return err
	}

	// 3. assign address
	assignAddr := c.commandFactory.Command("nsenter", "-t", pidStr, "-n", "ip", "addr", "add", networkConfig.Interface.IPv4.Address, "dev", networkConfig.Interface.Name)
	if err := assignAddr.Run(); err != nil {
		return err
	}

	// 4. up veth
	upVeth := c.commandFactory.Command("nsenter", "-t", pidStr, "-n", "ip", "link", "set", networkConfig.Interface.Name, "up")
	if err := upVeth.Run(); err != nil {
		return err
	}

	// 5. set gateway
	setGateway := c.commandFactory.Command("nsenter", "-t", pidStr, "-n", "ip", "route", "add", "default", "via", networkConfig.Interface.IPv4.Gateway)
	if err := setGateway.Run(); err != nil {
		return err
	}

	return nil
}
