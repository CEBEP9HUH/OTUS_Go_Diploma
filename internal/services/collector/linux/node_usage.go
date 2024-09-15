//go:build linux
// +build linux

package linux

import (
	"context"
	"errors"
	"fmt"

	"github.com/CEBEP9HUH/OTUS_Go_Diploma/internal/services/cmdexecutor"
	"github.com/CEBEP9HUH/OTUS_Go_Diploma/internal/services/cmdexecutor/standart"
	"github.com/CEBEP9HUH/OTUS_Go_Diploma/internal/services/collector"
	"github.com/CEBEP9HUH/OTUS_Go_Diploma/internal/services/statistic"
	"github.com/CEBEP9HUH/OTUS_Go_Diploma/internal/services/util"
)

type nodeUsageLinux struct {
	cmd        cmdexecutor.CmdExecutor
	colIDs     map[string]int
	colIDIsSet bool
}

func MakeNodeUsageStatCollector() collector.StatisticCollector {
	return &nodeUsageLinux{
		cmd: standart.MakeStandartCmdExecutor("/usr/bin/df", "-i"),
	}
}

func (nul *nodeUsageLinux) Collect(ctx context.Context) (statistic.Statistic, error) {
	out, err := nul.cmd.Run(ctx)
	if err != nil {
		return nil, err
	}
	return nul.extractStat(out)
}

// private section

func (nul *nodeUsageLinux) extractStat(data string) (statistic.Statistic, error) {
	const (
		useColHead   = "IUse%"
		fsColHead    = "Filesystem"
		mountColHead = "Mounted on"
	)
	table, err := util.ParseTable(data, lineSeparator, 0, 5)
	if err != nil {
		return nil, err
	}
	if !nul.colIDIsSet {
		err := nul.setCollIDs(table[0], []string{useColHead, fsColHead, mountColHead})
		if err != nil {
			return nil, errors.Join(collector.ErrNodeUsageInfoNotFound, err)
		}
	}

	res := statistic.NodeUsage{
		NodeUsage: make(map[string]statistic.FSNodeInfo, len(table)-1),
	}
	for _, record := range table[1:] {
		use := record[nul.colIDs[useColHead]]
		usage, err := util.GetTrimmedFloat(use, "", "%")
		if err != nil {
			return nil, errors.Join(collector.ErrNodeUsageInfoNotFound, err)
		}
		fs := record[nul.colIDs[fsColHead]]
		mount := record[nul.colIDs[mountColHead]]
		info := statistic.FSNodeInfo{
			FS:    fs,
			Usage: usage,
		}
		res.NodeUsage[mount] = info
	}

	return res, nil
}

func (nul *nodeUsageLinux) setCollIDs(header, heads []string) error {
	nul.colIDs = make(map[string]int, len(heads))
	headsIDs, err := util.GetHeadsIDs(header, heads)
	if err != nil {
		return err
	}
	for _, head := range heads {
		id, ok := headsIDs[head]
		if !ok {
			return fmt.Errorf("can't find head %q in header %v", head, header)
		}
		nul.colIDs[head] = id
	}
	return nil
}
