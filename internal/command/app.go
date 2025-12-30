package command

import (
	"github.com/urfave/cli/v2"
)

func NewApp() *cli.App {
	app := &cli.App{
		Name:  "droplet",
		Usage: "low-level container runtime",
		Commands: []*cli.Command{
			commandCreate(),
			commandStart(),
			commandKill(),
			commandDelete(),
			commandState(),
			commandRun(),
			commandExec(),
			commandSpec(),
			commandList(),
			commandInit(),
		},
	}

	// disable slice flag separator
	app.DisableSliceFlagSeparator = true

	return app
}
