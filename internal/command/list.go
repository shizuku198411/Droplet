package command

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

func commandList() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "lists containers",
		Action: func(ctx *cli.Context) error {
			// validate args
			if ctx.NArg() != 0 {
				return fmt.Errorf("usage: droplet list")
			}

			fmt.Println("container list")

			return nil
		},
	}
}
