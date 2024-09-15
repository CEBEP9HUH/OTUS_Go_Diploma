package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/CEBEP9HUH/OTUS_Go_Diploma/pkg/api/grpc/statistic/broadcaster"
)

const (
	timestampFormat = "15:04:05"
	lString         = "%-15.15s"
	rNumber         = "%15.2f"
)

func printHeader(heads ...string) string {
	var builder strings.Builder
	for _, h := range heads {
		builder.WriteString("|")
		builder.WriteString(fmt.Sprintf(lString, h))
	}
	builder.WriteString("|")
	header := builder.String()
	line := strings.Repeat("-", len(header))
	fmt.Println(header)
	fmt.Println(line)
	return line
}

func printLine(vals ...any) {
	var builder strings.Builder
	for _, val := range vals {
		switch v := val.(type) {
		case string:
			builder.WriteString("|")
			builder.WriteString(fmt.Sprintf(lString, v))
		case float32:
			builder.WriteString("|")
			builder.WriteString(fmt.Sprintf(rNumber, v))
		}
	}
	builder.WriteString("|")
	fmt.Println(builder.String())
}

func printCPUusage(stat <-chan *broadcaster.SysStat) {
	line := printHeader("Timestamp", "User Mode", "System Mode", "Idle")
	for snap := range stat {
		rec := snap.CpuUsage
		printLine(time.Now().Format(timestampFormat), rec.UserMode, rec.SystemMode, rec.Idle)
		fmt.Println(line)
	}
}

func printIOStat(stat <-chan *broadcaster.SysStat) {
	line := printHeader("Timestamp", "Device", "TPS", "Read (Kb/s)", "Write (Kb/s)")
	for snap := range stat {
		for dev, info := range snap.IoStat.Usage {
			printLine(time.Now().Format(timestampFormat), dev, info.Tps, info.ReadKbs, info.WriteKbs)
		}
		fmt.Println(line)
	}
}

func printLoadAvg(stat <-chan *broadcaster.SysStat) {
	line := printHeader("Timestamp", "Min1", "Min5", "Min15")
	for snap := range stat {
		rec := snap.LoadAvg
		printLine(time.Now().Format(timestampFormat), rec.Min1, rec.Min5, rec.Min15)
		fmt.Println(line)
	}
}

func printDiskUsage(stat <-chan *broadcaster.SysStat) {
	line := printHeader("Timestamp", "Mounted on", "FS", "Disk Usage (%)")
	for snap := range stat {
		for mount, info := range snap.DiskUsage.Usage {
			printLine(time.Now().Format(timestampFormat), mount, info.Fs, info.Percent)
		}
		fmt.Println(line)
	}
}

func printNodeUsage(stat <-chan *broadcaster.SysStat) {
	line := printHeader("Timestamp", "Mounted on", "FS", "Node Usage (%)")
	for snap := range stat {
		for mount, info := range snap.NodeUsage.Usage {
			printLine(time.Now().Format(timestampFormat), mount, info.Fs, info.Percent)
		}
		fmt.Println(line)
	}
}
