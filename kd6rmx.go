package kd6rmx

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// Sensor is a wrapper for control functions for the KD6RMX contact image sensor.
type Sensor struct {
	Port        string
	Logging     bool
	FileLogging bool
}

// CommunicationSpeed sets the communcation speed.
func (cis Sensor) CommunicationSpeed(baud int) error {
	var param string
	switch baud {
	case 9600:
		param = "00"
	case 19200:
		param = "01"
	case 115200:
		param = "02"
	default:
		return errors.New("invalid baud rate")
	}

	_, err := cis.SendCommand("BR", param)
	return err
}

// OutputFrequency sets the output frequency.
func (cis Sensor) OutputFrequency(freq float32) error {
	var val string
	switch freq {
	case 48.0:
		val = "00"
	case 50.7:
		val = "01"
	case 51.0:
		val = "02"
	case 51.4:
		val = "03"
	case 52.0:
		val = "04"
	case 52.8:
		val = "05"
	case 53.3:
		val = "06"
	case 54.0:
		val = "07"
	case 54.9:
		val = "08"
	case 56.0:
		val = "09"
	case 57.0:
		val = "0A"
	case 57.6:
		val = "0B"
	case 58.3:
		val = "0C"
	case 60.0:
		val = "0D"
	case 61.7:
		val = "0E"
	case 62.4:
		val = "0F"
	case 64.0:
		val = "10"
	case 65.1:
		val = "11"
	case 66.0:
		val = "12"
	case 67.2:
		val = "13"
	case 68.0:
		val = "14"
	case 68.6:
		val = "15"
	case 72.0:
		val = "16"
	case 76.0:
		val = "17"
	case 76.8:
		val = "18"
	case 78.0:
		val = "19"
	case 80.0:
		val = "1A"
	case 81.6:
		val = "1B"
	case 84.0:
		val = "1C"
	default:
		return errors.New("invalid output frequency")
	}

	result, err := cis.SendCommand("OF", val)
	if err != nil {
		return err
	}
	return checkError("SetOutputFrequency", result)
}

type PixelOutputBits int
type PixelOutputInterface int
type PixelOutputConfig int

const (
	PixelOutputBits10 PixelOutputBits = iota
	PixelOutputBits8
)

const (
	PixelOutputSerial PixelOutputInterface = iota
	PixelOutputParallel
)

const (
	PixelOutputBase PixelOutputConfig = iota
	PixelOutputMedium
)

// PixelOutputFormat sets the output format for pixels.
//
// For example:
//
//		cis.PixelOutputFormat(kd6rmx.PixelOutputBits10, kd6rmx.PixelOutputSerial, kd6rmx.PixelOutputBase, 1)
//
func (cis Sensor) PixelOutputFormat(bits PixelOutputBits, i PixelOutputInterface, conf PixelOutputConfig, number int) error {
	var param string
	switch bits {
	case PixelOutputBits10:
		switch i {
		case PixelOutputSerial:
			switch conf {
			case PixelOutputBase:
				if number == 1 {
					param = "00"
				}
			case PixelOutputMedium:
				switch number {
				case 1:
					param = "01"
				case 2:
					param = "02"
				case 3:
					param = "03"
				}
			}

		case PixelOutputParallel:
			switch conf {
			case PixelOutputBase:
				if number == 1 {
					param = "04"
				}
			case PixelOutputMedium:
				switch number {
				case 1:
					param = "05"
				case 2:
					param = "06"
				case 3:
					param = "07"
				}
			}

		}
	case PixelOutputBits8:
		switch i {
		case PixelOutputSerial:
			switch conf {
			case PixelOutputBase:
				if number == 1 {
					param = "08"
				}
			case PixelOutputMedium:
				switch number {
				case 1:
					param = "09"
				case 2:
					param = "0A"
				case 3:
					param = "0B"
				}
			}

		case PixelOutputParallel:
			switch conf {
			case PixelOutputBase:
				if number == 1 {
					param = "0C"
				}
			case PixelOutputMedium:
				switch number {
				case 1:
					param = "0D"
				case 2:
					param = "0E"
				case 3:
					param = "0F"
				}
			}
		}
	}

	if param == "" {
		return errors.New("invalid params for PixelOutputFormat")
	}

	result, err := cis.SendCommand("OC", param)
	if err != nil {
		return err
	}
	return checkError("PixelOutputFormat", result)
}

// PixelOverlap turns on/off the pixel overlap. Only to be used on CIS with 2 or 3 sensors.
func (cis Sensor) PixelOverlap(on bool) error {
	var param = "20"
	if on {
		param = "21"
	}
	result, err := cis.SendCommand("OC", param)
	if err != nil {
		return err
	}
	return checkError("PixelOverlap", result)
}

// PixelInterpolation turns on/off pixel interpolation.
func (cis Sensor) PixelInterpolation(on bool) error {
	var param = "40"
	if on {
		param = "41"
	}

	result, err := cis.SendCommand("OC", param)
	if err != nil {
		return err
	}
	return checkError("PixelInterpolation", result)
}

// PixelResolution sets the resolution for the sensor.
// Valid resolutions are 600, 300, 150, or 75 dpi.
func (cis Sensor) PixelResolution(res int) error {
	var param string
	switch res {
	case 600:
		param = "00"
	case 300:
		param = "01"
	case 150:
		param = "02"
	case 75:
		param = "03"
	default:
		return errors.New("invalid resolution")
	}

	result, err := cis.SendCommand("RC", param)
	if err != nil {
		return err
	}
	return checkError("PixelResolution", result)
}

// ExternalSync turns on the external sync.
func (cis Sensor) ExternalSync() error {
	result, err := cis.SendCommand("SS", "01")
	if err != nil {
		return err
	}
	return checkError("ExternalSync", result)
}

// InternalSync turns on the internal sync.
func (cis Sensor) InternalSync(val int) error {
	if val < 1 || val > 0xffff {
		errors.New("invalid sync clock value")
	}

	param := fmt.Sprintf("00%000X", val)
	result, err := cis.SendCommand("SS", param)
	if err != nil {
		return err
	}

	if len(result) < 5 {
		return errors.New("invalid result from set sync")
	}
	return nil
}

// LoadSettings loads the sensor's active settings with one of the memory presets.
//
// Valid presets are:
// 		0 (factory defaults)
// 		1 (user settings 1)
// 		2 (user settings 2)
// 		3 (user settings 3)
func (cis Sensor) LoadSettings(preset int) error {
	if preset < 0 || preset > 4 {
		return errors.New("invalid preset for LoadSettings")
	}

	param := fmt.Sprintf("%02X", preset)
	result, err := cis.SendCommand("DT", param)
	if err != nil {
		return err
	}
	return checkError("LoadSettings", result)
}

// SaveSettings saves the sensor's active settings to one of the memory presets.
//
// Valid presets are:
// 		1 (user settings 1)
// 		2 (user settings 2)
// 		3 (user settings 3)
func (cis Sensor) SaveSettings(preset int) error {
	if preset < 1 || preset > 4 {
		return errors.New("invalid preset for SaveSettings")
	}

	param := fmt.Sprintf("%02X", 0x80+preset)
	result, err := cis.SendCommand("DT", param)
	if err != nil {
		return err
	}
	return checkError("SaveSettings", result)
}

// LEDControl sets the LED settings. You can do the following:
// Turn on/off the A and B LEDS individually.
// Set the pulse divider to be used for both LEDS so they are:
//		1 = on all of the time
//		2 = on one half of the time
//		4 = on one quarter of the time
//		8 = on one eighth of the time
func (cis Sensor) LEDControl(leds string, on bool, pulsedivider int) error {
	var pd int
	switch pulsedivider {
	case 1:
		pd = 0
	case 2:
		pd = 1
	case 4:
		pd = 2
	case 8:
		pd = 3
	default:
		return errors.New("pulsedivider must be 1, 2, 4, or 8")
	}

	var val int
	if on {
		switch leds {
		case "ab", "AB":
			val = 3 + (pd * 4)
		case "a", "A":
			val = 1 + (pd * 4)
		case "b", "B":
			val = 2 + (pd * 4)
		default:
			return errors.New("invalid LEDs, must be 'A', 'B', or 'AB'")
		}
	} else {
		// turn both off.
		// TODO: turn back on anything that was on and we did not explicitly turn off.
		val = 0
	}

	param := fmt.Sprintf("%02X", val)
	result, err := cis.SendCommand("LC", param)
	if err != nil {
		return err
	}
	return checkError("LEDControl", result)
}

// LEDDutyCycle sets the duty cycle for each LED separately.
// The value for duty represents the raw value of register LC.
func (cis Sensor) LEDDutyCycle(led string, duty int) error {
	var ls string
	switch led {
	case "a", "A":
		ls = "20"
	case "b", "B":
		ls = "40"
	default:
		return errors.New("invalid LED for duty cycle")
	}

	if duty <= 0 || duty >= 4096 {
		return errors.New("invalid duty cycle register value")
	}
	param := fmt.Sprintf("%s%04X", ls, duty)

	result, err := cis.SendCommand("LC", param)
	if err != nil {
		return err
	}
	return checkError("LEDDutyCycle", result)
}

func (cis Sensor) DarkCorrectionEnabled(on bool) error {
	var param = "00"
	if on {
		param = "01"
	}

	result, err := cis.SendCommand("DC", param)
	if err != nil {
		return err
	}
	return checkError("DarkCorrectionEnabled", result)
}

func (cis Sensor) PerformDarkCorrection() error {
	result, err := cis.SendCommand("DC", "21")
	if err != nil {
		return err
	}
	return checkError("PerformDarkCorrection", result)
}

func (cis Sensor) WhiteCorrectionEnabled(on bool) error {
	var param = "00"
	if on {
		param = "01"
	}

	result, err := cis.SendCommand("WC", param)
	if err != nil {
		return err
	}
	return checkError("WhiteCorrectionEnabled", result)
}

func (cis Sensor) PerformWhiteCorrection() error {
	result, err := cis.SendCommand("WC", "21")
	if err != nil {
		return err
	}
	return checkError("PerformWhiteCorrection", result)
}

func (cis Sensor) WhiteCorrectionTarget(target int) error {
	if target > 255 {
		return errors.New("invalid white correction target")
	}

	h := fmt.Sprintf("%04X", target*16)

	result, err := cis.SendCommand("WC40", h)
	if err != nil {
		return err
	}
	// sleep for one second to make sure correction finished
	time.Sleep(1 * time.Second)
	return checkError("WhiteCorrectionTarget", result)
}

func (cis Sensor) GainAmplifierEnabled(on bool) error {
	var param = "00"
	if on {
		param = "01"
	}

	_, err := cis.SendCommand("PG", param)
	return err
}

func (cis Sensor) GainAmplifierLevel(gain int) error {
	var param string
	switch {
	case gain > 3071:
		return errors.New("invalid positive gain level")
	case gain >= 0:
		// positive gain
		param = fmt.Sprintf("20%04X", gain)
	case gain < -1027:
		return errors.New("invalid negative gain level")
	case gain < 0:
		// negative gain
		param = fmt.Sprintf("21%04X", gain)
	}

	result, err := cis.SendCommand("PG", param)
	if err != nil {
		return err
	}
	return checkError("GainAmplifierLevel", result)
}

func (cis Sensor) YCorrectionEnabled(on bool) error {
	var param = "00"
	if on {
		param = "01"
	}

	result, err := cis.SendCommand("OC", param)
	if err != nil {
		return err
	}
	return checkError("YCorrectionEnabled", result)
}

func (cis Sensor) TestPatternEnabled(on bool) error {
	var param = "00"
	if on {
		param = "01"
	}

	result, err := cis.SendCommand("TP", param)
	if err != nil {
		return err
	}

	return checkError("TestPatternEnabled", result)
}

type TestPatternType int

const (
	TestPatternStripe TestPatternType = iota
	TestPatternRamp
)

func (cis Sensor) TestPattern(pattern TestPatternType) error {
	var param = "20"
	if pattern == TestPatternRamp {
		param = "21"
	}

	result, err := cis.SendCommand("TP", param)
	if err != nil {
		return err
	}
	return checkError("TestPattern", result)
}

func (cis Sensor) SoftwareReset() error {
	result, err := cis.SendCommand("SR", "21")
	if err != nil {
		return err
	}
	if len(result) < 4 {
		return errors.New("invalid result from SoftwareReset")
	}

	time.Sleep(10 * time.Second)

	result, err = cis.SendCommand("SR", "01")
	if err != nil {
		return err
	}
	return checkError("SoftwareReset", result)
}

func (cis Sensor) SendCommand(cmd string, params string) (string, error) {

	f, err := os.OpenFile(cis.Port, os.O_RDWR|os.O_APPEND, 0777)
	if err != nil {
		return "", fmt.Errorf("error opening control port: %v", err)
	}
	defer f.Close()

	write_string := cmd + params + "\r"
	dt_string := time.Now().Format("2006-01-02 15:04:05.000000000")

	if cis.FileLogging {
		f_log, err := os.OpenFile("kd6cmd.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return "", fmt.Errorf("cannot write to log file: %v", err)
		}
		defer f_log.Close()

		if _, err := f_log.WriteString(dt_string + "\n"); err != nil {
			log.Println(err)
		}
		if _, err := f_log.WriteString("send: " + write_string + "\n"); err != nil {
			log.Println(err)
		}
	}

	if cis.Logging {
		fmt.Printf("%s\n", dt_string)
		fmt.Println("send: ", write_string)
	}

	_, err = f.Write([]byte(write_string))
	if err != nil {
		return "", fmt.Errorf("error sending command: %v", err)
	}

	buf := make([]byte, 5)
	var result string
	start := time.Now()
	for {

		n, err := f.Read(buf)
		if err != nil {
			if err == io.EOF {
				if time.Since(start) > time.Second*10 {
					return "", fmt.Errorf("timeout receiving result from command")
				}
				continue
			}

			// some other error
			return "", err
		}

		result += string(buf[:n])
		switch {
		case result[len(result)-1] == '\r':
			result = strings.Replace(result, "\r", "", -1)

			dt_string_rec := time.Now().Format("2006-01-02 15:04:05.000000000")
			if cis.FileLogging {
				f_log, err := os.OpenFile("kd6cmd.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					return "", fmt.Errorf("cannot write to log file: %v", err)
				}
				defer f_log.Close()
				if _, err := f_log.WriteString(dt_string_rec + "\n"); err != nil {
					log.Println(err)
				}
				if _, err := f_log.WriteString("received: " + write_string + "\n\n"); err != nil {
					log.Println(err)
				}
			}

			if cis.Logging {
				fmt.Printf("%s\n", dt_string_rec)
				fmt.Printf("received: %s\n\n", result)

			}
			return result, nil
		case n == 0:
			return "", fmt.Errorf("no data in result from command")
		case time.Since(start) > time.Second*10:
			return "", fmt.Errorf("timeout receiving result from command")
		}
	}
}

func checkError(funcname, result string) error {
	if len(result) < 4 {
		return fmt.Errorf("invalid result from %s: %s", funcname, result)
	}

	return nil
}

func (cis Sensor) ReadRegister(register string) error {
	return cis.ReadRegisterWithVal(register, "80")
}

func parseRate(short_res string) error {
	var val float32
	switch short_res {
	case "00":
		val = 9.6
	case "01":
		val = 19.2
	case "02":
		val = 115.2
	default:
		return errors.New("invalid output Baud Rate")
	}
	fmt.Printf("(UART Setting: Baud rate: %.1f MHz)\n", val)
	return nil
}

func parseFreq(short_res string) error {
	var val float32
	switch short_res {
	case "00":
		val = 48
	case "01":
		val = 50.7
	case "02":
		val = 51.0
	case "03":
		val = 51.4
	case "04":
		val = 52.0
	case "05":
		val = 52.8
	case "06":
		val = 53.3
	case "07":
		val = 54.0
	case "08":
		val = 54.9
	case "09":
		val = 56.0
	case "0A":
		val = 57.0
	case "0B":
		val = 57.6
	case "0C":
		val = 58.3
	case "0D":
		val = 60.0
	case "0E":
		val = 62.7
	case "0F":
		val = 62.4
	case "10":
		val = 64.0
	case "11":
		val = 65.1
	case "12":
		val = 66.0
	case "13":
		val = 67.2
	case "14":
		val = 68.0
	case "15":
		val = 68.6
	case "16":
		val = 72.0
	case "17":
		val = 76.0
	case "18":
		val = 76.8
	case "19":
		val = 78.0
	case "1A":
		val = 80.0
	case "1B":
		val = 81.6
	case "1C":
		val = 84.0
	default:
		return errors.New("invalid output frequency")
	}
	fmt.Printf("(Frequency: %.1f MHz)\n", val)
	return nil
}

func parseComm(short_res, val string) error {
	switch val {
	case "80":
		var val string
		switch short_res {
		case "00":
			val = "10bit Serial Base Configuration"
		case "01":
			val = "10bit Serial Medium Configuration"
		case "02":
			val = "10bit Serial Medium Configuration2"
		case "03":
			val = "10bit Serial Medium Configuration3"
		case "04":
			val = "10bit Parallel Base Configuration"
		case "05":
			val = "10bit Parallel Medium Configuration"
		case "06":
			val = "10bit Parallel Medium Configuration2"
		case "07":
			val = "10bit Parallel Medium Configuration3"
		case "08":
			val = "8bit Serial Base Configuration"
		case "09":
			val = "8bit Serial Medium Configuration"
		case "0A":
			val = "8bit Serial Medium Configuration2"
		case "0B":
			val = "8bit Serial Medium Configuration3"
		case "0C":
			val = "8bit Parallel Base Configuration"
		case "0D":
			val = "8bit Parallel Medium Configuration"
		case "0E":
			val = "8bit Parallel Medium Configuration2"
		case "0F":
			val = "8bit Parallel Medium Configuration3"
		default:
			return errors.New("invalid output format")
		}
		fmt.Printf("(Output Format: %s)\n", val)

	case "A0":
		var val string
		switch short_res {
		case "20":
			val = "Overlap Output Off"
		case "21":
			val = "Overlap Output On"
		default:
			return errors.New("invalid overlap output")
		}
		fmt.Printf("(%s)\n", val)

	case "C0":
		var val string
		switch short_res {
		case "40":
			val = "Interpolation function Off"
		case "41":
			val = "Interpolation function On"
		default:
			return errors.New("invalid interpolation output")
		}
		fmt.Printf("(%s)\n", val)
	default:
		return errors.New("invalid output format")
	}
	return nil
}
func parseRes(short_res string) error {
	var val int32
	switch short_res {
	case "00":
		val = 600
	case "01":
		val = 300
	case "02":
		val = 150
	case "03":
		val = 75
	default:
		return errors.New("invalid resolution")
	}
	fmt.Printf("(Resolution: %ddpi)\n", val)
	return nil
}

func parseSync(short_res string) error {
	var val string
	switch short_res {
	case "00":
		val = "Internal synchronization"
	case "01":
		val = "External synchronization"
	default:
		return errors.New("invalid synchronization mode")
	}
	fmt.Printf("(%s)\n", val)
	return nil
}

func parseLED(short_res, result, val string) error {
	switch val {
	case "80":
		var val string
		switch short_res {
		case "00":
			val = "Pulse1: OFF"
		case "01":
			val = "Pulse1: illumination A ON"
		case "02":
			val = "Pulse1: illumination B ON"
		case "03":
			val = "Pulse1: A and B ON"
		case "04":
			val = "Pulse2: OFF"
		case "05":
			val = "Pulse2: illumination A ON"
		case "06":
			val = "Pulse2: illumination B ON"
		case "07":
			val = "Pulse2: A and B ON"
		case "08":
			val = "Pulse4: OFF"
		case "09":
			val = "Pulse4: illumination A ON"
		case "0A":
			val = "Pulse4: illumination B ON"
		case "0B":
			val = "Pulse4: A and B ON"
		case "0C":
			val = "Pulse8: OFF"
		case "0D":
			val = "Pulse8: illumination A ON"
		case "0E":
			val = "ulse8: illumination B ON"
		case "0F":
			val = "Pulse8: A and B ON"
		default:
			return errors.New("invalid LED config")
		}
		fmt.Printf("(%s)\n", val)

	case "A0":
		if len(result) < 8 {
			fmt.Printf("(LED Period A: %d)", len(result))
			return errors.New("result too short")
		}
		hex := result[4:8]
		value, _ := strconv.ParseInt(hex, 16, 64)
		fmt.Printf("(LED Period A: %d)\n", value)
	case "C0":
		if len(result) < 8 {
			fmt.Printf("(LED Period A: %d)", len(result))
			return errors.New("result too short")
		}
		hex := result[4:8]
		value, _ := strconv.ParseInt(hex, 16, 64)
		fmt.Printf("(LED Period B: %d)\n", value)
	case "E0":
		// var val string
		// switch short_res {
		// case "40":
		// 	val = "Interpolation function Off"
		// case "41":
		// 	val = "Interpolation function On"
		// default:
		// 	return errors.New("invalid interpolation output")
		// }
		fmt.Printf("(period setting %s)\n", short_res)
	default:
		return errors.New("LED control setting")
	}
	return nil
}

func parseDarkCorr(short_res string) error {
	var val string
	switch short_res {
	case "00":
		val = "Black correction OFF"
	case "01":
		val = "Black correction ON"
	default:
		return errors.New("invalid dark correction mode")
	}
	fmt.Printf("(%s)\n", val)
	return nil
}

func parseWhiteCorr(short_res string) error {
	var val string
	switch short_res {
	case "00":
		val = "White correction OFF"
	case "01":
		val = "White correction ON"
	default:
		return errors.New("invalid white correction mode")
	}
	fmt.Printf("(%s)\n", val)
	return nil
}

func parseTestPattern(short_res, val string) error {
	switch val {
	case "80":
		var val string
		switch short_res {
		case "00":
			val = "Image output"
		case "01":
			val = "Test pattern output"
		default:
			return errors.New("output mode")
		}
		fmt.Printf("(%s)\n", val)
	case "A0":
		var val string
		switch short_res {
		case "20":
			val = "Stripe pattern output"
		case "21":
			val = "Ramp pattern output"
		default:
			return errors.New("invalid test pattern")
		}
		fmt.Printf("(%s)\n", val)
	default:
		return errors.New("Test Pattern")

	}

	return nil
}

func (cis Sensor) ReadRegisterWithVal(register, val string) error {
	result, err := cis.SendCommand(register, val)
	if err != nil {
		return errors.New("error: send command failed")
	}
	fmt.Printf("Reading %s register with parameter 0x%s ", register, val)

	if result[:2] != "00" {
		fmt.Printf("Reading FAIL. ")
		return errors.New("error: reading fal")
	}

	fmt.Print("Response from CIS ")

	for index := 0; index < len(result); index += 2 {
		fmt.Printf("0x%s ", result[index:index+2])
	}

	short_res := result[2:4]

	switch register {
	case "BR":
		return parseRate(short_res)

	case "OF":
		return parseFreq(short_res)

	case "OC":
		return parseComm(short_res, val)

	case "RC":
		return parseRes(short_res)

	case "SS":
		return parseSync(short_res)
	case "LC":
		return parseLED(short_res, result, val)

	case "DC":
		return parseDarkCorr(short_res)

	case "WC":
		return parseWhiteCorr(short_res)
	case "TP":
		parseTestPattern(short_res, val)

	default:
		fmt.Printf("(default)")
	}

	return nil
}
