//go:build integration
// +build integration

package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

const runtimeMainPkg = "../cmd/droplet"

var runtimeBin string

func TestMain(m *testing.M) {
	// build binary to tmp dir
	tmpDir, err := os.MkdirTemp("", "droplet-bin-*")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmpDir)

	runtimeBin = filepath.Join(tmpDir, "droplet")

	build := exec.Command("go", "build", "-o", runtimeBin, runtimeMainPkg)
	build.Stdout = os.Stdout
	build.Stderr = os.Stderr
	if err := build.Run(); err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func runRuntime(t *testing.T, root string, args ...string) error {
	t.Helper()

	cmd := exec.Command(runtimeBin, args...)
	// set root dir to TmpDir
	cmd.Env = append(os.Environ(), "RAIND_ROOT_DIR="+root)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func TestIntegration_CreateAndStart(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	containerId := "ct1"

	bundleDir := filepath.Join(root, containerId)
	if err := os.MkdirAll(bundleDir, 0o755); err != nil {
		t.Fatalf("mkdir bundle: %v", err)
	}

	// create config.json
	// command: droplet spec <options>...
	if err := runRuntime(t, root,
		"spec",
		"--rootfs", "/",
		"--cwd", "/",
		"--command", `/bin/sh -c "echo Hello World from Container!"`,
		"--hostname", containerId,
		"--if_name", "eth0",
		"--if_addr", "10.166.0.1/24", "--if_gateway", "10.166.0.254", "--dns", "8.8.8.8",
		"--image_layer", "/image/path",
		"--upper_dir", "/upper/path", "--work_dir", "/work/path", "--merge_dir", "/merge/path",
		"--output", bundleDir,
	); err != nil {
		t.Fatalf("spec failed: %v", err)
	}

	// create
	// command: droplet create ct1
	if err := runRuntime(t, root, "create", containerId); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	// start
	// command: droplet start ct1
	if err := runRuntime(t, root, "start", containerId); err != nil {
		t.Fatalf("start failed: %v", err)
	}
}

func TestIntegration_Run(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	containerId := "ct1"

	bundleDir := filepath.Join(root, containerId)
	if err := os.MkdirAll(bundleDir, 0o755); err != nil {
		t.Fatalf("mkdir bundle: %v", err)
	}

	// create config.json
	// command: droplet spec <options>...
	if err := runRuntime(t, root,
		"spec",
		"--rootfs", "/",
		"--cwd", "/",
		"--command", `/bin/sh -c "echo Hello World from Container!"`,
		"--hostname", containerId,
		"--if_name", "eth0",
		"--if_addr", "10.166.0.1/24", "--if_gateway", "10.166.0.254", "--dns", "8.8.8.8",
		"--image_layer", "/image/path",
		"--upper_dir", "/upper/path", "--work_dir", "/work/path", "--merge_dir", "/merge/path",
		"--output", bundleDir,
	); err != nil {
		t.Fatalf("spec failed: %v", err)
	}

	// run
	// command: droplet run ct1
	if err := runRuntime(t, root, "run", containerId); err != nil {
		t.Fatalf("run failed: %v", err)
	}
}

func TestIntegration_Run_Namespace(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	containerId := "ct1"

	bundleDir := filepath.Join(root, containerId)
	if err := os.MkdirAll(bundleDir, 0o755); err != nil {
		t.Fatalf("mkdir bundle: %v", err)
	}

	// create config.json
	// command: droplet spec <options>...
	if err := runRuntime(t, root,
		"spec",
		"--rootfs", "/",
		"--cwd", "/",
		"--command", `/bin/sh -c "echo Hello World from Container!"`,
		"--ns", "mount", "--ns", "network", "--ns", "uts", "--ns", "pid", "--ns", "ipc", "--ns", "user", "--ns", "cgroup",
		"--hostname", containerId,
		"--if_name", "eth0",
		"--if_addr", "10.166.0.1/24", "--if_gateway", "10.166.0.254", "--dns", "8.8.8.8",
		"--image_layer", "/image/path",
		"--upper_dir", "/upper/path", "--work_dir", "/work/path", "--merge_dir", "/merge/path",
		"--output", bundleDir,
	); err != nil {
		t.Fatalf("spec failed: %v", err)
	}

	// run
	// command: droplet run ct1
	if err := runRuntime(t, root, "run", containerId); err != nil {
		t.Fatalf("run failed: %v", err)
	}
}
