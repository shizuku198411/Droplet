package command

import (
	"testing"

	"github.com/urfave/cli/v2"
)

func TestCommandKill_1(t *testing.T) {
	// test case: droplet kill test-container

	// create application
	app := &cli.App{
		Name:     "droplet",
		Commands: []*cli.Command{commandKill()},
	}

	// execute command
	result := readStdout(t, func() {
		if err := app.Run([]string{"droplet", "kill", "test-container"}); err != nil {
			t.Errorf("error")
		}
	})

	// validate result
	expected := "kill container: test-container by SIGTERM\n"
	if result != expected {
		t.Errorf("TEST FAIL: expected = %q, result = %q", expected, result)
	}
}

func TestCommandKill_2(t *testing.T) {
	// test case: droplet kill test-container SIGKILL

	// create application
	app := &cli.App{
		Name:     "droplet",
		Commands: []*cli.Command{commandKill()},
	}

	// execute command
	result := readStdout(t, func() {
		if err := app.Run([]string{"droplet", "kill", "test-container", "SIGKILL"}); err != nil {
			t.Errorf("error")
		}
	})

	// validate result
	expected := "kill container: test-container by SIGKILL\n"
	if result != expected {
		t.Errorf("TEST FAIL: expected = %q, result = %q", expected, result)
	}
}
