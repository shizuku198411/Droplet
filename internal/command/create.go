package command

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

func commandCreate() *cli.Command {
	return &cli.Command{
		Name:      "create",
		Usage:     "create a container",
		ArgsUsage: "<container-id> [path-to-bundle]",
		Action: func(ctx *cli.Context) error {
			// validate args
			if ctx.NArg() < 1 || ctx.NArg() > 2 {
				return fmt.Errorf("usage: droplet create <container-id> [path-to-bundle]")
			}

			// get args
			container_id := ctx.Args().Get(0)
			bundle_path := ctx.Args().Get(1)
			if bundle_path == "" {
				bundle_path = "."
			}

			fmt.Println(container_id)
			fmt.Println(bundle_path)

			return nil
		},
	}
}
