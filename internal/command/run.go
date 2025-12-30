package command

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

func commandRun() *cli.Command {
	return &cli.Command{
		Name:      "run",
		Usage:     "run a container",
		ArgsUsage: "<container-id> [path-to-bundle]",
		Action: func(ctx *cli.Context) error {
			// validate args
			if ctx.NArg() < 1 || ctx.NArg() > 2 {
				return fmt.Errorf("usage: droplet run <container-id> [path-to-bundle]")
			}

			// get args
			container_id := ctx.Args().Get(0)
			bundle_path := ctx.Args().Get(1)
			if bundle_path == "" {
				bundle_path = "."
			}

			fmt.Println("run container: " + container_id + ", path: " + bundle_path)

			return nil
		},
	}
}
