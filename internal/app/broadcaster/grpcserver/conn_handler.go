package grpcserver

import (
	"errors"
	"time"

	"github.com/CEBEP9HUH/OTUS_Go_Diploma/internal/services/loggerwrapper"
	"github.com/CEBEP9HUH/OTUS_Go_Diploma/internal/services/statistic"
	"github.com/CEBEP9HUH/OTUS_Go_Diploma/internal/services/timedlist"
	"github.com/CEBEP9HUH/OTUS_Go_Diploma/pkg/api/grpc/statistic/broadcaster"
	"github.com/CEBEP9HUH/OTUS_Go_Diploma/pkg/api/grpc/statistic/data"
	"google.golang.org/grpc"
)

var ErrNotEnoughData = errors.New("not enough data")

type connHandler struct {
	storage          timedlist.TimedList[statistic.Snapshot]
	loger            loggerwrapper.Logger
	collectingPeriod time.Duration
	sendingPeriod    time.Duration
	timeShift        time.Duration
}

func (h *connHandler) handle(stream grpc.ServerStreamingServer[broadcaster.SysStat]) error {
	const (
		statErrTemplate = "Can't send statistic: %s"
	)
	h.loger.Info("Start statistic sending")
	defer h.loger.Info("Stop statistic sending")
	for {
		select {
		case <-stream.Context().Done():
			return nil
		case <-time.After(h.sendingPeriod):
			stat, err := h.getStatistic()
			if err != nil {
				h.loger.Warn(statErrTemplate, err)
				continue
			}
			if err := stream.Send(stat); err != nil {
				h.loger.Warn(statErrTemplate, err)
			}
		}
	}
}

func (h *connHandler) getStatistic() (*broadcaster.SysStat, error) {
	statData, oldest := h.storage.GetLast(h.collectingPeriod)
	if time.Until(oldest) > h.timeShift-h.collectingPeriod {
		return nil, ErrNotEnoughData
	}
	ret := &broadcaster.SysStat{
		CpuUsage:  makeCPUUSage(statData),
		LoadAvg:   makeLoadAvg(statData),
		IoStat:    makeIOStat(statData),
		NodeUsage: makeNodeUsage(statData),
		DiskUsage: makeDiskUsage(statData),
	}

	return ret, nil
}

func makeCPUUSage(statData []statistic.Snapshot) *data.CPUUsage {
	var ret data.CPUUsage
	size := float32(len(statData))
	for _, v := range statData {
		ret.Idle += v.CPUUsage.Idle / size
		ret.SystemMode += v.CPUUsage.SystemMode / size
		ret.UserMode += v.CPUUsage.UserMode / size
	}
	return &ret
}

func makeLoadAvg(statData []statistic.Snapshot) *data.LoadAvg {
	var ret data.LoadAvg
	size := float32(len(statData))
	for _, v := range statData {
		ret.Min1 += v.LoadAvg.Min1 / size
		ret.Min5 += v.LoadAvg.Min5 / size
		ret.Min15 += v.LoadAvg.Min15 / size
	}
	return &ret
}

func makeIOStat(statData []statistic.Snapshot) *data.IOStat {
	ret := data.IOStat{
		Usage: make(map[string]*data.IOStat_Info),
	}
	size := float32(len(statData))
	for _, v := range statData {
		for dev, info := range v.IOStat.Data {
			ret.Usage[dev] = &data.IOStat_Info{
				Tps:      info.TPS / size,
				ReadKbs:  info.ReadKbs / size,
				WriteKbs: info.WriteKbs / size,
			}
		}
	}
	return &ret
}

func makeNodeUsage(statData []statistic.Snapshot) *data.NodeUsage {
	ret := data.NodeUsage{
		Usage: make(map[string]*data.NodeUsage_Usage),
	}
	size := float32(len(statData))
	for _, v := range statData {
		for mount, usage := range v.NodeUsage.NodeUsage {
			if v, ok := ret.Usage[mount]; ok {
				ret.Usage[mount] = &data.NodeUsage_Usage{
					Fs:      v.Fs,
					Percent: v.Percent + usage.Usage/size,
				}
				continue
			}
			ret.Usage[mount] = &data.NodeUsage_Usage{
				Fs:      usage.FS,
				Percent: usage.Usage / size,
			}
		}
	}
	return &ret
}

func makeDiskUsage(statData []statistic.Snapshot) *data.DiskUsage {
	ret := data.DiskUsage{
		Usage: make(map[string]*data.DiskUsage_Usage),
	}
	size := float32(len(statData))
	for _, v := range statData {
		for mount, usage := range v.DiskUsage.BlockUsage {
			if v, ok := ret.Usage[mount]; ok {
				ret.Usage[mount] = &data.DiskUsage_Usage{
					Fs:      v.Fs,
					Percent: v.Percent + usage.Usage/size,
				}
				continue
			}
			ret.Usage[mount] = &data.DiskUsage_Usage{
				Fs:      usage.FS,
				Percent: usage.Usage / size,
			}
		}
	}
	return &ret
}
