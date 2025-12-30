package command

import (
	"testing"

	"github.com/urfave/cli/v2"
)

func TestCommandInit_1(t *testing.T) {
	// test case: droplet init

	// create application
	app := &cli.App{
		Name:     "droplet",
		Commands: []*cli.Command{commandInit()},
	}

	// execute command
	result := readStdout(t, func() {
		if err := app.Run([]string{"droplet", "init"}); err != nil {
			t.Errorf("error")
		}
	})

	// validate result
	expected := "init container\n"
	if result != expected {
		t.Errorf("TEST FAIL: expected = %q, result = %q", expected, result)
	}
}
