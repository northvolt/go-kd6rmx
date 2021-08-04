# go-kd6rmx

Go language command interface to the Mitsubishi Electric KD-MX series of contact image sensors.

http://www.mitsubishielectric.com/bu/contact_image/general/lineup/index.html

For information about this sensor contact Mitsubishi Electric.

## API

The API is meant to be used in Go programs. You would normally add it to your own project using Go modules like this:

```shell
go get github.com/northvolt/go-kd6rmx
```

```go
import "github.com/northvolt/go-kd6rmx"

...

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

### How to install

First you must obtain the git repo, and change into the new directory:

```shell
git clone https://github.com/northvolt/go-kd6rmx.git
cd go-kd6rmx
```

Now you can install the CLI

```shell
go install ./cmd/kd6ctl
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
  dark           Dark correction on/off/adjust.
  white          White correction on/off/adjust/target.
  led            Sets LEDs on sensor on or off.

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
kd6ctl led ab on
```

### How to build binaries

```shell
make build
```

### How to install binaries

```shell
make install
```

### How to tag and release binaries

```shell
make release bump=major
```
The automatic version format is vX.Y.Z
bump=major will create a tag with increament in X from previous tag.
bump=minor will create a tag with increament in Y from previous tag.
bump=which shows the current tag.
If bump is not specified, a tag with increment in Z is created.