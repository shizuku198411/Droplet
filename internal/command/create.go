package command

import (
	"droplet/internal/container"

	"github.com/urfave/cli/v2"
)

func commandCreate() *cli.Command {
	return &cli.Command{
		Name:      "create",
		Usage:     "create a container",
		ArgsUsage: "<container-id>",
		Action:    runCreate,
	}
}

func runCreate(ctx *cli.Context) error {
	// retrieve container ID
	containerId := ctx.Args().Get(0)

	err := container.CreateContainer(container.CreateOption{ContainerId: containerId})

	if err != nil {
		return err
	}

	return nil
}
