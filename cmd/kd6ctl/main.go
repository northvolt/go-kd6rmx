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
		logging     = rootFlagSet.Bool("log", false, "turn on debug logging")
	)

	version := &ffcli.Command{
		Name:       "version",
		ShortUsage: "kd6ctl version",
		ShortHelp:  "Show version of kd6ctl API.",
		Exec: func(_ context.Context, args []string) error {
			fmt.Println(kd6rmx.Version)
			return nil
		},
	}

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

			cis := kd6rmx.Sensor{Port: *port, Logging: *logging}
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

			cis := kd6rmx.Sensor{Port: *port, Logging: *logging}
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

			cis := kd6rmx.Sensor{Port: *port, Logging: *logging}
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

			cis := kd6rmx.Sensor{Port: *port, Logging: *logging}
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

			cis := kd6rmx.Sensor{Port: *port, Logging: *logging}
			return cis.PixelInterpolation(on)
		},
	}

	dark := &ffcli.Command{
		Name:       "dark",
		ShortUsage: "kd6ctl dark <on/off/adjust>",
		ShortHelp:  "Dark correction on/off/adjust.",
		Exec: func(_ context.Context, args []string) error {
			if n := len(args); n < 1 {
				return fmt.Errorf("dark correction requires a subcommand: 'on', 'off', or 'adjust'")
			}

			cis := kd6rmx.Sensor{Port: *port, Logging: *logging}

			switch args[0] {
			case "on":
				return cis.DarkCorrectionEnabled(true)
			case "off":
				return cis.DarkCorrectionEnabled(false)
			case "adjust":
				return cis.PerformDarkCorrection()
			default:
				return fmt.Errorf("invalid dark correction subcommand, must be 'on', 'off', or 'adjust'")
			}
		},
	}

	white := &ffcli.Command{
		Name:       "white",
		ShortUsage: "kd6ctl white <on/off/adjust/target>",
		ShortHelp:  "White correction on/off/adjust/target.",
		Exec: func(_ context.Context, args []string) error {
			if n := len(args); n < 1 {
				return fmt.Errorf("white correction requires a subcommand: 'on', 'off', 'adjust', or 'target'")
			}

			cis := kd6rmx.Sensor{Port: *port}

			switch args[0] {
			case "on":
				return cis.WhiteCorrectionEnabled(true)
			case "off":
				return cis.WhiteCorrectionEnabled(false)
			case "adjust":
				return cis.PerformWhiteCorrection()
			case "target":
				var target = 250
				if len(args) < 2 {
					fmt.Printf("no white correction target provided, using factory default of %d\n", target)
				} else {
					var err error
					target, err = strconv.Atoi(args[1])
					if err != nil {
						return fmt.Errorf("invalid value for white correction target")
					}
					return cis.WhiteCorrectionTarget(target)
				}
				return cis.PerformWhiteCorrection()
			default:
				return fmt.Errorf("invalid white correction subcommand, must be 'on', 'off', 'adjust', or 'target'")
			}
		},
	}

	leds := &ffcli.Command{
		Name:       "led",
		ShortUsage: "kd6ctl led <A/B/AB> <on/off> [pulse]",
		ShortHelp:  "Turn sensor LEDs on or off.",
		Exec: func(_ context.Context, args []string) error {
			if len(args) < 2 {
				return fmt.Errorf("led command requires specific LEDs either 'a', 'b', 'ab'. You must also specify to set LEDs 'on' or 'off'")
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

			cis := kd6rmx.Sensor{Port: *port, Logging: *logging}
			return cis.LEDControl(leds, on, pulse)
		},
	}

	duty := &ffcli.Command{
		Name:       "duty",
		ShortUsage: "kd6ctl duty <a/b> <duty>",
		ShortHelp:  "Set LED duty illumination period register value. Valid range 0 to 4095.",
		Exec: func(_ context.Context, args []string) error {
			if len(args) < 2 {
				return fmt.Errorf("duty command requires specific LEDs either 'a' or 'b'. You must also specify the duty to set LEDs to")
			}

			led := args[0]
			if led != "a" && led != "b" {
				return fmt.Errorf("invalid led value. must be 'a' or 'b'")
			}

			duty, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}

			cis := kd6rmx.Sensor{Port: *port}
			return cis.LEDDutyCycle(led, duty)
		},
	}

	gain := &ffcli.Command{
		Name:       "gain",
		ShortUsage: "kd6ctl gain <value/on/off>",
		ShortHelp:  "Enables the gain control and sets the specified value ",
		Exec: func(_ context.Context, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("adjust the gain number")
			}

			cis := kd6rmx.Sensor{Port: *port, Logging: *logging}
			switch args[0] {
			case "on":
				cis.GainAmplifierEnabled(true)
			case "off":
				cis.GainAmplifierEnabled(false)
			default:
				gain, err := strconv.Atoi(args[0])
				if err != nil {
					return err
				}
				return cis.GainAmplifierLevel(gain)
			}

			return nil
		},
	}

	dumpreg := &ffcli.Command{
		Name:       "dumpreg",
		ShortUsage: "kd6ctl dumpreg",
		ShortHelp:  "Dump the register values of CIS.",
		Exec: func(_ context.Context, args []string) error {

			cis := kd6rmx.Sensor{Port: *port}
			err := cis.ReadRegister("BR")
			err = cis.ReadRegister("OF")
			err = cis.ReadRegister("OC")
			err = cis.ReadRegisterWithVal("OC", "C0")
			err = cis.ReadRegister("RC")
			err = cis.ReadRegister("SS")
			err = cis.ReadRegister("DC")
			err = cis.ReadRegister("LC")
			err = cis.ReadRegisterWithVal("LC", "A0")
			err = cis.ReadRegisterWithVal("LC", "C0")
			err = cis.ReadRegister("WC")
			// cis.ReadRegister("PG")
			// cis.ReadRegister("GC")
			// cis.ReadRegister("TP")

			return err
		},
	}

	cmd := &ffcli.Command{
		Name:       "cmd",
		ShortUsage: "kd6ctl cmd <register> <value>",
		ShortHelp:  "Sends the specified command to sensor",
		Exec: func(_ context.Context, args []string) error {
			if len(args) < 2 {
				return fmt.Errorf("cmd command requires specific register and value you want to send")
			}

			register := args[0]
			command := args[1]

			cis := kd6rmx.Sensor{Port: *port}
			cis.SendCommand(register, command)
			return nil
		},
	}

	root := &ffcli.Command{
		ShortUsage:  "kd6ctl [flags] <subcommand>",
		ShortHelp:   "kd6ctl is a command line utility to change config on the KD6RMX contact image sensor.",
		FlagSet:     rootFlagSet,
		Subcommands: []*ffcli.Command{version, dumpreg, gain, load, save, outputfreq, outputfmt, interp, dark, white, leds, duty, cmd},
		Exec: func(context.Context, []string) error {
			return flag.ErrHelp
		},
	}

	if err := root.ParseAndRun(context.Background(), os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
