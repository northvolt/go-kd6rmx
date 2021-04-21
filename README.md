# go-kd6rmx

Go language command interface to the Mitsubishi Electric KD-MX series of contact image sensors.

http://www.mitsubishielectric.com/bu/contact_image/general/lineup/index.html

For information about this sensor contact Mitsubishi Electric.

## How to use:

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
