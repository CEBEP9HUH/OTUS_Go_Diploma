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

type ioStatLinux struct {
	cmd cmdexecutor.CmdExecutor

	colIDs     map[string]int
	colIDIsSet bool
}

func MakeIOStatStatCollector() collector.StatisticCollector {
	return &ioStatLinux{
		cmd: standart.MakeStandartCmdExecutor("/usr/bin/iostat", "-d", "-k"),
	}
}

func (iosl *ioStatLinux) Collect(ctx context.Context) (statistic.Statistic, error) {
	out, err := iosl.cmd.Run(ctx)
	if err != nil {
		return nil, err
	}
	return iosl.extractStat(out)
}

// private section

func (iosl *ioStatLinux) extractStat(data string) (statistic.Statistic, error) {
	const (
		devColHead      = "Device"
		tpsColHead      = "tps"
		readKbsColHead  = "kB_read/s"
		writeKbsColHead = "kB_wrtn/s"
	)
	table, err := util.ParseTable(data, lineSeparator, 2, 5)
	if err != nil {
		return nil, err
	}
	if !iosl.colIDIsSet {
		iosl.setCollIDs(table[0], []string{devColHead, tpsColHead, readKbsColHead, writeKbsColHead})
	}
	res := statistic.IOStat{
		Data: make(map[string]statistic.DevIOStat, len(table[1:])),
	}
	for _, record := range table[1:] {
		tps, err := util.GetFloat(record[iosl.colIDs[tpsColHead]])
		if err != nil {
			return nil, fmt.Errorf("%w: %w", collector.ErrTPSInfoNotFound, err)
		}
		readKbs, err := util.GetFloat(record[iosl.colIDs[readKbsColHead]])
		if err != nil {
			return nil, fmt.Errorf("%w: %w", collector.ErrReadKbsInfoNotFound, err)
		}
		writeKbs, err := util.GetFloat(record[iosl.colIDs[writeKbsColHead]])
		if err != nil {
			return nil, fmt.Errorf("%w: %w", collector.ErrWriteKbsInfoNotFound, err)
		}

		info := statistic.DevIOStat{
			TPS:      tps,
			ReadKbs:  readKbs,
			WriteKbs: writeKbs,
		}
		res.Data[record[iosl.colIDs[devColHead]]] = info
	}

	return res, nil
}

func (iosl *ioStatLinux) setCollIDs(header, heads []string) error {
	iosl.colIDs = make(map[string]int, len(heads))
	headsIDs, err := util.GetHeadsIDs(header, heads)
	if err != nil {
		return err
	}
	for _, head := range heads {
		id, ok := headsIDs[head]
		if !ok {
			return fmt.Errorf("can't find head %q in header %v", head, header)
		}
		iosl.colIDs[head] = id
	}
	return nil
}
