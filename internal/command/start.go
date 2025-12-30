package command

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

func commandStart() *cli.Command {
	return &cli.Command{
		Name:      "start",
		Usage:     "start a container",
		ArgsUsage: "<container-id>",
		Action: func(ctx *cli.Context) error {
			// validate args
			if ctx.NArg() != 1 {
				return fmt.Errorf("usage: droplet start <container-id>")
			}

			// get args
			container_id := ctx.Args().Get(0)

			fmt.Println("start container: " + container_id)

			return nil
		},
	}
}
