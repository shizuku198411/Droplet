package command

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

func commandExec() *cli.Command {
	return &cli.Command{
		Name:      "exec",
		Usage:     "exec a container",
		ArgsUsage: "<container-id> <command>",
		Action: func(ctx *cli.Context) error {
			// validate args
			if ctx.NArg() < 1 || ctx.NArg() > 2 {
				return fmt.Errorf("usage: droplet exec <container-id> <command>")
			}

			// get args
			container_id := ctx.Args().Get(0)
			command := ctx.Args().Get(1)

			fmt.Println("exec command: " + command + " in container: " + container_id)

			return nil
		},
	}
}
