package command

import (
	"testing"

	"github.com/urfave/cli/v2"
)

func TestCommandList_1(t *testing.T) {
	// test case: droplet list

	// create application
	app := &cli.App{
		Name:     "droplet",
		Commands: []*cli.Command{commandList()},
	}

	// execute command
	result := readStdout(t, func() {
		if err := app.Run([]string{"droplet", "list"}); err != nil {
			t.Errorf("error")
		}
	})

	// validate result
	expected := "container list\n"
	if result != expected {
		t.Errorf("TEST FAIL: expected = %q, result = %q", expected, result)
	}
}
