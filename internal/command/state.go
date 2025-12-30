package command

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

func commandState() *cli.Command {
	return &cli.Command{
		Name:      "state",
		Usage:     "query state a container",
		ArgsUsage: "<container-id>",
		Action: func(ctx *cli.Context) error {
			// validate args
			if ctx.NArg() != 1 {
				return fmt.Errorf("usage: droplet state <container-id>")
			}

			// get args
			container_id := ctx.Args().Get(0)

			fmt.Println("state container: " + container_id)

			return nil
		},
	}
}
