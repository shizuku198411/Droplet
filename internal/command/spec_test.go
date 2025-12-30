package command

import (
	"flag"
	"reflect"
	"slices"
	"testing"

	"droplet/internal/spec"

	"github.com/urfave/cli/v2"
)

// == test for func:parseMountFlag ==
func TestParseMountFlag_1(t *testing.T) {
	input := []string{
		"/src:/dst:opt1,opt2",
	}

	result, _ := parseMountFlag(input)

	expect := []spec.MountOption{
		{
			Destination: "/dst",
			Type:        "",
			Source:      "/src",
			Options: []string{
				"opt1",
				"opt2",
			},
		},
	}

	if !reflect.DeepEqual(result, expect) {
		t.Errorf("unexpected result.\nwant=%#v\ngot =%#v", expect, result)
	}
}

func TestParseMountFlag_2(t *testing.T) {
	input := []string{
		"/src:/dst",
	}

	result, _ := parseMountFlag(input)

	expect := []spec.MountOption{
		{
			Destination: "/dst",
			Type:        "",
			Source:      "/src",
			Options: []string{
				"bind",
			},
		},
	}

	if !reflect.DeepEqual(result, expect) {
		t.Errorf("unexpected result.\nwant=%#v\ngot =%#v", expect, result)
	}
}

func TestParseMountFlag_3(t *testing.T) {
	input := []string{
		"/src:/dst",
		"/src2:/dst2",
	}

	result, _ := parseMountFlag(input)

	expect := []spec.MountOption{
		{
			Destination: "/dst",
			Type:        "",
			Source:      "/src",
			Options: []string{
				"bind",
			},
		},
		{
			Destination: "/dst2",
			Type:        "",
			Source:      "/src2",
			Options: []string{
				"bind",
			},
		},
	}

	if !reflect.DeepEqual(result, expect) {
		t.Errorf("unexpected result.\nwant=%#v\ngot =%#v", expect, result)
	}
}

// ==================================

// == test for func:parseCommandFlag ==
func TestParseCommandFlag_1(t *testing.T) {
	input := "/bin/bash"

	result, _ := parseCommandFlag(input)

	expect := []string{
		"/bin/bash",
	}

	if !slices.Equal(result, expect) {
		t.Errorf("unexpected result.\nwant=%#v\ngot =%#v", expect, result)
	}
}

func TestParseCommandFlag_2(t *testing.T) {
	input := "echo \"Hello World\""

	result, _ := parseCommandFlag(input)

	expect := []string{
		"echo",
		"Hello World",
	}

	if !slices.Equal(result, expect) {
		t.Errorf("unexpected result.\nwant=%#v\ngot =%#v", expect, result)
	}
}

func TestParseCommandFlag_3(t *testing.T) {
	input := "apt install -y iputils-ping"

	result, _ := parseCommandFlag(input)

	expect := []string{
		"apt",
		"install",
		"-y",
		"iputils-ping",
	}

	if !slices.Equal(result, expect) {
		t.Errorf("unexpected result.\nwant=%#v\ngot =%#v", expect, result)
	}
}

// ====================================

// == test for func:createConfigOptions ==
func TestCreateConfigOptions_1(t *testing.T) {
	// dummy app
	app := &cli.App{}
	// set flag
	set := flag.NewFlagSet("testFlagSet", 0)
	set.String("rootfs", "", "")
	set.String("cwd", "", "")
	set.String("command", "", "")
	set.String("hostname", "", "")
	set.String("if_name", "", "")
	set.String("if_addr", "", "")
	set.String("if_gateway", "", "")
	set.String("image_layer", "", "")
	set.String("upper_dir", "", "")
	set.String("work_dir", "", "")
	set.String("merge_dir", "", "")

	mounts := cli.NewStringSlice()
	envs := cli.NewStringSlice()
	dns := cli.NewStringSlice()

	set.Var(mounts, "mount", "")
	set.Var(envs, "env", "")
	set.Var(dns, "dns", "")

	// input
	_ = set.Set("rootfs", "/")

	_ = set.Set("mount", "/src:/dst")
	_ = set.Set("mount", "/src2:/dst2")

	_ = set.Set("cwd", "/")
	_ = set.Set("env", "PATH=/usr/bin:/bin")
	_ = set.Set("env", "KEY=VALUE")
	_ = set.Set("command", "/bin/bash")

	_ = set.Set("hostname", "mycontainer")

	_ = set.Set("if_name", "eth0")
	_ = set.Set("if_addr", "10.166.0.1/24")
	_ = set.Set("if_gateway", "10.166.0.254")
	_ = set.Set("dns", "8.8.8.8")
	_ = set.Set("dns", "8.8.4.4")

	_ = set.Set("image_layer", "/image/path")
	_ = set.Set("upper_dir", "/upper/path")
	_ = set.Set("work_dir", "/work/path")
	_ = set.Set("merge_dir", "/merge/path")

	app.DisableSliceFlagSeparator = true

	// create context
	ctx := cli.NewContext(app, set, nil)

	// func
	result, _ := createConfigOptions(ctx)

	expect := spec.ConfigOptions{
		Rootfs: "/",
		Mounts: []spec.MountOption{
			{
				Destination: "/dst",
				Type:        "",
				Source:      "/src",
				Options: []string{
					"bind",
				},
			},
			{
				Destination: "/dst2",
				Type:        "",
				Source:      "/src2",
				Options: []string{
					"bind",
				},
			},
		},
		Process: spec.ProcessOption{
			Cwd: "/",
			Env: []string{
				"PATH=/usr/bin:/bin",
				"KEY=VALUE",
			},
			Args: []string{
				"/bin/bash",
			},
		},
		Hostname: "mycontainer",
		Net: spec.NetOption{
			InterfaceName: "eth0",
			Address:       "10.166.0.1/24",
			Gateway:       "10.166.0.254",
			Dns: []string{
				"8.8.8.8",
				"8.8.4.4",
			},
		},
		Image: spec.ImageOption{
			ImageLayer: "/image/path",
			UpperDir:   "/upper/path",
			WorkDir:    "/work/path",
			MergeDir:   "/merge/path",
		},
	}

	if !reflect.DeepEqual(result, expect) {
		t.Errorf("unexpected result.\nwant=%#v\ngot =%#v", expect, result)
	}
}

// =======================================
