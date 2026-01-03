package container

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultRootDir_DefaultReturn_Success(t *testing.T) {
	// == arrange ==

	// == act ==
	got := defaultRootDir()

	// == assert ==
	assert.Equal(t, "/etc/raind/container", got)
}

func TestDefaultRootDir_EnvSet_Success(t *testing.T) {
	// == arrange ==
	t.Setenv("RAIND_ROOT_DIR", "/path/to/root")

	// == act ==
	got := defaultRootDir()

	// == assert ==
	assert.Equal(t, "/path/to/root", got)
}

func TestContainerDir_Success(t *testing.T) {
	// == arrange ==
	containerId := "12345"

	// == act ==
	got := containerDir(containerId)

	// == assert ==
	assert.Equal(t, "/etc/raind/container/12345", got)
}

func TestConfigFilePath_Success(t *testing.T) {
	// == arrange ==
	containerId := "12345"

	// == act ==
	got := configFilePath(containerId)

	// == assert ==
	assert.Equal(t, "/etc/raind/container/12345/config.json", got)
}

func TestFifoPath_Success(t *testing.T) {
	// == arrange ==
	containerId := "12345"

	// == act ==
	got := fifoPath(containerId)

	// == assert ==
	assert.Equal(t, "/etc/raind/container/12345/exec.fifo", got)
}
