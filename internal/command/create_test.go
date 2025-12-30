package command

import (
	"testing"

	"github.com/urfave/cli/v2"
)

func TestCommandCreate_1(t *testing.T) {
	// test case: droplet create test-container

	// create application
	app := &cli.App{
		Name:     "droplet",
		Commands: []*cli.Command{commandCreate()},
	}

	// execute command
	result := readStdout(t, func() {
		if err := app.Run([]string{"droplet", "create", "test-container"}); err != nil {
			t.Errorf("error")
		}
	})

	// validate result
	expected := "test-container\n.\n"
	if result != expected {
		t.Errorf("TEST FAIL: expected = %q, result = %q", expected, result)
	}
}

func TestCommandCreate_2(t *testing.T) {
	// test case: droplet create test-container /tmp

	// create application
	app := &cli.App{
		Name:     "droplet",
		Commands: []*cli.Command{commandCreate()},
	}

	// execute command
	result := readStdout(t, func() {
		if err := app.Run([]string{"droplet", "create", "test-container", "/tmp"}); err != nil {
			t.Errorf("error")
		}
	})

	// validate result
	expected := "test-container\n/tmp\n"
	if result != expected {
		t.Errorf("TEST FAIL: expected = %q, result = %q", expected, result)
	}
}
