package command

import (
	"testing"

	"github.com/urfave/cli/v2"
)

func TestCommandStart_1(t *testing.T) {
	// test case: droplet start test-container

	// create application
	app := &cli.App{
		Name:     "droplet",
		Commands: []*cli.Command{commandStart()},
	}

	// execute command
	result := readStdout(t, func() {
		if err := app.Run([]string{"droplet", "start", "test-container"}); err != nil {
			t.Errorf("error")
		}
	})

	// validate result
	expected := "start container: test-container\n"
	if result != expected {
		t.Errorf("TEST FAIL: expected = %q, result = %q", expected, result)
	}
}
