package command

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

func commandKill() *cli.Command {
	return &cli.Command{
		Name:      "kill",
		Usage:     "kill a container",
		ArgsUsage: "<container-id> [signal]",
		Action: func(ctx *cli.Context) error {
			// validate args
			if ctx.NArg() < 1 || ctx.NArg() > 2 {
				return fmt.Errorf("usage: droplet kill <container-id> [signal]")
			}

			// get args
			container_id := ctx.Args().Get(0)
			signal := ctx.Args().Get(1)
			if signal == "" {
				signal = "SIGTERM"
			}

			fmt.Println("kill container: " + container_id + " by " + signal)

			return nil
		},
	}
}
