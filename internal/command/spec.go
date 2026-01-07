package command

import (
	"fmt"
	"strings"

	"droplet/internal/spec"

	"github.com/google/shlex"
	"github.com/urfave/cli/v2"
)

func commandSpec() *cli.Command {
	return &cli.Command{
		Name:  "spec",
		Usage: "create a new specification file",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "rootfs",
				Usage: "path to container root filesystem",
				Value: "rootfs",
			},
			&cli.StringSliceFlag{
				Name:  "mount",
				Usage: "mount info (source:dest:options)",
			},
			&cli.StringFlag{
				Name:  "cwd",
				Usage: "container working directory",
				Value: "/",
			},
			&cli.StringSliceFlag{
				Name:  "env",
				Usage: "environment variables (KEY=VALUE)",
			},
			&cli.StringFlag{
				Name:  "command",
				Usage: "container entrypoint",
				Value: "sh",
			},
			&cli.StringSliceFlag{
				Name:  "ns",
				Usage: "namespace target [mount|network|uts|pid|ipc|user|cgroup]",
			},
			&cli.StringFlag{
				Name:  "hostname",
				Usage: "container hostname",
			},

			// network
			&cli.StringFlag{
				Name:  "host_if_name",
				Usage: "host interface name",
				Value: "eth0",
			},
			&cli.StringFlag{
				Name:  "bridge_if_name",
				Usage: "bridge interface name",
				Value: "raind_br0",
			},
			&cli.StringFlag{
				Name:  "if_name",
				Usage: "container interface name",
				Value: "eth0",
			},
			&cli.StringFlag{
				Name:  "if_addr",
				Usage: "container interface address",
				Value: "172.16.0.1/24",
			},
			&cli.StringFlag{
				Name:  "if_gateway",
				Usage: "container interface gateway",
				Value: "172.16.0.254",
			},
			&cli.StringSliceFlag{
				Name:  "dns",
				Usage: "dns server",
			},

			// layer
			&cli.StringSliceFlag{
				Name:  "image_layer",
				Usage: "image layer path",
			},
			&cli.StringFlag{
				Name:  "upper_dir",
				Usage: "upper directory",
			},
			&cli.StringFlag{
				Name:  "work_dir",
				Usage: "work directory",
			},

			// hook
			&cli.StringSliceFlag{
				Name:  "hook-prestart",
				Usage: "prestart hook (format: path[,arg1,arg2,...])",
			},
			&cli.StringSliceFlag{
				Name:  "hook-create-runtime",
				Usage: "createRuntime hook (format: path[,arg1,arg2,...])",
			},
			&cli.StringSliceFlag{
				Name:  "hook-create-container",
				Usage: "createContainer hook (format: path[,arg1,arg2,...])",
			},
			&cli.StringSliceFlag{
				Name:  "hook-start-container",
				Usage: "startContainer hook (format: path[,arg1,arg2,...])",
			},
			&cli.StringSliceFlag{
				Name:  "hook-poststart",
				Usage: "poststart hook (format: path[,arg1,arg2,...])",
			},
			&cli.StringSliceFlag{
				Name:  "hook-poststop",
				Usage: "poststop hook (format: path[,arg1,arg2,...])",
			},

			&cli.StringFlag{
				Name:  "output",
				Usage: "output path",
				Value: ".",
			},
		},
		Action: runCreateConfigFile,
	}
}

func runCreateConfigFile(ctx *cli.Context) error {
	// create config options
	configOptions, err := createConfigOptions(ctx)
	if err != nil {
		return err
	}

	// build configuration file(config.json)
	if err := spec.CreateConfigFile(ctx.String("output")+"/config.json", configOptions); err != nil {
		return err
	}

	return nil
}

func createConfigOptions(ctx *cli.Context) (spec.ConfigOptions, error) {
	// parse flags and create ConfigOptions
	// rootfs
	rootfs := ctx.String("rootfs")

	// mount
	mounts, err := parseMountFlag(ctx.StringSlice("mount"))
	if err != nil {
		return spec.ConfigOptions{}, err
	}

	// process
	// cwd
	cwd := ctx.String("cwd")
	// env
	env := ctx.StringSlice("env")
	// args
	args, err := parseCommandFlag(ctx.String("command"))
	if err != nil {
		return spec.ConfigOptions{}, err
	}

	// namespace
	namespace := ctx.StringSlice("ns")

	// hostname
	hostname := ctx.String("hostname")

	// net
	// host interface name
	hostIfName := ctx.String("host_if_name")
	// bridge interface name
	brIfName := ctx.String("bridge_if_name")
	// interface name
	ifName := ctx.String("if_name")
	// interface address
	ifAddr := ctx.String("if_addr")
	// gateway
	ifGateway := ctx.String("if_gateway")
	// dns
	dns := ctx.StringSlice("dns")

	// image
	// image layer
	imageLayer := ctx.StringSlice("image_layer")
	// upper dir
	upperDir := ctx.String("upper_dir")
	// work dir
	workDir := ctx.String("work_dir")

	// hook
	// prestart
	prestartHook, err := parseHookFlag(ctx.StringSlice("hook-prestart"))
	if err != nil {
		return spec.ConfigOptions{}, err
	}
	createRuntimeHook, err := parseHookFlag(ctx.StringSlice("hook-create-runtime"))
	if err != nil {
		return spec.ConfigOptions{}, err
	}
	createContainerHook, err := parseHookFlag(ctx.StringSlice("hook-create-container"))
	if err != nil {
		return spec.ConfigOptions{}, err
	}
	startContainerHook, err := parseHookFlag(ctx.StringSlice("hook-start-container"))
	if err != nil {
		return spec.ConfigOptions{}, err
	}
	poststartHook, err := parseHookFlag(ctx.StringSlice("hook-poststart"))
	if err != nil {
		return spec.ConfigOptions{}, err
	}
	poststopHook, err := parseHookFlag(ctx.StringSlice("hook-poststop"))
	if err != nil {
		return spec.ConfigOptions{}, err
	}

	return spec.ConfigOptions{
		Rootfs: rootfs,
		Mounts: mounts,
		Process: spec.ProcessOption{
			Cwd:  cwd,
			Env:  env,
			Args: args,
		},
		Namespace: namespace,
		Hostname:  hostname,
		Net: spec.NetOption{
			HostInterface:       hostIfName,
			BridgeInterfaceName: brIfName,
			InterfaceName:       ifName,
			Address:             ifAddr,
			Gateway:             ifGateway,
			Dns:                 dns,
		},
		Image: spec.ImageOption{
			ImageLayer: imageLayer,
			UpperDir:   upperDir,
			WorkDir:    workDir,
		},
		Hooks: spec.HookLifecycleOption{
			Prestart:        prestartHook,
			CreateRuntime:   createRuntimeHook,
			CreateContainer: createContainerHook,
			StartContainer:  startContainerHook,
			Poststart:       poststartHook,
			Poststop:        poststopHook,
		},
	}, nil
}

func parseMountFlag(mounts []string) ([]spec.MountOption, error) {
	var mountOption []spec.MountOption
	for _, mount := range mounts {
		parts := strings.SplitN(mount, ":", 3)
		if len(parts) < 2 {
			return []spec.MountOption{}, fmt.Errorf("invalid mount format")
		}

		// source, deestination
		src := parts[0]
		dst := parts[1]

		// options
		var opts []string
		if len(parts) == 3 && parts[2] != "" {
			opts = strings.Split(parts[2], ",")
		} else {
			opts = append(opts, "bind")
		}

		mountOption = append(mountOption, spec.MountOption{
			Destination: dst,
			Type:        "",
			Source:      src,
			Options:     opts,
		})
	}
	return mountOption, nil
}

func parseCommandFlag(s string) ([]string, error) {
	args, err := shlex.Split(s)
	if err != nil {
		return []string{}, err
	}
	return args, nil
}

func parseHookFlag(value []string) ([]spec.HookOption, error) {
	var hooks []spec.HookOption

	for _, v := range value {
		if v == "" {
			continue
		}
		parts := strings.Split(v, ",")
		if len(parts) == 0 || parts[0] == "" {
			return []spec.HookOption{}, fmt.Errorf("invalid hook: %q (path is required)", v)
		}
		h := spec.HookOption{
			Path: parts[0],
		}
		if len(parts) > 1 {
			h.Args = parts[1:]
		}
		hooks = append(hooks, h)
	}
	return hooks, nil
}
