package kd6rmx

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

// Sensor is a wrapper for control functions for the KD6RMX contact image sensor.
type Sensor struct {
	Port string
}

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

	_, err := cis.sendCommand("BR", param)
	return err
}

func (cis Sensor) SetOutputFrequency(freq float32) error {
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

	result, err := cis.sendCommand("OF", val)
	if err != nil {
		return err
	}
	if len(result) < 5 {
		return errors.New("invalid result from SetOutputFrequency")
	}
	return nil
}

func (cis Sensor) PixelOutputFormat(conf string) error {
	result, err := cis.sendCommand("OC", conf)
	if err != nil {
		return err
	}
	if len(result) < 5 {
		return errors.New("invalid result from PixelOutputFormat")
	}
	return nil
}

func (cis Sensor) Resolution(res int) error {
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

	result, err := cis.sendCommand("RC", param)
	if err != nil {
		return err
	}
	if len(result) < 5 {
		return errors.New("invalid result from Resolution")
	}
	return nil
}

func (cis Sensor) ExternalSync() error {
	result, err := cis.sendCommand("SS", "01")
	if err != nil {
		return err
	}
	if len(result) < 5 {
		return errors.New("invalid result from ExternalSync")
	}
	return nil
}

func (cis Sensor) InternalSync(upper, lower string) error {
	if len(upper) != 2 || len(lower) != 2 {
		errors.New("invalid upper sync value")
	}

	result, err := cis.sendCommand("SS", "00"+upper+lower)
	if err != nil {
		return err
	}

	if len(result) < 5 {
		return errors.New("invalid result from set sync")
	}
	return nil
}

func (cis Sensor) LoadSettings(preset string) error {
	var val string
	switch preset {
	case "0", "1", "2", "3":
		val = "0" + preset
	default:
		return errors.New("invalid preset for LoadSettings: " + preset)
	}

	result, err := cis.sendCommand("DT", val)
	if err != nil {
		return err
	}
	if len(result) < 5 {
		return errors.New("invalid result from LoadSettings")
	}
	if string(result[3:4]) != preset {
		return errors.New("invalid result code from LoadSettings: " + string(result))
	}
	return nil
}

func (cis Sensor) SaveSettings(preset string) error {
	var val string
	switch preset {
	case "1", "2", "3":
		val = "8" + preset
	default:
		return errors.New("invalid preset for SaveSettings: " + preset)
	}

	result, err := cis.sendCommand("DT", val)
	if err != nil {
		return err
	}
	if len(result) < 5 {
		return errors.New("invalid result from SaveSettings")
	}
	if string(result[2:4]) != preset {
		return errors.New("invalid result code from SaveSettings: " + string(result))
	}
	return nil
}

func (cis Sensor) DarkCorrectionEnabled(on bool) error {
	var param = "00"
	if on {
		param = "01"
	}

	result, err := cis.sendCommand("DC", param)
	if err != nil {
		return err
	}
	if len(result) < 5 {
		return errors.New("invalid result from DarkCorrection")
	}
	return nil
}

func (cis Sensor) LightControl(a, b bool) error {
	var param string
	switch {
	case a && b:
		param = "03"
	case a:
		param = "01"
	case b:
		param = "02"
	default:
		param = "00"
	}

	result, err := cis.sendCommand("LC", param)
	if err != nil {
		return err
	}
	if len(result) < 5 {
		return errors.New("invalid result from LightControl")
	}
	return nil
}

func (cis Sensor) WhiteCorrectionEnabled(on bool) error {
	var param = "00"
	if on {
		param = "01"
	}

	result, err := cis.sendCommand("WC", param)
	if err != nil {
		return err
	}
	if len(result) < 5 {
		return errors.New("invalid result from WhiteCorrection")
	}
	return nil
}

func (cis Sensor) PerformWhiteCorrection() error {
	result, err := cis.sendCommand("WC", "21")
	if err != nil {
		return err
	}
	if len(result) < 5 {
		return errors.New("invalid result from WhiteCorrection")
	}
	return nil
}

func (cis Sensor) WhiteCorrectionTarget(target int) error {
	if target > 255 {
		return errors.New("invalid white correction target")
	}
	h := fmt.Sprintf("%0X", target)

	result, err := cis.sendCommand("WC40", h)
	if err != nil {
		return err
	}
	if len(result) < 5 {
		return errors.New("invalid result from WhiteCorrectionTarget")
	}
	return nil
}

func (cis Sensor) GainAmplifierEnabled(on bool) error {
	var param = "00"
	if on {
		param = "01"
	}

	_, err := cis.sendCommand("PG", param)
	return err
}

func (cis Sensor) GainAmplifierLevel(gain int) error {
	if gain > 3071 {
		return errors.New("invalid gain level")
	}

	result, err := cis.sendCommand("PG20", fmt.Sprintf("%000X", gain))
	if err != nil {
		return err
	}
	if len(result) < 5 {
		return errors.New("invalid result from WhiteCorrectionTarget")
	}

	return nil
}

func (cis Sensor) YCorrectionEnabled(on bool) error {
	var param = "00"
	if on {
		param = "01"
	}

	result, err := cis.sendCommand("GC", param)
	if err != nil {
		return err
	}
	if len(result) < 5 {
		return errors.New("invalid result from YCorrection")
	}
	return nil
}

func (cis Sensor) TestPatternEnabled(on bool) error {
	var param = "00"
	if on {
		param = "01"
	}

	result, err := cis.sendCommand("TP", param)
	if err != nil {
		return err
	}
	if len(result) < 5 {
		return errors.New("invalid result from TestPattern")
	}
	return nil
}

func (cis Sensor) InterpolationEnabled(on bool) error {
	var param = "40"
	if on {
		param = "41"
	}

	result, err := cis.sendCommand("OC", param)
	if err != nil {
		return err
	}
	if len(result) < 5 {
		return errors.New("invalid result from Interpolation")
	}
	return nil
}

func (cis Sensor) SoftwareReset() error {
	result, err := cis.sendCommand("SR", "21")
	if err != nil {
		return err
	}
	if len(result) < 5 {
		return errors.New("invalid result from SoftwareReset")
	}

	time.Sleep(10 * time.Second)

	result, err = cis.sendCommand("SR", "01")
	if err != nil {
		return err
	}
	if len(result) < 5 {
		return errors.New("invalid result from SoftwareReset")
	}

	return nil
}

func (cis Sensor) sendCommand(cmd string, params string) (string, error) {
	f, err := os.OpenFile(cis.Port, os.O_RDWR|os.O_APPEND, 0777)
	if err != nil {
		return "", fmt.Errorf("error opening control port: %v", err)
	}
	defer f.Close()

	_, err = f.Write([]byte(cmd + params + "\r"))
	if err != nil {
		return "", fmt.Errorf("error sending command: %v", err)
	}

	buf := make([]byte, 5)
	_, err = io.ReadFull(f, buf)
	if err != nil {
		return "", fmt.Errorf("error receiving result from command: %v", err)
	}

	return string(buf), nil
}
