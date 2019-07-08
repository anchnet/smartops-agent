package system

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"gitlab.51idc.com/smartops/smartcat-agent/pkg/collector/check"
	"gitlab.51idc.com/smartops/smartcat-agent/pkg/collector/core"
	"time"
)

const cpuCheckName = "cpu"

type CPUCheck struct {
	core.CheckBase
	id        string
	name      string
	interval  time.Duration
	cores     int32
	lastCycle float64
	lastTimes cpu.TimesStat
}

func (c *CPUCheck) Run() error {
	cpuTimes, _ := cpu.Times(false)
	t := cpuTimes[0]
	cpuInfo, _ := cpu.Info()
	c.cores = cpuInfo[0].Cores

	cycle := t.Total() / float64(c.cores)

	if c.lastCycle != 0 {
		toPercent := 100 / (cycle - c.lastCycle)

		user := ((t.User + t.Nice) - (c.lastTimes.User + c.lastTimes.Nice)) / float64(c.cores)

		fmt.Println("cpu.user: ", user*toPercent)
	}
	c.lastCycle = cycle
	c.lastTimes = t
	return nil
}

func (c *CPUCheck) Configure() {
	panic("implement me")
}

func CPUFactory() check.Check {
	return &CPUCheck{
		CheckBase: core.NewCheckBase(cpuCheckName),
	}
}

func init() {
	core.RegisterCheck(cpuCheckName, CPUFactory)
}
