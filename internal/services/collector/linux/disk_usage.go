//go:build linux
// +build linux

package linux

import (
	"context"
	"fmt"

	"github.com/CEBEP9HUH/OTUS_Go_Diploma/internal/services/cmdexecutor"
	"github.com/CEBEP9HUH/OTUS_Go_Diploma/internal/services/cmdexecutor/standart"
	"github.com/CEBEP9HUH/OTUS_Go_Diploma/internal/services/collector"
	"github.com/CEBEP9HUH/OTUS_Go_Diploma/internal/services/statistic"
	"github.com/CEBEP9HUH/OTUS_Go_Diploma/internal/services/util"
)

type diskUsageLinux struct {
	cmd        cmdexecutor.CmdExecutor
	colIDs     map[string]int
	colIDIsSet bool
}

func MakeDiskUsageStatCollector() collector.StatisticCollector {
	return &diskUsageLinux{
		cmd: standart.MakeStandardCmdExecutor("/usr/bin/df", "-k"),
	}
}

func (dul *diskUsageLinux) Collect(ctx context.Context) (statistic.Statistic, error) {
	out, err := dul.cmd.Run(ctx)
	if err != nil {
		return nil, err
	}
	return dul.extractStat(out)
}

// private section

func (dul *diskUsageLinux) extractStat(data string) (statistic.Statistic, error) {
	const (
		useColHead   = "Use%"
		fsColHead    = "Filesystem"
		mountColHead = "Mounted on"
	)
	table, err := util.ParseTable(data, lineSeparator, 0, 5)
	if err != nil {
		return nil, err
	}
	if !dul.colIDIsSet {
		dul.setCollIDs(table[0], []string{useColHead, fsColHead, mountColHead})
	}
	res := statistic.DiskUsage{
		BlockUsage: make(map[string]statistic.FSDiskInfo, len(table)-1),
	}
	for _, record := range table[1:] {
		use := record[dul.colIDs[useColHead]]
		usage, err := util.GetTrimmedFloat(use, "", "%")
		if err != nil {
			return nil, fmt.Errorf("%w: %w", collector.ErrDiskUsageInfoNotFound, err)
		}
		fs := record[dul.colIDs[fsColHead]]
		mount := record[dul.colIDs[mountColHead]]
		info := statistic.FSDiskInfo{
			FS:    fs,
			Usage: usage,
		}
		res.BlockUsage[mount] = info
	}

	return res, nil
}

func (dul *diskUsageLinux) setCollIDs(header, heads []string) error {
	dul.colIDs = make(map[string]int, len(heads))
	headsIDs, err := util.GetHeadsIDs(header, heads)
	if err != nil {
		return err
	}
	for _, head := range heads {
		id, ok := headsIDs[head]
		if !ok {
			return fmt.Errorf("can't find head %q in header %v", head, header)
		}
		dul.colIDs[head] = id
	}
	return nil
}
