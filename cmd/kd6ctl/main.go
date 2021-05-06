// kd6ctl is a command line utility for configuring the KD6RMX contact image sensor.

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

func main() {
	var (
		rootFlagSet = flag.NewFlagSet("kd6ctl", flag.ExitOnError)
		port        = rootFlagSet.String("p", "/dev/corser/XtiumCLMX41_s0", "port of KD6RMX sensor to use")
	)

	load := &ffcli.Command{
		Name:       "load",
		ShortUsage: "kd6ctl load <preset>",
		ShortHelp:  "Load user settings.",
		Exec: func(_ context.Context, args []string) error {
			if n := len(args); n != 1 {
				return fmt.Errorf("load requires the number of the preset you want ot load")
			}

			preset, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			cis := kd6rmx.Sensor{Port: *port}
			return cis.LoadSettings(preset)
		},
	}

	save := &ffcli.Command{
		Name:       "save",
		ShortUsage: "kd6ctl save <preset>",
		ShortHelp:  "Save current settings into a user preset.",
		Exec: func(_ context.Context, args []string) error {
			if n := len(args); n != 1 {
				return fmt.Errorf("load requires the number of the preset you want ot load")
			}

			preset, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			cis := kd6rmx.Sensor{Port: *port}
			return cis.SaveSettings(preset)
		},
	}

	outputfreq := &ffcli.Command{
		Name:       "frequency",
		ShortUsage: "kd6ctl frequency <freq>",
		ShortHelp:  "Change output frequency (in Mhz).",
		Exec: func(_ context.Context, args []string) error {
			if n := len(args); n != 1 {
				return fmt.Errorf("outputfreq requires providing the desired frequency")
			}

			freq, err := strconv.ParseFloat(args[0], 32)
			if err != nil {
				return err
			}

			cis := kd6rmx.Sensor{Port: *port}
			return cis.OutputFrequency(float32(freq))
		},
	}

	outputfmt := &ffcli.Command{
		Name:       "format",
		ShortUsage: "kd6ctl format <bits> <interface> <config> <num>",
		ShortHelp:  "Change output format.",
		Exec: func(_ context.Context, args []string) error {
			if n := len(args); n != 4 {
				return fmt.Errorf("outputfmt requires providing all the needed params")
			}

			var bits kd6rmx.PixelOutputBits
			switch args[0] {
			case "10":
				bits = kd6rmx.PixelOutputBits10
			case "8":
				bits = kd6rmx.PixelOutputBits8
			default:
				return fmt.Errorf("invalid number of bits")
			}

			var intf kd6rmx.PixelOutputInterface
			switch args[1] {
			case "serial":
				intf = kd6rmx.PixelOutputSerial
			case "parallel":
				intf = kd6rmx.PixelOutputParallel
			default:
				return fmt.Errorf("invalid interface, must be serial or parallel")
			}

			var conf kd6rmx.PixelOutputConfig
			switch args[2] {
			case "base":
				conf = kd6rmx.PixelOutputBase
			case "medium":
				conf = kd6rmx.PixelOutputMedium
			default:
				return fmt.Errorf("invalid config, must be base or medium")
			}

			num, err := strconv.Atoi(args[3])
			if err != nil {
				return err
			}

			cis := kd6rmx.Sensor{Port: *port}
			return cis.PixelOutputFormat(bits, intf, conf, num)
		},
	}

	interp := &ffcli.Command{
		Name:       "interpolation",
		ShortUsage: "kd6ctl interpolation <on/off>",
		ShortHelp:  "Set interpolation on/off.",
		Exec: func(_ context.Context, args []string) error {
			if n := len(args); n != 1 {
				return fmt.Errorf("interp nust be either 'on' or 'off'")
			}

			var on bool
			switch args[0] {
			case "on":
				on = true
			case "off":
				on = false
			default:
				return fmt.Errorf("invalid interpolation, must be on or off")
			}

			cis := kd6rmx.Sensor{Port: *port}
			return cis.PixelInterpolation(on)
		},
	}

	dark := &ffcli.Command{
		Name:       "dark",
		ShortUsage: "kd6ctl dark <on/off>",
		ShortHelp:  "Dark correction on/off.",
		Exec: func(_ context.Context, args []string) error {
			if n := len(args); n != 1 {
				return fmt.Errorf("dark correction must be either 'on' or 'off'")
			}

			var on bool
			switch args[0] {
			case "on":
				on = true
			case "off":
				on = false
			default:
				return fmt.Errorf("invalid dark correction, must be on or off")
			}

			cis := kd6rmx.Sensor{Port: *port}
			return cis.DarkCorrectionEnabled(on)
		},
	}

	white := &ffcli.Command{
		Name:       "white",
		ShortUsage: "kd6ctl white <on/off>",
		ShortHelp:  "White correction on/off.",
		Exec: func(_ context.Context, args []string) error {
			if n := len(args); n != 1 {
				return fmt.Errorf("white correction must be either 'on' or 'off'")
			}

			var on bool
			switch args[0] {
			case "on":
				on = true
			case "off":
				on = false
			default:
				return fmt.Errorf("invalid white correction, must be on or off")
			}

			cis := kd6rmx.Sensor{Port: *port}
			return cis.WhiteCorrectionEnabled(on)
		},
	}

	leds := &ffcli.Command{
		Name:       "led",
		ShortUsage: "kd6ctl led <A/B/AB> <on/off> [pulse]",
		ShortHelp:  "Sets LEDs on sensor on or off.",
		Exec: func(_ context.Context, args []string) error {
			if len(args) < 2 {
				return fmt.Errorf("invalid led command params")
			}

			leds := args[0]
			if len(leds) == 0 {
				return fmt.Errorf("invalid leds")
			}

			var on bool
			switch args[1] {
			case "on":
				on = true
			case "off":
				on = false
			default:
				return fmt.Errorf("invalid led setting, must be on or off")
			}

			var pulse int
			if len(args) < 3 {
				pulse = 1
			} else {
				var err error
				pulse, err = strconv.Atoi(args[2])
				if err != nil {
					return err
				}
			}

			cis := kd6rmx.Sensor{Port: *port}
			return cis.LEDControl(leds, on, pulse)
		},
	}

	root := &ffcli.Command{
		ShortUsage:  "kd6ctl [flags] <subcommand>",
		ShortHelp:   "kd6ctl is a command line utility to change config on the KD6RMX contact image sensor.",
		FlagSet:     rootFlagSet,
		Subcommands: []*ffcli.Command{load, save, outputfreq, outputfmt, interp, dark, white, leds},
		Exec: func(context.Context, []string) error {
			return flag.ErrHelp
		},
	}

	if err := root.ParseAndRun(context.Background(), os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
