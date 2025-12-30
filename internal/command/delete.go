package command

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

func commandDelete() *cli.Command {
	return &cli.Command{
		Name:      "delete",
		Usage:     "delete a container",
		ArgsUsage: "<container-id>",
		Action: func(ctx *cli.Context) error {
			// validate args
			if ctx.NArg() != 1 {
				return fmt.Errorf("usage: droplet delete <container-id>")
			}

			// get args
			container_id := ctx.Args().Get(0)

			fmt.Println("delete container: " + container_id)

			return nil
		},
	}
}
