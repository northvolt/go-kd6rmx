// +build windows

package kd6rmx

import (
	"fmt"
	"go.bug.st/serial"
	"strings"
	"time"
)

func (cis Sensor) sendCommand(cmd string, params string) (string, error) {
	// open serial port
	p, err := serial.Open(cis.Port, &serial.Mode{BaudRate: 9600})
	if err != nil {
		return "", err
	}
	defer p.Close()

	_, err = p.Write([]byte(cmd + params + "\r"))
	if err != nil {
		return "", fmt.Errorf("error sending command: %v", err)
	}

	buf := make([]byte, 5)
	var result string
	start := time.Now()
	for {
		n, err := p.Read(buf)
		if err != nil {
			return "", fmt.Errorf("error receiving result from command: %v", err)
		}

		result += string(buf[:n])
		switch {
		case result[len(result)-1] == '\r':
			result = strings.Replace(result, "\r", "", -1)
			return result, nil
		case n == 0:
			return "", fmt.Errorf("no data in result from command")
		case time.Since(start) > time.Second*10:
			return "", fmt.Errorf("timeout receiving result from command")
		}
	}
}
