package command

import (
	"testing"

	"github.com/urfave/cli/v2"
)

func TestCommandDelete_1(t *testing.T) {
	// test case: droplet delete test-container

	// create application
	app := &cli.App{
		Name:     "droplet",
		Commands: []*cli.Command{commandDelete()},
	}

	// execute command
	result := readStdout(t, func() {
		if err := app.Run([]string{"droplet", "delete", "test-container"}); err != nil {
			t.Errorf("error")
		}
	})

	// validate result
	expected := "delete container: test-container\n"
	if result != expected {
		t.Errorf("TEST FAIL: expected = %q, result = %q", expected, result)
	}
}
