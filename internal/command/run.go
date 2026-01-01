package command

import (
	"droplet/internal/container"

	"github.com/urfave/cli/v2"
)

func commandRun() *cli.Command {
	return &cli.Command{
		Name:      "run",
		Usage:     "run a container",
		ArgsUsage: "<container-id>",
		Action:    runRun,
	}
}

func runRun(ctx *cli.Context) error {
	// retrieve container ID
	containerId := ctx.Args().Get(0)

	containerRun := container.NewContainerRun()
	err := containerRun.Run(
		container.RunOption{
			ContainerId: containerId,
		},
	)

	if err != nil {
		return err
	}

	return nil
}
