package command

import (
	"droplet/internal/container"

	"github.com/urfave/cli/v2"
)

func commandInit() *cli.Command {
	return &cli.Command{
		Name:      "init",
		Usage:     "initialize a container",
		ArgsUsage: "<fifo-path> <entrypoint>",
		Hidden:    true,
		Action:    runInit,
	}
}

func runInit(ctx *cli.Context) error {
	// retrieve fifo and entrypoint
	fifo := ctx.Args().Get(0)
	args := ctx.Args().Slice()
	entrypoint := args[1:]

	err := container.InitContainer(container.InitOption{
		Fifo:       fifo,
		Entrypoint: entrypoint,
	})
	if err != nil {
		return err
	}

	return nil
}
