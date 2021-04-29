# go-kd6rmx

Go language command interface to the Mitsubishi Electric KD-MX series of contact image sensors.

http://www.mitsubishielectric.com/bu/contact_image/general/lineup/index.html

For information about this sensor contact Mitsubishi Electric.

## API

```go
// initialize sensor
cis := kd6rmx.Sensor{Port: "/dev/your-port-here"}

// load settings from user preset 1
return cis.LoadSettings(1)

// change some values 
cis.PixelOverlap(true)
cis.PixelInterpolation(true)
cis.PixelResolution(600)

// save current settings to user preset 2
cis.SaveSettings(2)
```

## CLI

`kd6ctl` is a command line interface tool to allow for user configuration.

### How to build

```shell
go build -o build/kd6ctl ./cmd/kd6ctl
```

### How to use

```shell
kd6ctl help
```

Will output a list of commands:

```shell
USAGE
  kd6ctl [flags] <subcommand>

SUBCOMMANDS
  load           Load user settings.
  save           Save current settings into a user preset.
  frequency      Change output frequency (in Mhz).
  format         Change output format.
  interpolation  Set interpolation on/off.
  dark           Dark correction on/off.
  white          White correction on/off.
  led            Sets LED on sensor.

FLAGS
  -p /dev/corser/XtiumCLMX41_s0  port of KD6RMX sensor to use
```

How to set params:

```shell
kd6ctl frequency 60.0
kd6ctl format 10 serial base
kd6ctl interpolation on
kd6ctl dark on
kd6ctl white on
```
