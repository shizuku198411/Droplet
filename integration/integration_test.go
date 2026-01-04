//go:build integration
// +build integration

package integration

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"testing"
	"time"

	"droplet/internal/spec"

	"github.com/stretchr/testify/assert"
)

var binary string

// in integration testing, tests are executed using pre-built binary
// TestMain() is used to build the binary into a temporary directory
// when all tests have finished, the temporary directory is deleted
// regardless of whether the test cases passed or failed.
func TestMain(m *testing.M) {
	// build binary
	tmpDir, err := os.MkdirTemp("", "droplet-integration-*")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmpDir)

	binary = filepath.Join(tmpDir, "droplet")

	build := exec.Command("go", "build", "-o", binary, "../cmd/droplet")
	build.Stdout = os.Stdout
	build.Stderr = os.Stderr
	if err := build.Run(); err != nil {
		panic("go build failed: " + err.Error())
	}

	test := m.Run()
	os.Exit(test)
}

// == helper define from ==
// isLocalEnvironment is to determine whether execution is in a local environment
func isLocalEnvironment(t *testing.T) bool {
	t.Helper()

	if v := os.Getenv("RAIND_INTEGRATION_LOCAL"); v != "TRUE" {
		return false
	}
	return true
}

// determine whether sudo can be executed without a password
func isSudoNotNeedPassword(t *testing.T) bool {
	t.Helper()

	if err := exec.Command("sudo", "-n", "true"); err != nil {
		return false
	}
	return true
}

// extract PID from the value
func extractPidFromOutput(t *testing.T, v []byte) int {
	t.Helper()

	re := regexp.MustCompile(`pid:\s*(\d+)`)
	m := re.FindSubmatch(v)
	if len(m) < 2 {
		t.Fatalf("failed to parse pid from value:\n%s", string(v))
	}

	pid, err := strconv.Atoi(string(m[1]))
	if err != nil {
		t.Fatalf("invalid pid in v: %q (%v)", string(m[1]), err)
	}
	return pid
}

// determine whether a process exists
func isProcessExists(t *testing.T, pid int) bool {
	t.Helper()

	// confirm that the process was successfully terminated using `ps`
	cmd := exec.Command("ps", "-p", strconv.Itoa(pid), "-o", "pid=")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

// cleanupInitProcess sends SIGTERM/SIGKILL to the init process
// and confirms that the process has terminated
func cleanupInitProcess(t *testing.T, pid int) {
	t.Helper()

	pidStr := strconv.Itoa(pid)

	// send SIGTERM
	_ = exec.Command("sudo", "kill", "-TERM", pidStr).Run()
	// wait for killing process
	time.Sleep(500 * time.Millisecond)
	// send SIGKILL
	_ = exec.Command("sudo", "kill", "-KILL", pidStr).Run()

	// confirm that the process was successfully terminated using `ps`
	cmd := exec.Command("ps", "-p", pidStr, "-o", "pid=")
	if err := cmd.Run(); err == nil {
		t.Fatalf("init process (pid=%d) is still running after cleanup", pid)
	}
}

// == helper define end ==

func TestDropletRun_Spec_Success(t *testing.T) {
	// skip in environments where sudo cannot be used without a password.
	// if password input is possible in the local environment, run the test
	if !isLocalEnvironment(t) {
		if !isSudoNotNeedPassword(t) {
			t.Skip("sudo -n not available, skipping integration test")
		}
	}

	// == arrange ==
	path := t.TempDir()
	configJsonPath := filepath.Join(path, "config.json")
	// config.json parameter
	command := `/bin/sh -c "echo Hello World"`
	hostname := "mycontainer"
	cmd := exec.Command("sudo", binary,
		"spec",
		"--command", command,
		"--hostname", hostname,
		"--output", path,
	)

	// == act ==
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("droplet spec failed: %v, output: \n%s", err, string(out))
	}

	// == assert ==
	// file exists
	assert.FileExists(t, configJsonPath)

	// contents validation
	// read config.json
	data, readFileErr := os.ReadFile(configJsonPath)
	if readFileErr != nil {
		t.Fatalf("read config.json failed")
	}
	var spec spec.Spec
	if err := json.Unmarshal(data, &spec); err != nil {
		t.Fatalf("json.Unmarshal failed")
	}
	// command: /bin/sh -c "echo Hello World"
	assert.Equal(t, []string{"/bin/sh", "-c", "echo Hello World"}, spec.Process.Args)
}

func TestDropletRun_Create(t *testing.T) {
	// skip in environments where sudo cannot be used without a password.
	// if password input is possible in the local environment, run the test
	if !isLocalEnvironment(t) {
		if !isSudoNotNeedPassword(t) {
			t.Skip("sudo -n not available, skipping integration test")
		}
	}

	// == arrange ==
	path := t.TempDir()
	containerId := "12345"
	rootDir := filepath.Join(path, containerId)
	// mkdir
	if err := os.Mkdir(rootDir, 0o755); err != nil {
		t.Fatalf("mkdir failed, path: %s", rootDir)
	}

	// create config.json
	// config.json parameter
	command := `/bin/sh -c "echo Hello World"`
	hostname := "mycontainer"
	cmdSpec := exec.Command("sudo", binary,
		"spec",
		"--command", command,
		"--hostname", hostname,
		"--output", rootDir,
	)
	if err := cmdSpec.Run(); err != nil {
		t.Fatalf("droplet spec: create config.json failed, path: %s\n%v", rootDir, err)
	}

	// create command
	cmdCreate := exec.Command("sudo", "env", "RAIND_ROOT_DIR="+path,
		binary, "create",
		"--print-pid",
		containerId,
	)

	// == act ==
	// 1. create
	createOut, createErr := cmdCreate.CombinedOutput()
	if createErr != nil {
		t.Fatalf("droplet create failed: %v, output: \n%s", createErr, string(createOut))
	}
	// 2. extract pid and set cleanup
	pid := extractPidFromOutput(t, createOut)
	defer cleanupInitProcess(t, pid)

	// == assert ==
	// fifo.exec exists
	assert.FileExists(t, filepath.Join(rootDir, "exec.fifo"))
	// init process exists
	assert.True(t, isProcessExists(t, pid))
}

func TestDropletRun_CreateAndStart(t *testing.T) {
	// skip in environments where sudo cannot be used without a password.
	// if password input is possible in the local environment, run the test
	if !isLocalEnvironment(t) {
		if !isSudoNotNeedPassword(t) {
			t.Skip("sudo -n not available, skipping integration test")
		}
	}

	// == arrange ==
	path := t.TempDir()
	containerId := "12345"
	rootDir := filepath.Join(path, containerId)
	// mkdir
	if err := os.Mkdir(rootDir, 0o755); err != nil {
		t.Fatalf("mkdir failed, path: %s", rootDir)
	}

	// create config.json
	// config.json parameter
	command := `/bin/sh -c "echo Hello World"`
	hostname := "mycontainer"
	cmdSpec := exec.Command("sudo", binary,
		"spec",
		"--command", command,
		"--hostname", hostname,
		"--output", rootDir,
	)
	if err := cmdSpec.Run(); err != nil {
		t.Fatalf("droplet spec: create config.json failed, path: %s\n%v", rootDir, err)
	}

	// create command
	cmdCreate := exec.Command("sudo", "env", "RAIND_ROOT_DIR="+path,
		binary, "create",
		"--print-pid",
		containerId,
	)
	// start command
	cmdStart := exec.Command("sudo", "env", "RAIND_ROOT_DIR="+path,
		binary, "start", containerId,
	)

	// == act ==
	// 1. create
	createOut, createErr := cmdCreate.CombinedOutput()
	if createErr != nil {
		t.Fatalf("droplet create failed: %v, output: \n%s", createErr, string(createOut))
	}
	// 2. extract pid and set cleanup
	pid := extractPidFromOutput(t, createOut)
	defer cleanupInitProcess(t, pid)

	// 3. start
	startOut, startErr := cmdStart.CombinedOutput()
	if startErr != nil {
		t.Fatalf("droplet start failed: %v, output: \n%s", startErr, string(startOut))
	}

	// == assert ==
	// fifo.exec not exists
	assert.NoFileExists(t, filepath.Join(rootDir, "exec.fifo"))
	// init process not exists
	assert.False(t, isProcessExists(t, pid))
}

func TestDropletRun_Run(t *testing.T) {
	// skip in environments where sudo cannot be used without a password.
	// if password input is possible in the local environment, run the test
	if !isLocalEnvironment(t) {
		if !isSudoNotNeedPassword(t) {
			t.Skip("sudo -n not available, skipping integration test")
		}
	}

	// == arrange ==
	path := t.TempDir()
	containerId := "12345"
	rootDir := filepath.Join(path, containerId)
	// mkdir
	if err := os.Mkdir(rootDir, 0o755); err != nil {
		t.Fatalf("mkdir failed, path: %s", rootDir)
	}

	// create config.json
	// config.json parameter
	command := `/bin/sh -c "echo Hello World"`
	hostname := "mycontainer"
	cmdSpec := exec.Command("sudo", binary,
		"spec",
		"--command", command,
		"--hostname", hostname,
		"--output", rootDir,
	)
	if err := cmdSpec.Run(); err != nil {
		t.Fatalf("droplet spec: create config.json failed, path: %s\n%v", rootDir, err)
	}

	// run command
	cmdRun := exec.Command("sudo", "env", "RAIND_ROOT_DIR="+path,
		binary, "run",
		"--print-pid",
		containerId,
	)

	// == act ==
	// run
	runOut, runErr := cmdRun.CombinedOutput()
	if runErr != nil {
		t.Fatalf("droplet create failed: %v, output: \n%s", runErr, string(runOut))
	}
	// extract pid and set cleanup
	pid := extractPidFromOutput(t, runOut)
	defer cleanupInitProcess(t, pid)

	// == assert ==
	// fifo.exec not exists
	assert.NoFileExists(t, filepath.Join(rootDir, "exec.fifo"))
	// init process not exists
	assert.False(t, isProcessExists(t, pid))
}
