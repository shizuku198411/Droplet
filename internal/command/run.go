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
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "interactive",
				Usage:   "Run the container in interactive mode",
				Aliases: []string{"i"},
			},
			&cli.BoolFlag{
				Name:   "print-pid",
				Hidden: true,
				Value:  false,
			},
		},
		Action: runRun,
	}
}

func runRun(ctx *cli.Context) error {
	// retrieve container ID
	containerId := ctx.Args().Get(0)
	// options
	// interactive
	interactive := ctx.Bool("interactive")
	// print-pid
	printPidFlag := ctx.Bool("print-pid")

	containerRun := container.NewContainerRun()
	err := containerRun.Run(
		container.RunOption{
			ContainerId:  containerId,
			Interactive:  interactive,
			PrintPidFlag: printPidFlag,
		},
	)

	if err != nil {
		return err
	}

	return nil
}
