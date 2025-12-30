package command

import (
	"testing"

	"github.com/urfave/cli/v2"
)

func TestCommandExec_1(t *testing.T) {
	// test case: droplet exec test-container /bin/bash

	// create application
	app := &cli.App{
		Name:     "droplet",
		Commands: []*cli.Command{commandExec()},
	}

	// execute command
	result := readStdout(t, func() {
		if err := app.Run([]string{"droplet", "exec", "test-container", "/bin/bash"}); err != nil {
			t.Errorf("error")
		}
	})

	// validate result
	expected := "exec command: /bin/bash in container: test-container\n"
	if result != expected {
		t.Errorf("TEST FAIL: expected = %q, result = %q", expected, result)
	}
}
