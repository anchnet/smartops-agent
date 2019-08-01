package system

import (
	"fmt"
	"github.com/DataDog/gopsutil/cpu"
	"github.com/DataDog/gopsutil/process"
	"github.com/anchnet/smartops-agent/pkg/metric"
	"github.com/shirou/gopsutil/mem"
	"strconv"

	"strings"
	"time"
)

var (
	lastProcs       map[int32]*process.FilledProcess
	lastProcCPUTime cpu.TimesStat
)

const (
	procMetric = "system.proc.%s"
)

func runProcCheck(t time.Time) ([]metric.MetricSample, error) {
	var samples []metric.MetricSample

	cpuTimes, err := cpu.Times(false)
	if err != nil {
		return samples, err
	}
	procs, err := process.AllProcesses()
	if err != nil {
		return samples, err
	}
	if lastProcs == nil {
		lastProcs = procs
		lastProcCPUTime = cpuTimes[0]
		return samples, nil
	}
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return samples, err
	}

	totalMem := memInfo.Total
	for _, p := range procs {

		if _, ok := lastProcs[p.Pid]; !ok {
			continue
		}
		tag := make(map[string]string, 4)
		tag["pid"] = strconv.Itoa(int(p.Pid))
		tag["name"] = p.Name
		tag["user"] = p.Username
		tag["cmdline"] = strings.Join(p.Cmdline, "")

		cpuUsage, cpuUser, cpuSys := formatCPU(p.CpuTime, lastProcs[p.Pid].CpuTime, cpuTimes[0], lastProcCPUTime)

		samples = append(samples, metric.NewServerMetricSample(fmt.Sprintf(procMetric, "cpu.usage"), float64(cpuUsage), metric.UnitByte, t, tag))
		samples = append(samples, metric.NewServerMetricSample(fmt.Sprintf(procMetric, "cpu.user"), float64(cpuUser), metric.UnitByte, t, tag))
		samples = append(samples, metric.NewServerMetricSample(fmt.Sprintf(procMetric, "cpu.system"), float64(cpuSys), metric.UnitByte, t, tag))
		samples = append(samples, metric.NewServerMetricSample(fmt.Sprintf(procMetric, "mem.rss"), float64(p.MemInfo.RSS), metric.UnitByte, t, tag))
		samples = append(samples, metric.NewServerMetricSample(fmt.Sprintf(procMetric, "mem.vms"), float64(p.MemInfo.VMS), metric.UnitByte, t, tag))
		samples = append(samples, metric.NewServerMetricSample(fmt.Sprintf(procMetric, "mem.pct"), float64(p.MemInfo.VMS/totalMem*100), metric.UnitPercent, t, tag))
		samples = append(samples, metric.NewServerMetricSample(fmt.Sprintf(procMetric, "mem.pct"), float64(p.MemInfo.VMS/totalMem*100), metric.UnitByte, t, tag))
		samples = append(samples, metric.NewServerMetricSample(fmt.Sprintf(procMetric, "thread.count"), float64(p.NumThreads), "", t, tag))

	}

	return samples, nil
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
