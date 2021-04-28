package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/northvolt/go-kd6rmx"
	"github.com/peterbourgon/ff/v3/ffcli"
)

// kd6ctl is a command line utility to change config on the KD6RMX contact image sensor.

func main() {
	var (
		rootFlagSet = flag.NewFlagSet("kd6ctl", flag.ExitOnError)
		port        = rootFlagSet.String("p", "/dev/corser/XtiumCLMX41_s0", "port of KD6RMX sensor to use")
	)

	load := &ffcli.Command{
		Name:       "load",
		ShortUsage: "kd6ctl load <arg>",
		ShortHelp:  "Load user settings.",
		Exec: func(_ context.Context, args []string) error {
			if n := len(args); n != 1 {
				return fmt.Errorf("load requires exactly 1 argument, but you provided %d", n)
			}
			preset, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			cis := kd6rmx.Sensor{Port: *port}
			return cis.LoadSettings(preset)
		},
	}

	ledson := &ffcli.Command{
		Name:       "on",
		ShortUsage: "kd6ctl led on",
		ShortHelp:  "Turn on LED.",
		Exec: func(_ context.Context, args []string) error {
			return flag.ErrHelp
		},
	}

	ledsoff := &ffcli.Command{
		Name:       "off",
		ShortUsage: "kd6ctl led off",
		ShortHelp:  "Turn on LED.",
		Exec: func(_ context.Context, args []string) error {
			return flag.ErrHelp
		},
	}

	leds := &ffcli.Command{
		Name:        "led",
		ShortUsage:  "kd6ctl led <arg>",
		ShortHelp:   "Sets LED on sensor.",
		Subcommands: []*ffcli.Command{ledson, ledsoff},
		Exec: func(_ context.Context, args []string) error {
			return flag.ErrHelp
		},
	}

	root := &ffcli.Command{
		ShortUsage:  "kd6ctl [flags] <subcommand>",
		ShortHelp:   "kd6ctl is a command line utility to change config on the KD6RMX contact image sensor.",
		FlagSet:     rootFlagSet,
		Subcommands: []*ffcli.Command{load, leds},
		Exec: func(context.Context, []string) error {
			return flag.ErrHelp
		},
	}

	if err := root.ParseAndRun(context.Background(), os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
