package main

import (
	"fmt"
	mp "github.com/mackerelio/go-mackerel-plugin"
	"os"
	"os/exec"
	"syscall"
)

func checkDeviceStatus(dev string) (int, error) {
	cmd := exec.Command("smartctl", "-a", dev)
	if err := cmd.Run(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				return status.ExitStatus(), nil
			}
		}
		return -1, err
	}
	return 0, nil
}

type Plugin struct {
}

func (p *Plugin) FetchMetrics() (map[string]float64, error) {
	ret := make(map[string]float64)
	for c := 'a'; c <= 'z'; c++ {
		dev := fmt.Sprintf("/dev/sd%c", c)
		if _, err := os.Stat(dev); err != nil {
			continue
		}
		status, err := checkDeviceStatus(dev)
		if err != nil {
			return nil, err
		}
		ret[fmt.Sprintf("diskhealth.sd%c.errors", c)] = float64(status)
	}
	return ret, nil
}

func (p *Plugin) GraphDefinition() map[string]mp.Graphs {
	return map[string]mp.Graphs{
		"diskhealth.#": {
			Label: "smartctl exit status",
			Unit:  mp.UnitInteger,
			Metrics: []mp.Metrics{
				{Name: "errors", Label: "Errors"},
			},
		},
	}
}

func main() {
	p := mp.NewMackerelPlugin(&Plugin{})
	p.Run()
}
