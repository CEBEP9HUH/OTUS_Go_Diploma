//go:build linux
// +build linux

package collectorsmanager

import (
	sc "github.com/CEBEP9HUH/OTUS_Go_Diploma/internal/services/collector"
	"github.com/CEBEP9HUH/OTUS_Go_Diploma/internal/services/collector/linux"
)

func getCollectorsMap(l CollectorsList) []sc.StatisticCollector {
	var collectors []sc.StatisticCollector
	if l.EnableCPUUsage {
		collectors = append(collectors, linux.MakeCPUUsageStatCollector())
	}
	if l.EnableDiskUsage {
		collectors = append(collectors, linux.MakeDiskUsageStatCollector())
	}
	if l.EnableNodeUsage {
		collectors = append(collectors, linux.MakeNodeUsageStatCollector())
	}
	if l.EnableLoadAvg {
		collectors = append(collectors, linux.MakeLoadAvgStatCollector())
	}
	if l.EnableIOStat {
		collectors = append(collectors, linux.MakeIOStatStatCollector())
	}
	return collectors
}
