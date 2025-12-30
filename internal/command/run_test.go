package command

import (
	"testing"

	"github.com/urfave/cli/v2"
)

func TestCommandRun_1(t *testing.T) {
	// test case: droplet run test-container

	// create application
	app := &cli.App{
		Name:     "droplet",
		Commands: []*cli.Command{commandRun()},
	}

	// execute command
	result := readStdout(t, func() {
		if err := app.Run([]string{"droplet", "run", "test-container"}); err != nil {
			t.Errorf("error")
		}
	})

	// validate result
	expected := "run container: test-container, path: .\n"
	if result != expected {
		t.Errorf("TEST FAIL: expected = %q, result = %q", expected, result)
	}
}

func TestCommandRun_2(t *testing.T) {
	// test case: droplet run test-container /tmp

	// create application
	app := &cli.App{
		Name:     "droplet",
		Commands: []*cli.Command{commandRun()},
	}

	// execute command
	result := readStdout(t, func() {
		if err := app.Run([]string{"droplet", "run", "test-container", "/tmp"}); err != nil {
			t.Errorf("error")
		}
	})

	// validate result
	expected := "run container: test-container, path: /tmp\n"
	if result != expected {
		t.Errorf("TEST FAIL: expected = %q, result = %q", expected, result)
	}
}
