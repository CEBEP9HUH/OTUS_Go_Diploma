package linux

import (
	"context"
	"errors"
	"strings"

	"github.com/CEBEP9HUH/OTUS_Go_Diploma/internal/services/cmdexecutor"
	"github.com/CEBEP9HUH/OTUS_Go_Diploma/internal/services/cmdexecutor/standart"
	"github.com/CEBEP9HUH/OTUS_Go_Diploma/internal/services/collector"
	"github.com/CEBEP9HUH/OTUS_Go_Diploma/internal/services/statistic"
	"github.com/CEBEP9HUH/OTUS_Go_Diploma/internal/services/util"
)

type loadAvgLinux struct {
	cmd cmdexecutor.CmdExecutor
}

func MakeLoadAvgStatCollector() collector.StatisticCollector {
	return &loadAvgLinux{
		cmd: standart.MakeStandartCmdExecutor("/usr/bin/top", "-b", "-n1"),
	}
}

func (lal *loadAvgLinux) Collect(ctx context.Context) (statistic.Statistic, error) {
	out, err := lal.cmd.Run(ctx)
	if err != nil {
		return nil, err
	}
	return lal.extractStat(out)
}

// private section

func (lal *loadAvgLinux) extractStat(data string) (statistic.Statistic, error) {
	const (
		loadAvgPrefix = "load average:"
	)
	loadAvgLine, err := util.GetLineInfo(data, loadAvgPrefix, lineSeparator)
	if err != nil {
		return nil, collector.ErrLoadAvgInfoNotFound
	}

	return extractloadAvgUsage(loadAvgLine)
}

func extractloadAvgUsage(loadAvgLine string) (statistic.Statistic, error) {
	const (
		min1Id = iota
		min5Id
		min15Id
		size
	)
	loadAvgLine = strings.TrimSpace(loadAvgLine)
	if len(loadAvgLine) == 0 {
		return nil, collector.ErrLoadAvgInfoNotFound
	}
	var res statistic.LoadAvg
	loadAvgInfo := strings.Split(loadAvgLine, " ")
	if len(loadAvgInfo) != size {
		return nil, collector.ErrLoadAvgInfoNotFound
	}
	var parseErr error
	var dest *float32
	for i, info := range loadAvgInfo {
		switch i {
		case min1Id:
			parseErr = collector.ErrAvgMin1NotFound
			dest = &res.Min1
		case min5Id:
			parseErr = collector.ErrAvgMin5NotFound
			dest = &res.Min5
		case min15Id:
			parseErr = collector.ErrAvgMin15NotFound
			dest = &res.Min15
		default:
			continue
		}
		v, err := util.GetTrimmedFloat(info, "", ",")
		if err != nil {
			return nil,
				errors.Join(collector.ErrLoadAvgInfoNotFound, parseErr, err)
		}
		*dest = v
	}
	return res, nil
}
