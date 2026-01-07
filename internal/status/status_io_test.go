package status

import (
	"droplet/internal/spec"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateStatusFile_Success(t *testing.T) {
	// == arrange ==
	path := t.TempDir()
	t.Setenv("RAIND_ROOT_DIR", path)
	containerId := "12345"
	pid := 11111
	containerStatus := CREATED
	rootfs := "/path/to/rootfs"
	bundle := "/path/to/bundle"
	annotation := spec.AnnotationObject{
		Version: "0.1.0",
		Image:   "imageannotation",
		Net:     "netannotation",
	}
	containerStatusHandler := &StatusHandler{}
	if err := os.MkdirAll(filepath.Join(path, containerId), 0o755); err != nil {
		t.Fatalf("create directory failed")
	}

	// == act ==
	err := containerStatusHandler.CreateStatusFile(containerId, pid, containerStatus, rootfs, bundle, annotation)

	// == assert ==
	// file created
	assert.FileExists(t, filepath.Join(path, containerId, "state.json"))

	// file content veryfi
	var content StatusObject
	data, readErr := os.ReadFile(filepath.Join(path, containerId, "state.json"))
	if readErr != nil {
		t.Fatalf("read file failed")
	}
	json.Unmarshal(data, &content)

	assert.Equal(t, "12345", content.Id)
	assert.Equal(t, 11111, content.Pid)
	assert.Equal(t, CREATED.String(), content.Status)
	assert.Equal(t, "/path/to/bundle", content.Bundle)

	// error is nil
	assert.Nil(t, err)
}

func TestUpdateStatus_StatusUpdateSuccess(t *testing.T) {
	// == arrange ==
	path := t.TempDir()
	t.Setenv("RAIND_ROOT_DIR", path)
	containerId := "12345"
	pid := 11111
	containerStatus := CREATED
	bundle := "/path/to/bundle"
	containerStatusObject := StatusObject{
		Id:     containerId,
		Pid:    pid,
		Status: containerStatus.String(),
		Bundle: bundle,
	}
	// create state.json to tmp dir
	if err := os.MkdirAll(filepath.Join(path, containerId), 0o755); err != nil {
		t.Fatalf("create directory failed")
	}
	f, createErr := os.Create(filepath.Join(path, containerId, "state.json"))
	if createErr != nil {
		t.Fatalf("create file failed")
	}
	defer f.Close()
	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "    ")
	encoder.Encode(containerStatusObject)

	updateStatus := RUNNING
	containerStatusHandler := &StatusHandler{}

	// == act ==
	err := containerStatusHandler.UpdateStatus(containerId, updateStatus, -1)

	// == assert ==
	// file content veryfi
	var content StatusObject
	data, readErr := os.ReadFile(filepath.Join(path, containerId, "state.json"))
	if readErr != nil {
		t.Fatalf("read file failed")
	}
	json.Unmarshal(data, &content)

	// status changed from CREATED to RUNNING
	assert.Equal(t, RUNNING.String(), content.Status)
	// pid not changed
	assert.Equal(t, 11111, content.Pid)

	// error is nil
	assert.Nil(t, err)
}

func TestUpdateStatus_PidUpdateSuccess(t *testing.T) {
	// == arrange ==
	path := t.TempDir()
	t.Setenv("RAIND_ROOT_DIR", path)
	containerId := "12345"
	pid := 11111
	containerStatus := CREATED
	bundle := "/path/to/bundle"
	containerStatusObject := StatusObject{
		Id:     containerId,
		Pid:    pid,
		Status: containerStatus.String(),
		Bundle: bundle,
	}
	// create state.json to tmp dir
	if err := os.MkdirAll(filepath.Join(path, containerId), 0o755); err != nil {
		t.Fatalf("create directory failed")
	}
	f, createErr := os.Create(filepath.Join(path, containerId, "state.json"))
	if createErr != nil {
		t.Fatalf("create file failed")
	}
	defer f.Close()
	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "    ")
	encoder.Encode(containerStatusObject)

	updatePid := 22222
	containerStatusHandler := &StatusHandler{}

	// == act ==
	err := containerStatusHandler.UpdateStatus(containerId, -1, updatePid)

	// == assert ==
	// file content veryfi
	var content StatusObject
	data, readErr := os.ReadFile(filepath.Join(path, containerId, "state.json"))
	if readErr != nil {
		t.Fatalf("read file failed")
	}
	json.Unmarshal(data, &content)

	// status not changed
	assert.Equal(t, CREATED.String(), content.Status)
	// pid changed from 11111 to 22222
	assert.Equal(t, 22222, content.Pid)

	// error is nil
	assert.Nil(t, err)
}

func TestUpdateStatus_StatusAndPidUpdateSuccess(t *testing.T) {
	// == arrange ==
	path := t.TempDir()
	t.Setenv("RAIND_ROOT_DIR", path)
	containerId := "12345"
	pid := 11111
	containerStatus := CREATED
	bundle := "/path/to/bundle"
	containerStatusObject := StatusObject{
		Id:     containerId,
		Pid:    pid,
		Status: containerStatus.String(),
		Bundle: bundle,
	}
	// create state.json to tmp dir
	if err := os.MkdirAll(filepath.Join(path, containerId), 0o755); err != nil {
		t.Fatalf("create directory failed")
	}
	f, createErr := os.Create(filepath.Join(path, containerId, "state.json"))
	if createErr != nil {
		t.Fatalf("create file failed")
	}
	defer f.Close()
	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "    ")
	encoder.Encode(containerStatusObject)

	updateStatus := RUNNING
	updatePid := 22222
	containerStatusHandler := &StatusHandler{}

	// == act ==
	err := containerStatusHandler.UpdateStatus(containerId, updateStatus, updatePid)

	// == assert ==
	// file content verify
	var content StatusObject
	data, readErr := os.ReadFile(filepath.Join(path, containerId, "state.json"))
	if readErr != nil {
		t.Fatalf("read file failed")
	}
	json.Unmarshal(data, &content)

	// status changed from created to running
	assert.Equal(t, RUNNING.String(), content.Status)
	// pid changed from 11111 to 22222
	assert.Equal(t, 22222, content.Pid)

	// error is nil
	assert.Nil(t, err)
}

func TestGetPidFromId_Success(t *testing.T) {
	// == arrange ==
	path := t.TempDir()
	t.Setenv("RAIND_ROOT_DIR", path)
	containerId := "12345"
	pid := 11111
	containerStatus := CREATED
	bundle := "/path/to/bundle"
	containerStatusObject := StatusObject{
		Id:     containerId,
		Pid:    pid,
		Status: containerStatus.String(),
		Bundle: bundle,
	}
	// create state.json to tmp dir
	if err := os.MkdirAll(filepath.Join(path, containerId), 0o755); err != nil {
		t.Fatalf("create directory failed")
	}
	f, createErr := os.Create(filepath.Join(path, containerId, "state.json"))
	if createErr != nil {
		t.Fatalf("create file failed")
	}
	defer f.Close()
	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "    ")
	encoder.Encode(containerStatusObject)

	containerStatusHandler := &StatusHandler{}

	// == act ==
	got, err := containerStatusHandler.GetPidFromId(containerId)

	// == assert ==
	// pid: 11111
	assert.Equal(t, 11111, got)

	// error is nil
	assert.Nil(t, err)
}

func TestGetStatusFromId_Success(t *testing.T) {
	// == arrange ==
	path := t.TempDir()
	t.Setenv("RAIND_ROOT_DIR", path)
	containerId := "12345"
	pid := 11111
	containerStatus := CREATED
	bundle := "/path/to/bundle"
	containerStatusObject := StatusObject{
		Id:     containerId,
		Pid:    pid,
		Status: containerStatus.String(),
		Bundle: bundle,
	}
	// create state.json to tmp dir
	if err := os.MkdirAll(filepath.Join(path, containerId), 0o755); err != nil {
		t.Fatalf("create directory failed")
	}
	f, createErr := os.Create(filepath.Join(path, containerId, "state.json"))
	if createErr != nil {
		t.Fatalf("create file failed")
	}
	defer f.Close()
	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "    ")
	encoder.Encode(containerStatusObject)

	containerStatusHandler := &StatusHandler{}

	// == act ==
	got, err := containerStatusHandler.GetStatusFromId(containerId)

	// == assert ==
	// status: CREATED
	assert.Equal(t, CREATED, got)

	// error is nil
	assert.Nil(t, err)
}
