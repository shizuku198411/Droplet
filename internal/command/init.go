package command

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

func commandInit() *cli.Command {
	return &cli.Command{
		Name:   "init",
		Usage:  "initialize a container",
		Hidden: true,
		Action: func(ctx *cli.Context) error {
			// validate args
			if ctx.NArg() != 0 {
				return fmt.Errorf("usage: droplet init")
			}

			fmt.Println("init container")

			return nil
		},
	}
}
