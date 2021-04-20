# go-kd6rmx

Go language command interface to the Mitsubishi Electric KD-MX series of contact image sensors.

http://www.mitsubishielectric.com/bu/contact_image/general/lineup/index.html

For information about this sensor contact Mitsubishi Electric.

## How to use:

```go
	cis := kd6rmx.Sensor{Port: "/dev/your-port-here"}

	// load settings from user preset 1
	return cis.LoadSettings("1")
```
