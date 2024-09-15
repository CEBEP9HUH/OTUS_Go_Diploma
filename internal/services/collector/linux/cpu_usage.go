//go:build linux
// +build linux

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

type cpuUsageLinux struct {
	cmd cmdexecutor.CmdExecutor
}

func MakeCPUUsageStatCollector() collector.StatisticCollector {
	return &cpuUsageLinux{
		cmd: standart.MakeStandartCmdExecutor("/usr/bin/top", "-b", "-n1"),
	}
}

func (cul *cpuUsageLinux) Collect(ctx context.Context) (statistic.Statistic, error) {
	out, err := cul.cmd.Run(ctx)
	if err != nil {
		return nil, err
	}
	return cul.extractStat(out)
}

// private section

func (cul *cpuUsageLinux) extractStat(data string) (statistic.Statistic, error) {
	const (
		cpuPrefix = "cpu(s):"
	)
	cpuLine, err := util.GetLineInfo(data, cpuPrefix, lineSeparator)
	if err != nil {
		return nil, collector.ErrCPUInfoNotFound
	}

	return extractCPUUsage(cpuLine)
}

func extractCPUUsage(cpuLine string) (statistic.Statistic, error) {
	const (
		userMode   = "us"
		systemMode = "sy"
		idle       = "id"
	)
	if len(cpuLine) == 0 {
		return nil, collector.ErrCPUInfoNotFound
	}
	var res statistic.CPUUsage
	params := strings.Split(cpuLine, ", ")
	f := func(params []string, paramType string, dst *float32, err error) error {
		for _, param := range params {
			if strings.Contains(param, paramType) {
				v, parseErr := util.GetTrimmedFloat(param, "", paramType)
				if parseErr != nil {
					return errors.Join(err, parseErr)
				}
				*dst = v
				return nil
			}
		}
		return err
	}
	err := f(params, userMode, &res.UserMode, collector.ErrUserModeNotFound)
	if err != nil {
		return nil, errors.Join(collector.ErrCPUInfoNotFound, err)
	}
	err = f(params, systemMode, &res.SystemMode, collector.ErrSystemModeNotFound)
	if err != nil {
		return nil, errors.Join(collector.ErrCPUInfoNotFound, err)
	}
	err = f(params, idle, &res.Idle, collector.ErrIdleNotFound)
	if err != nil {
		return nil, errors.Join(collector.ErrCPUInfoNotFound, err)
	}

	return res, nil
}
