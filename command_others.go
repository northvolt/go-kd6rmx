// +build linux darwin

package kd6rmx

import (
	"fmt"
	"os"
	"strings"
	"time"
)

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
	var result string
	start := time.Now()
	for {
		n, err := f.Read(buf)
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
