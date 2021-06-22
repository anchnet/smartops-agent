// +build !windows
package system

/*
#include <stdio.h>
#include <unistd.h>
#include <utmp.h>
#include <fcntl.h>

int get_users(char users[10][UT_NAMESIZE], int *len)
{
    struct utmp utbuf;
    int utmpfd;

    if ((utmpfd = open(UTMP_FILE, O_RDONLY)) == -1)
    {
        perror(UTMP_FILE);
        exit(1);
    }

    while (read(utmpfd, &utbuf, sizeof(utbuf)) == sizeof(utbuf))
    {
        if (utbuf.ut_type != USER_PROCESS)
            continue;

        strcpy(users[*len], utbuf.ut_name);
        (*len)++;
    }

    close(utmpfd);
    return 0;
}
*/
import "C"

import (
	"fmt"
	"strconv"

	"time"

	"github.com/anchnet/smartops-agent/pkg/collector/core"
	"github.com/anchnet/smartops-agent/pkg/metric"
	"github.com/shirou/gopsutil/v3/host"
)

type AlarmCheck struct {
	name string
}

func (c *AlarmCheck) Name() string {
	return "alarm"
}

func (c *AlarmCheck) Collect(t time.Time) ([]metric.MetricSample, error) {
	var samples []metric.MetricSample

	var argv [10][32]C.char
	var len C.int
	var golen int
	C.get_users(&argv[0], &len)
	golen = int(len)

	info, err := host.Info()
	if err != nil {
		return nil, err
	}

	tags := make(map[string]string)
	tags["uptime"] = strconv.FormatUint(info.Uptime, 10)
	tags["hostname"] = info.Hostname
	tags["platform"] = info.Platform
	tags["kernelVersion"] = info.KernelVersion
	tags["version"] = ""
	samples = append(samples, metric.NewServerMetricSample(c.formatMetric("info"), 0, "", t, tags))
	samples = append(samples, metric.NewServerMetricSample(c.formatMetric("user_count"), float64(golen), metric.UnitPercent, t, nil))

	return samples, nil
}

func (c AlarmCheck) formatMetric(name string) string {
	format := "system.alarm.%s"
	return fmt.Sprintf(format, name)
}

func init() {
	core.RegisterCheck(&AlarmCheck{
		name: "alarm",
	})
}
