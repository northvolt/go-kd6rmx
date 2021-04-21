package kd6rmx

import (
	"testing"
)

func TestSensor(t *testing.T) {
	cis := Sensor{}
	if cis.Port != "" {
		t.Error("Sensor should not have default port value")
	}
}
