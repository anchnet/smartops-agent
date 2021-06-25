package system

import (
	"fmt"

	"github.com/anchnet/smartops-agent/pkg/collector/core"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/process"

	//"github.com/anchnet/smartops-agent/pkg/collector/core"
	"strconv"
	"strings"
	"time"

	"github.com/anchnet/smartops-agent/pkg/metric"
	"github.com/shirou/gopsutil/v3/mem"
)

type ProcCheck struct {
	name            string
	lastProcs       map[int32]*cpu.TimesStat
	lastProcCPUTime cpu.TimesStat
}

func (c *ProcCheck) Name() string {
	return c.name
}

func (c *ProcCheck) Collect(t time.Time) ([]metric.MetricSample, error) {
	var samples []metric.MetricSample

	cpuTimes, err := cpu.Times(false)
	if err != nil {
		return samples, err
	}
	procs, err := process.Processes()
	if err != nil {
		return samples, err
	}
	if c.lastProcs != nil {
		memInfo, err := mem.VirtualMemory()
		if err != nil {
			return samples, err
		}

		totalMem := memInfo.Total
		for _, p := range procs {

			if _, ok := c.lastProcs[p.Pid]; !ok {
				continue
			}

			tag := make(map[string]string, 4)
			tag["pid"] = strconv.Itoa(int(p.Pid))

			name, _ := p.Name()
			user, _ := p.Username()
			cmdline, _ := p.CmdlineSlice()
			tag["name"] = name
			tag["user"] = user
			tag["cmdline"] = strings.Join(cmdline, "")

			minfo, err := p.MemoryInfo()
			if err != nil {
				continue
			}

			numThreads, err := p.NumThreads()
			if err != nil {
				continue
			}

			nowptime, err := p.Times()
			if err != nil {
				continue
			}
			preptime := c.lastProcs[p.Pid]
			if preptime == nil {
				continue
			}

			cpuUsage, cpuUser, cpuSys := formatCPU(*nowptime, *preptime, cpuTimes[0], c.lastProcCPUTime)

			samples = append(samples, metric.NewServerMetricSample(c.formatMetric("cpu.usage"), float64(cpuUsage), metric.UnitPercent, t, tag))
			samples = append(samples, metric.NewServerMetricSample(c.formatMetric("cpu.user"), float64(cpuUser), metric.UnitPercent, t, tag))
			samples = append(samples, metric.NewServerMetricSample(c.formatMetric("cpu.system"), float64(cpuSys), metric.UnitPercent, t, tag))
			samples = append(samples, metric.NewServerMetricSample(c.formatMetric("mem.rss"), float64(minfo.RSS), metric.UnitByte, t, tag))
			samples = append(samples, metric.NewServerMetricSample(c.formatMetric("mem.vms"), float64(minfo.VMS), metric.UnitByte, t, tag))
			samples = append(samples, metric.NewServerMetricSample(c.formatMetric("mem.pct"), float64(minfo.VMS/totalMem*100), metric.UnitPercent, t, tag))
			samples = append(samples, metric.NewServerMetricSample(c.formatMetric("thread.count"), float64(numThreads), "", t, tag))
		}
	}

	// Store the last state for comparison on the next run
	for _, p := range procs {
		times, _ := p.Times()
		c.lastProcs[p.Pid] = times
	}

	c.lastProcCPUTime = cpuTimes[0]

	return samples, nil
}
func (c ProcCheck) formatMetric(name string) string {
	format := "system.proc.%s"
	return fmt.Sprintf(format, name)
}

func formatCPU(t2, t1, syst2, syst1 cpu.TimesStat) (float32, float32, float32) {

	deltaSys := syst2.Total() - syst1.Total()

	totalPct := calculatePct((t2.User-t1.User)+(t2.System-t1.System), deltaSys)
	userPct := calculatePct(t2.User-t1.User, deltaSys)
	sysPct := calculatePct(t2.System-t1.System, deltaSys)
	return totalPct, userPct, sysPct

}

func calculatePct(deltaProc, deltaTime float64) float32 {
	if deltaTime == 0 {
		return 0
	}
	overalPct := (deltaProc / deltaTime) * 100
	if overalPct > 100 {
		overalPct = 100
	}

	return float32(overalPct)
}

func init() {
	core.RegisterCheck(&ProcCheck{
		lastProcs: make(map[int32]*cpu.TimesStat),
		name:      "proc",
	})
}
